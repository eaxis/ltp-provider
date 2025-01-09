package services

import (
	"context"
	"fmt"
	"github.com/eaxis/ltp-provider/internal/cache"
	"github.com/eaxis/ltp-provider/internal/domain"
	"github.com/eaxis/ltp-provider/internal/kraken"
	"strings"
)

// should this be un-hardcoded (moved to config)?
var supportedPairs = []string{
	"BTC/USD", "BTC/EUR", "BTC/CHF",
}

type LtpService struct {
	cache  *cache.LTPCache
	client *kraken.Client
}

func NewLtpService(ltpCache *cache.LTPCache, client *kraken.Client) *LtpService {
	return &LtpService{
		cache:  ltpCache,
		client: client,
	}
}

func (s *LtpService) Retrieve(ctx context.Context, pairs []string) (entries []domain.Ltp, err error) {
	const op = "services.LtpService.Retrieve"

	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)

		if !s.isPairSupported(pair) {
			return nil, domain.ErrNotSupported
		}

		if entry, found := s.cache.Get(pair); found {
			entries = append(entries, *entry)
			continue
		}

		price, krakenErr := s.client.GetLastTradedPrice(ctx, pair)

		if krakenErr != nil {
			return nil, fmt.Errorf("(%s) cannot obtain relevant data: %w", op, krakenErr)
		}

		entry := domain.Ltp{
			Pair:   pair,
			Amount: price,
		}

		s.cache.Set(entry)

		entries = append(entries, entry)
	}

	return entries, nil
}

func (s *LtpService) isPairSupported(pair string) bool {
	for _, supportedPair := range supportedPairs {
		if pair == supportedPair {
			return true
		}
	}

	return false
}
