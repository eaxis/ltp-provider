package server

import (
	"context"
	"github.com/eaxis/ltp-provider/internal/domain"
)

type LtpService interface {
	Retrieve(ctx context.Context, pairs []string) ([]domain.Ltp, error)
}
