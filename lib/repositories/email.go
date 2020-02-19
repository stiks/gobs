package repositories

import (
	"context"

	"github.com/matcornic/hermes/v2"
)

// EmailRepository ...
type EmailRepository interface {
	SendEmail(ctx context.Context, to string, subject string, email hermes.Email) error
}
