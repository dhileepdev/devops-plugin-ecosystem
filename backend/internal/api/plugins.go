package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"plugins.dev/backend/internal/models"
	"plugins.dev/backend/internal/registry"
)

type PluginHandler struct {
	reg *registry.Registry
}

func NewPluginHandler(reg *registry.Registry) *PluginHandler {
	return &PluginHandler{reg: reg}
}

func (h *PluginHandler) List(w http.ResponseWriter, r *http.Request) {
	platform := r.URL.Query().Get("platform")
	var plugins []models.Plugin
	if platform != "" {
		plugins = h.reg.ByPlatform(platform)
	} else {
		plugins = h.reg.List()
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(plugins)
}

func (h *PluginHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	p, err := h.reg.Get(id)
	if err != nil {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(p)
}
