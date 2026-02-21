package user

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type Handler struct {
	service Service
	logger  *zap.Logger
}

func NewHandler(service Service, logger *zap.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request", zap.Error(err))
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	user, err := h.service.Register(r.Context(), req)
	if err != nil {
		h.logger.Error("failed to register user", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, "failed to register user")
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request", zap.Error(err))
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	response, err := h.service.Login(r.Context(), req)
	if err != nil {
		h.logger.Error("failed to login", zap.Error(err))
		respondWithError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
