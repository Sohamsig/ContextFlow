package records

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

type upsertRecordRequest struct {
	ExternalID  string          `json:"external_id"`
	Content     string          `json:"content"`
	ContentType string          `json:"content_type"`
	Metadata    json.RawMessage `json:"metadata"`
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/api/v1/sources/{sourceID}/records", h.upsert)
	r.Get("/api/v1/sources/{sourceID}/records", h.get)
}

func (h *Handler) upsert(w http.ResponseWriter, r *http.Request) {
	sourceID := chi.URLParam(r, "sourceID")

	if !isUUID(sourceID) {
		writeError(w, http.StatusBadRequest, "invalid source ID")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var request upsertRecordRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, "request must contain one JSON object")
		return
	}

	request.ExternalID = strings.TrimSpace(request.ExternalID)
	request.ContentType = strings.TrimSpace(request.ContentType)

	if request.ExternalID == "" {
		writeError(w, http.StatusBadRequest, "external_id is required")
		return
	}

	if request.Content == "" {
		writeError(w, http.StatusBadRequest, "content is required")
		return
	}

	if request.ContentType == "" {
		request.ContentType = "text"
	}

	if len(request.Metadata) == 0 {
		request.Metadata = json.RawMessage(`{}`)
	}

	var metadataObject map[string]json.RawMessage

	if err := json.Unmarshal(request.Metadata, &metadataObject); err != nil || metadataObject == nil {
		writeError(w, http.StatusBadRequest, "metadata must be a JSON object")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	record, err := h.repository.Upsert(
		ctx,
		sourceID,
		request.ExternalID,
		request.Content,
		request.ContentType,
		request.Metadata,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to ingest context record")
		return
	}

	writeJSON(w, http.StatusOK, record)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	sourceID := chi.URLParam(r, "sourceID")

	if !isUUID(sourceID) {
		writeError(w, http.StatusBadRequest, "invalid source ID")
		return
	}

	externalID := strings.TrimSpace(r.URL.Query().Get("external_id"))

	if externalID == "" {
		writeError(w, http.StatusBadRequest, "external_id is required")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	record, err := h.repository.GetByExternalID(ctx, sourceID, externalID)

	if errors.Is(err, ErrNotFound) {
		writeError(w, http.StatusNotFound, "context record not found")
		return
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to retrieve context record")
		return
	}

	writeJSON(w, http.StatusOK, record)
}

func isUUID(value string) bool {
	const pattern = `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`

	return regexp.MustCompile(pattern).MatchString(value)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
