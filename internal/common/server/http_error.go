package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func InternalError(message string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, message, w, r, "Internal server error", http.StatusInternalServerError)
}

/*
func Unauthorised(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, "Unauthorised", http.StatusUnauthorized)
}
*/

func BadRequest(message string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, message, w, r, "Bad request", http.StatusBadRequest)
}

func NotFound(message string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, message, w, r, "Not found", http.StatusBadRequest)
}

func httpRespondWithError(err error, message string, w http.ResponseWriter, r *http.Request, msg string, status int) {
	log.Printf("error: %s, slug: %s, msg: %s", err, message, msg)

	resp := ErrorResponse{Message: message, httpStatus: status}
	if os.Getenv("DEBUG_ERRORS") != "" && err != nil {
		resp.Error = err.Error()
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}

type ErrorResponse struct {
	Message    string `json:"message"`
	Error      string `json:"error,omitempty"`
	httpStatus int
}

func (e ErrorResponse) Render(w http.ResponseWriter, _ *http.Request) error {
	w.WriteHeader(e.httpStatus)
	return nil
}
