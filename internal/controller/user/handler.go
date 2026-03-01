package user

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"example/internal/domain"
	svcuser "example/internal/service/user"
	"example/internal/utils/httputil"
	"example/pkg/dto"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc *svcuser.Service
	log *slog.Logger
}

func New(svc *svcuser.Service, log *slog.Logger) *Handler {
	return &Handler{svc: svc, log: log}
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	httputil.HandleFlow(w, h.log, func() error {
		var req dto.ListRange
		if err := httputil.Decode(r, &req); err != nil {
			return fmt.Errorf("%w: invalid body", domain.ErrInvalidInput)
		}
		if err := validateListRange(&req); err != nil {
			return fmt.Errorf("%w: invalid range", domain.ErrInvalidInput)
		}
		users, err := h.svc.List(r.Context(), req.Limit, req.Offset)
		if err != nil {
			return fmt.Errorf("svc.List: %w", err)
		}

		httputil.JSON(w, http.StatusOK, toUserListResponse(users))
		return nil
	})
}

func (h *Handler) getByID(w http.ResponseWriter, r *http.Request) {
	httputil.HandleFlow(w, h.log, func() error {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			return fmt.Errorf("%w: invalid id", domain.ErrInvalidInput)
		}
		u, err := h.svc.GetByID(r.Context(), id)
		if err != nil {
			return fmt.Errorf("svc.GetByID: %w", err)
		}
		httputil.JSON(w, http.StatusOK, toUserResponse(u))
		return nil
	})
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	httputil.HandleFlow(w, h.log, func() error {
		var req dto.CreateUserRequest
		if err := httputil.Decode(r, &req); err != nil {
			return fmt.Errorf("%w: invalid body", domain.ErrInvalidInput)
		}
		u, err := h.svc.Create(r.Context(), svcuser.CreateInput{
			Name:  req.Name,
			Email: req.Email,
		})
		if err != nil {
			return fmt.Errorf("svc.Create: %w", err)
		}
		httputil.JSON(w, http.StatusCreated, toUserResponse(u))
		return nil
	})
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	httputil.HandleFlow(w, h.log, func() error {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			return fmt.Errorf("%w: invalid id", domain.ErrInvalidInput)
		}
		var req dto.UpdateUserRequest
		if err := httputil.Decode(r, &req); err != nil {
			return fmt.Errorf("%w: invalid body", domain.ErrInvalidInput)
		}
		u, err := h.svc.Update(r.Context(), id, svcuser.UpdateInput{
			Name:  req.Name,
			Email: req.Email,
		})
		if err != nil {
			return fmt.Errorf("svc.Update: %w", err)
		}
		httputil.JSON(w, http.StatusOK, toUserResponse(u))
		return nil
	})
}

func (h *Handler) upsert(w http.ResponseWriter, r *http.Request) {
	httputil.HandleFlow(w, h.log, func() error {
		var req dto.UpsertUserRequest
		if err := httputil.Decode(r, &req); err != nil {
			return fmt.Errorf("%w: invalid body", domain.ErrInvalidInput)
		}
		u, err := h.svc.Upsert(r.Context(), svcuser.CreateInput{
			Name:  req.Name,
			Email: req.Email,
		})
		if err != nil {
			return fmt.Errorf("svc.Upsert: %w", err)
		}
		httputil.JSON(w, http.StatusOK, toUserResponse(u))
		return nil
	})
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	httputil.HandleFlow(w, h.log, func() error {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			return fmt.Errorf("%w: invalid id", domain.ErrInvalidInput)
		}
		if err := h.svc.Delete(r.Context(), id); err != nil {
			return fmt.Errorf("svc.Delete: %w", err)
		}
		w.WriteHeader(http.StatusNoContent)
		return nil
	})
}
