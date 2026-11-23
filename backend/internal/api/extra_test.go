package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	r := NewRouter()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/health", http.NoBody)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestOperatorAnalyzeEndpoint(t *testing.T) {
	r := NewRouter()
	body := `{"manifest":"kind: Role\nmetadata:\n  name: r\nrules:\n- verbs: [get]\n  resources: [pods]\n","target_version":"v1.0.0"}`
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/operator/analyze", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var out map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	if _, ok := out["upgrade_path"]; !ok {
		t.Fatal("missing upgrade_path")
	}
}

func TestNotifyIngestEndpoint(t *testing.T) {
	r := NewRouter()
	body := `{"id":"x","event_type":"pr.merge","source":"gh","severity":"info","message":"merged","timestamp":"2026-11-23T10:00:00Z"}`
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/notifications", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", rec.Code)
	}
	// stats
	req2 := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/notifications/stats", http.NoBody)
	rec2 := httptest.NewRecorder()
	r.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec2.Code)
	}
}
