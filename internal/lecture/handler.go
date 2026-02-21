package lecture

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreateLectureRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request", zap.Error(err))
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	lecture, err := h.service.Create(r.Context(), userID, req)
	if err != nil {
		h.logger.Error("failed to create lecture", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, "failed to create lecture")
		return
	}

	respondWithJSON(w, http.StatusCreated, lecture)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid lecture id")
		return
	}

	lecture, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.logger.Error("failed to get lecture", zap.Error(err))
		respondWithError(w, http.StatusNotFound, "lecture not found")
		return
	}

	respondWithJSON(w, http.StatusOK, lecture)
}

func (h *Handler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	lectures, err := h.service.GetByUserID(r.Context(), userID)
	if err != nil {
		h.logger.Error("failed to get lectures", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, "failed to get lectures")
		return
	}

	respondWithJSON(w, http.StatusOK, lectures)
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
