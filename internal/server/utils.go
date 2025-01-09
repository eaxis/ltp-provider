package server

import (
	"github.com/eaxis/ltp-provider/internal/domain"
)

func toResponseLtp(entries []domain.Ltp) LtpResponse {
	return LtpResponse{
		Entries: entries,
	}
}
