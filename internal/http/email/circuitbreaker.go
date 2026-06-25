package email

import (
	"context"
	"log"

	"github.com/sony/gobreaker/v2"
)

// CBService wraps an email Service with a circuit breaker.
// When the breaker is open (too many consecutive failures), calls fail fast
// without hitting the underlying service.
type CBService struct {
	inner Mailer
	cb    *gobreaker.CircuitBreaker[struct{}]
}

func NewCBService(inner Mailer) *CBService {
	settings := gobreaker.Settings{
		Name:        "email-service",
		MaxRequests: 1,
		Interval:    0,
		Timeout:     30_000_000_000,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 3
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			log.Printf("[circuit-breaker] %s: %s → %s", name, from, to)
		},
	}

	return &CBService{
		inner: inner,
		cb:    gobreaker.NewCircuitBreaker[struct{}](settings),
	}
}

func (s *CBService) SendInvitation(ctx context.Context, toEmail, teamName string) error {
	_, err := s.cb.Execute(func() (struct{}, error) {
		return struct{}{}, s.inner.SendInvitation(ctx, toEmail, teamName)
	})
	return err
}
