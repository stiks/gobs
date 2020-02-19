package services

import (
	"context"

	"github.com/matcornic/hermes/v2"

	"github.com/stiks/gobs/lib/repositories"
)

type emailService struct {
	repo repositories.EmailRepository
}

// EmailService ...
type EmailService interface {
	SendEmail(ctx context.Context, to string, subject string, email hermes.Email) error
}

// NewEmailService ...
func NewEmailService(repo repositories.EmailRepository) EmailService {
	return &emailService{
		repo: repo,
	}
}

// SendEmail ...
func (s *emailService) SendEmail(ctx context.Context, to string, subject string, email hermes.Email) error {
	return s.repo.SendEmail(ctx, to, subject, email)
}
