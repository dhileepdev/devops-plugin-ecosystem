package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"plugins.dev/backend/internal/github"
	"plugins.dev/backend/internal/notify"
	"plugins.dev/backend/internal/registry"
)

func NewRouter() http.Handler {
	return NewRouterWithRegistry(registry.New())
}

// NewRouterWithRegistry preserves the existing signature used by tests.
func NewRouterWithRegistry(reg *registry.Registry) http.Handler {
	return NewRouterWithDeps(reg, github.NewClient(""), notify.NewEngine(5*time.Minute))
}

func NewRouterWithDeps(reg *registry.Registry, gh *github.Client, eng *notify.Engine) http.Handler {
	ph := NewPluginHandler(reg)
	hh := NewHealthHandler(gh)
	nh := NewNotificationHandler(eng)
	mh := NewManifestHandler("../../manifests")

	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Route("/health", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"status":"ok"}`))
			})
			r.Get("/terraform", hh.Terraform)
			r.Get("/operator", hh.Operator)
		})
		r.Route("/plugins", func(r chi.Router) {
			r.Get("/", ph.List)
			r.Get("/{id}", ph.Get)
		})
		r.Route("/manifests", func(r chi.Router) {
			r.Get("/", mh.List)
			r.Get("/{name}", mh.Get)
		})
		r.Route("/notifications", func(r chi.Router) {
			r.Post("/", nh.Ingest)
			r.Get("/", nh.List)
			r.Get("/stats", nh.Stats)
		})
		r.Post("/operator/analyze", hh.OperatorAnalyze)
	})
	return r
}
