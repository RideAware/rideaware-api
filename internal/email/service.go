package email

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strconv"
)

type Service struct {
	smtpServer   string
	smtpPort     int
	smtpUser     string
	smtpPassword string
	from         string
}

func NewService() *Service {
	smtpServer := os.Getenv("SMTP_SERVER")
	if smtpServer == "" {
		smtpServer = "localhost"
	}

	smtpPort := os.Getenv("SMTP_PORT")
	if smtpPort == "" {
		smtpPort = "587"
	}

	port, err := strconv.Atoi(smtpPort)
	if err != nil {
		port = 587
	}

	smtpUser := os.Getenv("SMTP_USER")
	if smtpUser == "" {
		smtpUser = "noreply@rideaware.app"
	}

	smtpPassword := os.Getenv("SMTP_PASSWORD")

	from := os.Getenv("SENDER_EMAIL")
	if from == "" {
		from = "noreply@rideaware.app"
	}

	log.Printf("📧 Email service initialized: %s@%s:%d", smtpUser, smtpServer, port)

	return &Service{
		smtpServer:   smtpServer,
		smtpPort:     port,
		smtpUser:     smtpUser,
		smtpPassword: smtpPassword,
		from:         from,
	}
}

func (s *Service) sendEmail(to []string, subject, htmlBody string) error {
	log.Printf("📧 Preparing to send email to: %s (Subject: %s)", to[0], subject)

	// Create message
	headers := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n",
		s.from,
		to[0],
		subject,
	)

	message := headers + htmlBody

	// SMTP server address
	addr := fmt.Sprintf("%s:%d", s.smtpServer, s.smtpPort)
	log.Printf("📧 Connecting to SMTP: %s", addr)

	// TLS configuration
	tlsConfig := &tls.Config{
		ServerName: s.smtpServer,
	}

	// Create connection
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		log.Printf("❌ TLS connection failed: %v", err)
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	log.Printf("📧 TLS connection established")

	// Create SMTP client
	client, err := smtp.NewClient(conn, s.smtpServer)
	if err != nil {
		log.Printf("❌ SMTP client creation failed: %v", err)
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	log.Printf("📧 SMTP client created")

	// Authenticate
	auth := smtp.PlainAuth("", s.smtpUser, s.smtpPassword, s.smtpServer)
	if err = client.Auth(auth); err != nil {
		log.Printf("❌ SMTP authentication failed: %v", err)
		return fmt.Errorf("SMTP authentication failed: %w", err)
	}

	log.Printf("📧 SMTP authentication successful")

	// Send email
	if err = client.Mail(s.from); err != nil {
		log.Printf("❌ SMTP Mail command failed: %v", err)
		return fmt.Errorf("SMTP Mail command failed: %w", err)
	}

	if err = client.Rcpt(to[0]); err != nil {
		log.Printf("❌ SMTP Rcpt command failed: %v", err)
		return fmt.Errorf("SMTP Rcpt command failed: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		log.Printf("❌ SMTP Data command failed: %v", err)
		return fmt.Errorf("SMTP Data command failed: %w", err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Printf("❌ Writing email body failed: %v", err)
		return fmt.Errorf("writing email body failed: %w", err)
	}

	err = w.Close()
	if err != nil {
		log.Printf("❌ Closing email data failed: %v", err)
		return fmt.Errorf("closing email data failed: %w", err)
	}

	client.Quit()

	log.Printf("✅ Email sent successfully to: %s", to[0])
	return nil
}

func (s *Service) SendPasswordResetEmail(email, username, resetLink string) error {
	log.Printf("🔑 Sending password reset email to: %s", email)

	subject := "Reset Your RideAware Password"
	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.header { background: linear-gradient(135deg, #1e4e9c 0%, #337cf2 100%); color: white; padding: 20px; border-radius: 8px; }
				.content { padding: 20px; background: #f9f9f9; margin: 20px 0; border-radius: 8px; }
				.button { background: linear-gradient(135deg, #1e4e9c 0%, #337cf2 100%); color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; display: inline-block; margin: 20px 0; }
				.footer { text-align: center; color: #666; font-size: 12px; margin-top: 20px; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h2>Password Reset Request</h2>
				</div>
				<div class="content">
					<p>Hi %s,</p>
					<p>We received a request to reset your password. Click the button below to create a new password:</p>
					<p><a href="%s" class="button">Reset Password</a></p>
					<p><strong>Note:</strong> This link will expire in 1 hour.</p>
					<p>If you didn't request this, you can safely ignore this email.</p>
				</div>
				<div class="footer">
					<p>&copy; 2025 RideAware. All rights reserved.</p>
				</div>
			</div>
		</body>
		</html>
	`, username, resetLink)

	return s.sendEmail([]string{email}, subject, htmlBody)
}

func (s *Service) SendWelcomeEmail(email, username string) error {
	log.Printf("👋 Sending welcome email to: %s", email)

	subject := "Welcome to RideAware"
	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.header { background: linear-gradient(135deg, #1e4e9c 0%, #337cf2 100%); color: white; padding: 20px; border-radius: 8px; }
				.content { padding: 20px; background: #f9f9f9; margin: 20px 0; border-radius: 8px; }
				.button { background: linear-gradient(135deg, #1e4e9c 0%, #337cf2 100%); color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; display: inline-block; margin: 20px 0; }
				.footer { text-align: center; color: #666; font-size: 12px; margin-top: 20px; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h2>Welcome to RideAware</h2>
				</div>
				<div class="content">
					<p>Hi %s,</p>
					<p>Your account has been created successfully! 🚀</p>
					<p>You're now ready to:</p>
					<ul>
						<li>Track your cycling performance</li>
						<li>Manage your equipment</li>
						<li>Create custom training zones</li>
						<li>Plan structured workouts</li>
					</ul>
					<p>Get started by logging in to your account and setting up your profile.</p>
					<p><a href="https://dev.rideaware.org" class="button">Go to RideAware</a></p>
				</div>
				<div class="footer">
					<p>&copy; 2025 RideAware. All rights reserved.</p>
				</div>
			</div>
		</body>
		</html>
	`, username)

	return s.sendEmail([]string{email}, subject, htmlBody)
}

func (s *Service) SendNewsletterEmail(email, subject, htmlBody string) error {
	log.Printf("📬 Sending newsletter email to: %s", email)
	return s.sendEmail([]string{email}, subject, htmlBody)
}