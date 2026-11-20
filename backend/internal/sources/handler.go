package sources

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	repository *Repository
}

func NewHandler(repository *Repository) *Handler {
	return &Handler{repository: repository}
}

type createSourceRequest struct {
	Name   string          `json:"name"`
	Type   string          `json:"type"`
	Config json.RawMessage `json:"config"`
}

type errorResponse struct {
	Error string `json:"error"`
}

var uuidPattern = regexp.MustCompile(
	`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
)

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/api/v1/sources", h.create)
	r.Get("/api/v1/sources", h.list)
	r.Get("/api/v1/sources/{sourceID}", h.get)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var request createSourceRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON request")
		return
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		writeError(w, http.StatusBadRequest, "request must contain one JSON object")
		return
	}

	request.Name = strings.TrimSpace(request.Name)
	request.Type = strings.TrimSpace(request.Type)

	if request.Name == "" || request.Type == "" {
		writeError(w, http.StatusBadRequest, "name and type are required")
		return
	}

	if len(request.Config) == 0 {
		request.Config = json.RawMessage(`{}`)
	}

	var configObject map[string]json.RawMessage
	if err := json.Unmarshal(request.Config, &configObject); err != nil || configObject == nil {
		writeError(w, http.StatusBadRequest, "config must be a JSON object")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	source, err := h.repository.Create(
		ctx,
		request.Name,
		request.Type,
		request.Config,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create source")
		return
	}

	writeJSON(w, http.StatusCreated, source)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	sources, err := h.repository.List(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list sources")
		return
	}

	writeJSON(w, http.StatusOK, sources)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "sourceID")

	if !uuidPattern.MatchString(id) {
		writeError(w, http.StatusBadRequest, "invalid source ID")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	source, err := h.repository.GetByID(ctx, id)
	if errors.Is(err, ErrNotFound) {
		writeError(w, http.StatusNotFound, "source not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to retrieve source")
		return
	}

	writeJSON(w, http.StatusOK, source)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}
