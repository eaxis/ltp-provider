package cache

import (
	"github.com/eaxis/ltp-provider/internal/domain"
	"sync"
	"time"
)

type item struct {
	ltp       domain.Ltp
	expiresAt time.Time
}

type LTPCache struct {
	data       map[string]item
	mu         sync.Mutex
	expiration time.Duration
}

func NewLTPCache(expiration time.Duration) *LTPCache {
	return &LTPCache{
		data:       make(map[string]item),
		expiration: expiration, // 1 * time.Minute
	}
}

func (s *LTPCache) Get(pair string) (*domain.Ltp, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, found := s.data[pair]

	if !found {
		return nil, false
	}

	if time.Now().After(entry.expiresAt) {
		delete(s.data, pair)
		return nil, false
	}

	return &entry.ltp, true
}

func (s *LTPCache) Set(ltp domain.Ltp) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[ltp.Pair] = item{
		ltp:       ltp,
		expiresAt: time.Now().Add(s.expiration),
	}
}
