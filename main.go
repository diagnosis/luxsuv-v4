package main

import (
	"database/sql"
	"fmt"
	"log"
	"luxsuv-v4/handlers"
	"luxsuv-v4/middleware"
	"luxsuv-v4/services"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	// Database connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://username:password@localhost/luxsuv?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// JWT secret
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production"
		log.Println("Warning: Using default JWT secret. Set JWT_SECRET environment variable in production.")
	}

	// Initialize services
	authService := services.NewAuthService(db, jwtSecret)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)
	
	// Rate limiter: 5 requests per second with burst of 10
	rateLimiter := middleware.NewRateLimiter(200*time.Millisecond, 10)

	// Routes
	mux := http.NewServeMux()

	// Public routes (with rate limiting)
	mux.HandleFunc("/auth/register", rateLimiter.RateLimit(authHandler.Register))
	mux.HandleFunc("/auth/login", rateLimiter.RateLimit(authHandler.Login))

	// Protected routes
	mux.HandleFunc("/users/me", authMiddleware.RequireAuth(authHandler.GetCurrentUser))

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
	})

	// CORS middleware
	corsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      corsHandler(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(server.ListenAndServe())
}