package api

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

// ManifestHandler serves plugin manifest files for individual discovery.
type ManifestHandler struct {
	dir string
}

func NewManifestHandler(dir string) *ManifestHandler {
	if dir == "" {
		dir = "manifests"
	}
	candidates := []string{dir, "../manifests", "../../manifests", "../../../manifests"}
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			dir = c
			break
		}
	}
	return &ManifestHandler{dir: dir}
}

func (h *ManifestHandler) List(w http.ResponseWriter, r *http.Request) {
	entries, err := os.ReadDir(h.dir)
	if err != nil {
		http.Error(w, `{"error":"manifests unavailable"}`, http.StatusInternalServerError)
		return
	}
	var manifests []map[string]interface{}
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(h.dir, e.Name()))
		if err != nil {
			continue
		}
		var m map[string]interface{}
		if json.Unmarshal(data, &m) == nil {
			manifests = append(manifests, m)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(manifests)
}

func (h *ManifestHandler) Get(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	path := filepath.Join(h.dir, name+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		http.Error(w, `{"error":"manifest not found"}`, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}
