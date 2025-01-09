package server

import (
	"errors"
	"fmt"
	"github.com/eaxis/ltp-provider/internal/common/server"
	"github.com/eaxis/ltp-provider/internal/domain"
	"net/http"
	"strings"
)

func (s *HttpServer) GetLtp(w http.ResponseWriter, r *http.Request) {
	const op = "server.HttpServer.GetLtp"

	pairsQuery := r.URL.Query().Get("pairs")

	if pairsQuery == "" {
		pairsQuery = "BTC/USD,BTC/CHF,BTC/EUR" // default if none specified
	}

	pairs := strings.Split(pairsQuery, ",")

	if len(pairs) == 0 {
		server.BadRequest("no pairs specified", fmt.Errorf("(%s) no pairs specified", op), w, r)
	}

	entries, err := s.ltpService.Retrieve(r.Context(), pairs)

	if err != nil {
		if errors.Is(err, domain.ErrNotSupported) {
			wrappedErr := fmt.Errorf("(%s) pair(-s) %s aint not supported", op, pairsQuery)

			server.NotFound("some of the given pairs aren't supported yet", wrappedErr, w, r)
			return
		}
		server.InternalError("internal error", err, w, r)
		return
	}

	response := toResponseLtp(entries)

	server.RespondOK(response, w, r)
}
