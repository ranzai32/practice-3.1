package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"practice4/practice-4/internal/usecase"
	"practice4/practice-4/pkg/apperrors"
)

type UserHandler struct {
	uc usecase.UserUsecase
}

func NewUserHandler(uc usecase.UserUsecase) *UserHandler {
	return &UserHandler{uc: uc}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func errorResponse(w http.ResponseWriter, err error) {
    switch {
    case errors.Is(err, apperrors.ErrNotFound):
        writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
    case errors.Is(err, apperrors.ErrValidation):
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
    default:
        writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
    }
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.uc.GetAll(r.Context())
	if err != nil {
		errorResponse(w, err)
		return
	}
	writeJSON(w, http.StatusOK, users)
}

