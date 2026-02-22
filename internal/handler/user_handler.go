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

// GetAll godoc
// @Summary Get all users
// @Tags users
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} modules.PaginatedUsers
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /users [get]
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	offset, err := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)
	if err != nil || offset < 0 {
		offset = 0
	}

	result, err := h.uc.GetAll(r.Context(), limit, offset)
	if err != nil {
		errorResponse(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// GetByID godoc
// @Summary Get user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} modules.User
// @Failure 404 {object} map[string]string
// @Security ApiKeyAuth
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
		return
	}

	user, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		errorResponse(w, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

// Create godoc
// @Summary Create user
// @Tags users
// @Accept json
// @Produce json
// @Param user body modules.UserInput true "User"
// @Success 201 {object} map[string]int64
// @Failure 400 {object} map[string]string
// @Security ApiKeyAuth
// @Router /users [post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input modules.UserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}

	id, err := h.uc.Create(r.Context(), &modules.User{Name: input.Name, Email: input.Email})
	if err != nil {
		errorResponse(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}

// Update godoc
// @Summary Update user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body modules.UserInput true "User"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security ApiKeyAuth
// @Router /users/{id} [put]
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
		return
	}

	var input modules.UserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}

	if err := h.uc.Update(r.Context(), &modules.User{ID: id, Name: input.Name, Email: input.Email}); err != nil {
		errorResponse(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "user updated"})
}

// Delete godoc
// @Summary Soft delete user
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Security ApiKeyAuth
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
		return
	}

	if err := h.uc.Delete(r.Context(), id); err != nil {
		errorResponse(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// CreateWithAudit godoc
// @Summary Create user with audit log
// @Tags users
// @Accept json
// @Produce json
// @Param user body modules.UserInput true "User"
// @Success 201 {object} map[string]int64
// @Failure 400 {object} map[string]string
// @Security ApiKeyAuth
// @Router /users/audit [post]
func (h *UserHandler) CreateWithAudit(w http.ResponseWriter, r *http.Request) {
	var input modules.UserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}

	id, err := h.uc.CreateUserWithAudit(r.Context(), &modules.User{Name: input.Name, Email: input.Email})
	if err != nil {
		errorResponse(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}
