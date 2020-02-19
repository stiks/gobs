package mock

import (
	"context"

	"github.com/matcornic/hermes/v2"

	"github.com/stiks/gobs/lib/repositories"
)

type emailRepository struct {
}

// NewEmailRepository ...
func NewEmailRepository() repositories.EmailRepository {
	return &emailRepository{}
}

// SendEmail ...
func (e *emailRepository) SendEmail(ctx context.Context, to string, subject string, email hermes.Email) error {
	return nil
}
