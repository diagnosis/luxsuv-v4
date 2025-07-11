package email

import (
	"fmt"
	"crypto/tls"

	"github.com/diagnosis/luxsuv-v4/internal/logger"
	"gopkg.in/gomail.v2"
)

type Service struct {
	host     string
	port     int
	username string
	password string
	from     string
	logger   *logger.Logger
}

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func NewService(config Config, logger *logger.Logger) *Service {
	return &Service{
		host:     config.Host,
		port:     config.Port,
		username: config.Username,
		password: config.Password,
		from:     config.From,
		logger:   logger,
	}
}

// SendPasswordResetEmail sends a password reset email
func (s *Service) SendPasswordResetEmail(to, resetToken string) error {
	subject := "Password Reset Request"
	
	// Create the reset URL (you'll need to adjust this based on your frontend URL)
	resetURL := fmt.Sprintf("http://localhost:3000/reset-password?token=%s", resetToken)
	
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Password Reset</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #f8f9fa; padding: 20px; text-align: center; border-radius: 5px; }
        .content { padding: 20px 0; }
        .button { 
            display: inline-block; 
            padding: 12px 24px; 
            background-color: #007bff; 
            color: white; 
            text-decoration: none; 
            border-radius: 5px; 
            margin: 20px 0;
        }
        .footer { font-size: 12px; color: #666; margin-top: 30px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Password Reset Request</h1>
        </div>
        <div class="content">
            <p>Hello,</p>
            <p>You have requested to reset your password. Click the button below to reset your password:</p>
            <p style="text-align: center;">
                <a href="%s" class="button">Reset Password</a>
            </p>
            <p>If the button doesn't work, you can copy and paste this link into your browser:</p>
            <p><a href="%s">%s</a></p>
            <p><strong>This link will expire in 1 hour.</strong></p>
            <p>If you didn't request this password reset, please ignore this email.</p>
        </div>
        <div class="footer">
            <p>This is an automated message, please do not reply to this email.</p>
        </div>
    </div>
</body>
</html>
	`, resetURL, resetURL, resetURL)

	return s.sendEmail(to, subject, body)
}

// sendEmail sends an email using SMTP
func (s *Service) sendEmail(to, subject, body string) error {
	s.logger.Info(fmt.Sprintf("Attempting to send email to %s via %s:%d", to, s.host, s.port))
	
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(s.host, s.port, s.username, s.password)
	
	// Configure TLS
	d.TLSConfig = &tls.Config{InsecureSkipVerify: false}
	
	s.logger.Info(fmt.Sprintf("Connecting to SMTP server %s:%d with username %s", s.host, s.port, s.username))

	if err := d.DialAndSend(m); err != nil {
		s.logger.Err(fmt.Sprintf("Failed to send email to %s via %s:%d - Error: %s", to, s.host, s.port, err.Error()))
		return fmt.Errorf("failed to send email: %w", err)
	}

	s.logger.Info(fmt.Sprintf("Email sent successfully to %s via %s:%d", to, s.host, s.port))
	return nil
}

// SendWelcomeEmail sends a welcome email to new users
func (s *Service) SendWelcomeEmail(to, username string) error {
	subject := "Welcome to LuxSUV!"
	
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Welcome to LuxSUV</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #28a745; color: white; padding: 20px; text-align: center; border-radius: 5px; }
        .content { padding: 20px 0; }
        .footer { font-size: 12px; color: #666; margin-top: 30px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to LuxSUV!</h1>
        </div>
        <div class="content">
            <p>Hello %s,</p>
            <p>Welcome to LuxSUV! Your account has been successfully created.</p>
            <p>You can now log in and start using our premium ride-sharing service.</p>
            <p>If you have any questions, feel free to contact our support team.</p>
            <p>Thank you for choosing LuxSUV!</p>
        </div>
        <div class="footer">
            <p>This is an automated message, please do not reply to this email.</p>
        </div>
    </div>
</body>
</html>
	`, username)

	return s.sendEmail(to, subject, body)
}