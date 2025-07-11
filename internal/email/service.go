package email

import (
	"context"
	"fmt"
	"time"

	"github.com/diagnosis/luxsuv-v4/internal/logger"
	"github.com/mailersend/mailersend-go"
)

type Service struct {
	client    *mailersend.Mailersend
	fromEmail string
	fromName  string
	logger    *logger.Logger
}

type Config struct {
	APIKey    string
	FromEmail string
	FromName  string
}

func NewService(config Config, logger *logger.Logger) *Service {
	client := mailersend.NewMailersend(config.APIKey)
	
	return &Service{
		client:    client,
		fromEmail: config.FromEmail,
		fromName:  config.FromName,
		logger:    logger,
	}
}

// SendPasswordResetEmail sends a password reset email using MailerSend
func (s *Service) SendPasswordResetEmail(to, resetToken string) error {
	subject := "Password Reset Request - LuxSUV"
	
	// Create the reset URL (you'll need to adjust this based on your frontend URL)
	resetURL := fmt.Sprintf("http://localhost:3000/reset-password?token=%s", resetToken)
	
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Password Reset</title>
    <style>
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif; 
            line-height: 1.6; 
            color: #333; 
            margin: 0; 
            padding: 0; 
            background-color: #f8f9fa;
        }
        .container { 
            max-width: 600px; 
            margin: 40px auto; 
            background: white; 
            border-radius: 8px; 
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        .header { 
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); 
            color: white; 
            padding: 40px 30px; 
            text-align: center; 
        }
        .header h1 { 
            margin: 0; 
            font-size: 28px; 
            font-weight: 600; 
        }
        .content { 
            padding: 40px 30px; 
        }
        .content h2 {
            color: #2d3748;
            font-size: 20px;
            margin-bottom: 20px;
        }
        .content p {
            margin-bottom: 20px;
            color: #4a5568;
        }
        .button { 
            display: inline-block; 
            padding: 16px 32px; 
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); 
            color: white; 
            text-decoration: none; 
            border-radius: 6px; 
            font-weight: 600;
            text-align: center;
            margin: 20px 0;
            transition: transform 0.2s;
        }
        .button:hover {
            transform: translateY(-1px);
        }
        .link-fallback {
            background-color: #f7fafc;
            padding: 15px;
            border-radius: 6px;
            border-left: 4px solid #667eea;
            margin: 20px 0;
            word-break: break-all;
        }
        .footer { 
            background-color: #f8f9fa;
            padding: 30px;
            text-align: center;
            font-size: 14px; 
            color: #718096; 
        }
        .security-notice {
            background-color: #fff5f5;
            border: 1px solid #fed7d7;
            border-radius: 6px;
            padding: 15px;
            margin: 20px 0;
        }
        .security-notice strong {
            color: #c53030;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚗 LuxSUV</h1>
        </div>
        <div class="content">
            <h2>Password Reset Request</h2>
            <p>Hello,</p>
            <p>We received a request to reset your password for your LuxSUV account. Click the button below to create a new password:</p>
            
            <div style="text-align: center; margin: 30px 0;">
                <a href="%s" class="button">Reset My Password</a>
            </div>
            
            <p>If the button doesn't work, you can copy and paste this link into your browser:</p>
            <div class="link-fallback">
                <a href="%s" style="color: #667eea; text-decoration: none;">%s</a>
            </div>
            
            <div class="security-notice">
                <p><strong>⏰ This link will expire in 1 hour</strong> for your security.</p>
                <p>If you didn't request this password reset, please ignore this email and your password will remain unchanged.</p>
            </div>
            
            <p>For your security, this request was made from IP address and will be logged.</p>
        </div>
        <div class="footer">
            <p><strong>LuxSUV - Premium Ride Sharing</strong></p>
            <p>This is an automated message, please do not reply to this email.</p>
            <p>If you need help, contact our support team.</p>
        </div>
    </div>
</body>
</html>
	`, resetURL, resetURL, resetURL)

	text := fmt.Sprintf(`
Password Reset Request - LuxSUV

Hello,

We received a request to reset your password for your LuxSUV account.

Reset your password by visiting this link:
%s

This link will expire in 1 hour for your security.

If you didn't request this password reset, please ignore this email.

---
LuxSUV - Premium Ride Sharing
This is an automated message, please do not reply.
	`, resetURL)

	return s.sendEmail(to, subject, html, text)
}

// SendWelcomeEmail sends a welcome email to new users
func (s *Service) SendWelcomeEmail(to, username string) error {
	subject := "Welcome to LuxSUV - Your Premium Ride Experience Awaits!"
	
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Welcome to LuxSUV</title>
    <style>
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif; 
            line-height: 1.6; 
            color: #333; 
            margin: 0; 
            padding: 0; 
            background-color: #f8f9fa;
        }
        .container { 
            max-width: 600px; 
            margin: 40px auto; 
            background: white; 
            border-radius: 8px; 
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        .header { 
            background: linear-gradient(135deg, #48bb78 0%%, #38a169 100%%); 
            color: white; 
            padding: 40px 30px; 
            text-align: center; 
        }
        .header h1 { 
            margin: 0; 
            font-size: 28px; 
            font-weight: 600; 
        }
        .content { 
            padding: 40px 30px; 
        }
        .welcome-message {
            text-align: center;
            margin: 30px 0;
        }
        .welcome-message h2 {
            color: #2d3748;
            font-size: 24px;
            margin-bottom: 10px;
        }
        .features {
            background-color: #f7fafc;
            border-radius: 8px;
            padding: 30px;
            margin: 30px 0;
        }
        .feature {
            display: flex;
            align-items: center;
            margin-bottom: 20px;
        }
        .feature:last-child {
            margin-bottom: 0;
        }
        .feature-icon {
            width: 40px;
            height: 40px;
            background: linear-gradient(135deg, #48bb78 0%%, #38a169 100%%);
            border-radius: 50%%;
            display: flex;
            align-items: center;
            justify-content: center;
            margin-right: 15px;
            font-size: 18px;
        }
        .footer { 
            background-color: #f8f9fa;
            padding: 30px;
            text-align: center;
            font-size: 14px; 
            color: #718096; 
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚗 Welcome to LuxSUV!</h1>
        </div>
        <div class="content">
            <div class="welcome-message">
                <h2>Hello %s! 👋</h2>
                <p>Welcome to LuxSUV - where luxury meets convenience. Your account has been successfully created and you're ready to experience premium ride-sharing.</p>
            </div>
            
            <div class="features">
                <h3 style="color: #2d3748; margin-bottom: 20px;">What you can do now:</h3>
                
                <div class="feature">
                    <div class="feature-icon">🚙</div>
                    <div>
                        <strong>Book Premium Rides</strong><br>
                        <span style="color: #718096;">Access our fleet of luxury vehicles</span>
                    </div>
                </div>
                
                <div class="feature">
                    <div class="feature-icon">⭐</div>
                    <div>
                        <strong>Rate Your Experience</strong><br>
                        <span style="color: #718096;">Help us maintain our high standards</span>
                    </div>
                </div>
                
                <div class="feature">
                    <div class="feature-icon">🎯</div>
                    <div>
                        <strong>Track Your Rides</strong><br>
                        <span style="color: #718096;">Real-time updates and notifications</span>
                    </div>
                </div>
                
                <div class="feature">
                    <div class="feature-icon">💳</div>
                    <div>
                        <strong>Secure Payments</strong><br>
                        <span style="color: #718096;">Safe and convenient payment options</span>
                    </div>
                </div>
            </div>
            
            <p>If you have any questions or need assistance, our support team is here to help 24/7.</p>
            <p>Thank you for choosing LuxSUV for your premium transportation needs!</p>
        </div>
        <div class="footer">
            <p><strong>LuxSUV - Premium Ride Sharing</strong></p>
            <p>This is an automated message, please do not reply to this email.</p>
            <p>Need help? Contact our support team anytime.</p>
        </div>
    </div>
</body>
</html>
	`, username)

	text := fmt.Sprintf(`
Welcome to LuxSUV!

Hello %s,

Welcome to LuxSUV - where luxury meets convenience. Your account has been successfully created and you're ready to experience premium ride-sharing.

What you can do now:
• Book Premium Rides - Access our fleet of luxury vehicles
• Rate Your Experience - Help us maintain our high standards  
• Track Your Rides - Real-time updates and notifications
• Secure Payments - Safe and convenient payment options

If you have any questions or need assistance, our support team is here to help 24/7.

Thank you for choosing LuxSUV for your premium transportation needs!

---
LuxSUV - Premium Ride Sharing
This is an automated message, please do not reply.
	`, username)

	return s.sendEmail(to, subject, html, text)
}

// sendEmail sends an email using MailerSend API
func (s *Service) sendEmail(to, subject, html, text string) error {
	s.logger.Info(fmt.Sprintf("Attempting to send email to %s via MailerSend", to))
	
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	from := mailersend.From{
		Name:  s.fromName,
		Email: s.fromEmail,
	}

	recipients := []mailersend.Recipient{
		{
			Name:  to, // Use email as name if no name provided
			Email: to,
		},
	}

	message := s.client.Email.NewMessage()
	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(subject)
	message.SetHTML(html)
	message.SetText(text)
	message.SetTags([]string{"password-reset", "luxsuv"})

	s.logger.Info(fmt.Sprintf("Sending email via MailerSend API to %s", to))

	res, err := s.client.Email.Send(ctx, message)
	if err != nil {
		s.logger.Err(fmt.Sprintf("Failed to send email to %s via MailerSend - Error: %s", to, err.Error()))
		return fmt.Errorf("failed to send email: %w", err)
	}

	messageID := res.Header.Get("X-Message-Id")
	s.logger.Info(fmt.Sprintf("Email sent successfully to %s via MailerSend (Message ID: %s)", to, messageID))
	return nil
}