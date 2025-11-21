package email

import (
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

type Service struct {
	client *resend.Client
	from   string
}

func NewService() *Service {
	senderEmail := os.Getenv("SENDER_EMAIL")
	if senderEmail == "" {
		senderEmail = "noreply@rideaware.app"
	}

	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		apiKey = "re_test"
	}

	return &Service{
		client: resend.NewClient(apiKey),
		from:   senderEmail,
	}
}

func (s *Service) SendPasswordResetEmail(email, username, resetLink string) error {
	params := &resend.SendEmailRequest{
		From:    s.from,
		To:      []string{email},
		Subject: "Reset Your RideAware Password",
		Html: fmt.Sprintf(`
			<h2>Password Reset Request</h2>
			<p>Hi %s,</p>
			<p>We received a request to reset your password. Click the link below to create a new password:</p>
			<p><a href="%s">Reset Password</a></p>
			<p>This link will expire in 1 hour.</p>
			<p>If you didn't request this, you can ignore this email.</p>
		`, username, resetLink),
	}

	sent, err := s.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if sent.Id == "" {
		return fmt.Errorf("failed to send email")
	}

	return nil
}

func (s *Service) SendWelcomeEmail(email, username string) error {
	params := &resend.SendEmailRequest{
		From:    s.from,
		To:      []string{email},
		Subject: "Welcome to RideAware",
		Html: fmt.Sprintf(`
			<h2>Welcome to RideAware</h2>
			<p>Hi %s,</p>
			<p>Your account has been created successfully!</p>
			<p>Start tracking your rides and improve your performance.</p>
		`, username),
	}

	sent, err := s.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if sent.Id == "" {
		return fmt.Errorf("failed to send email")
	}

	return nil
}