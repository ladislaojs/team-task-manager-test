package email

import (
	"context"
	"fmt"
	"log"
)

type Mailer interface {
	SendInvitation(ctx context.Context, toEmail, teamName string) error
}

type MockMailer struct{}

func NewMockMailer() *MockMailer {
	return &MockMailer{}
}

func (s *MockMailer) SendInvitation(_ context.Context, toEmail, teamName string) error {
	log.Printf("email invitation sent to %s for team %q", toEmail, teamName)

	if toEmail == "" {
		return fmt.Errorf("email address is empty")
	}

	return nil
}
