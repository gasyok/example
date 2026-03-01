package httputil

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"example/internal/domain"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func Decode(r *http.Request, v any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

func HandleFlow(w http.ResponseWriter, log *slog.Logger, flow func() error) {
	switch err := flow(); {
	case err == nil:
		return
	case errors.Is(err, domain.ErrUnauthorized):
		log.Warn("unauthorized", slog.String("error", err.Error()))
		JSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	case errors.Is(err, domain.ErrInvalidInput):
		log.Warn("bad request", slog.String("error", err.Error()))
		JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	case errors.Is(err, domain.ErrNotFound):
		log.Warn("not found", slog.String("error", err.Error()))
		JSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
	case errors.Is(err, domain.ErrForbidden):
		log.Warn("forbidden", slog.String("error", err.Error()))
		JSON(w, http.StatusForbidden, map[string]string{"error": "forbidden"})
	case errors.Is(err, domain.ErrAlreadyExists):
		log.Warn("conflict", slog.String("error", err.Error()))
		JSON(w, http.StatusConflict, map[string]string{"error": "already exists"})
	default:
		log.Error("internal error", slog.String("error", err.Error()))
		JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
	}
}
