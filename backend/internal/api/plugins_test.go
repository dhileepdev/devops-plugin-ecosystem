package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"plugins.dev/backend/internal/models"
	"plugins.dev/backend/internal/registry"
)

func TestPluginsList(t *testing.T) {
	reg := registry.New()
	_ = reg.Create(models.Plugin{ID: "p1", Name: "provider-aws", Platform: "terraform"})
	r := NewRouterWithRegistry(reg)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/plugins", http.NoBody)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var plugins []models.Plugin
	if err := json.NewDecoder(rec.Body).Decode(&plugins); err != nil {
		t.Fatalf("bad json: %v", err)
	}
	if len(plugins) != 1 {
		t.Fatalf("expected 1 plugin, got %d", len(plugins))
	}
}
