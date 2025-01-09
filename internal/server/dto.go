package server

import "github.com/eaxis/ltp-provider/internal/domain"

type LtpResponse struct {
	Entries []domain.Ltp `json:"ltp"`
}
