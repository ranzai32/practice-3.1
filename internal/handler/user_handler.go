package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"practice4/practice-4/internal/usecase"
	"practice4/practice-4/pkg/apperrors"
	"practice4/practice-4/pkg/modules"
	"strconv"
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

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
		return
	}

	user, err := h.uc.GetByID(r.Context(), idInt)
	if err != nil {
		errorResponse(w, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user modules.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}

	id, err := h.uc.Create(r.Context(), &user)
	if err != nil {
		errorResponse(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
		return
	}

	var user modules.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	user.ID = idInt

	if err := h.uc.Update(r.Context(), &user); err != nil {
		errorResponse(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "user updated"})
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest,	 map[string]string{"error": "invalid user ID"})
		return
	}

	if err := h.uc.Delete(r.Context(), idInt); err != nil {
		errorResponse(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}	