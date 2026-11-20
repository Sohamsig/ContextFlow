package records

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func newTestRouter() http.Handler {
	router := chi.NewRouter()
	handler := NewHandler(nil)
	handler.RegisterRoutes(router)
	return router
}

func TestUpsertRejectsInvalidRequests(t *testing.T) {
	tests := []struct {
		name       string
		sourceID   string
		body       string
		wantStatus int
	}{
		{
			name:       "invalid source ID",
			sourceID:   "not-a-uuid",
			body:       `{}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "malformed JSON",
			sourceID:   "97abb98c-b761-46e0-915f-204e13a55b83",
			body:       `{"external_id":`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "unknown field",
			sourceID:   "97abb98c-b761-46e0-915f-204e13a55b83",
			body:       `{"external_id":"note-1","content":"hello","unexpected":true}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing external ID",
			sourceID:   "97abb98c-b761-46e0-915f-204e13a55b83",
			body:       `{"content":"hello"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing content",
			sourceID:   "97abb98c-b761-46e0-915f-204e13a55b83",
			body:       `{"external_id":"note-1"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "metadata must be a JSON object",
			sourceID:   "97abb98c-b761-46e0-915f-204e13a55b83",
			body:       `{"external_id":"note-1","content":"hello","metadata":[]}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "metadata must be valid JSON",
			sourceID:   "97abb98c-b761-46e0-915f-204e13a55b83",
			body:       `{"external_id":"note-1","content":"hello","metadata":"invalid"}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(
				http.MethodPost,
				"/api/v1/sources/"+test.sourceID+"/records",
				strings.NewReader(test.body),
			)

			response := httptest.NewRecorder()
			newTestRouter().ServeHTTP(response, request)

			if response.Code != test.wantStatus {
				t.Fatalf(
					"expected status %d, got %d; body: %s",
					test.wantStatus,
					response.Code,
					response.Body.String(),
				)
			}
		})
	}
}
