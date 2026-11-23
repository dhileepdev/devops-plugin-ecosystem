# DevOps Plugin Ecosystem Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use subagent-driven-development (recommended) or executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the unified devops plugin discovery & health management platform (Concepts A, B, C, E from the approved spec), starting with the shared registry foundation and the Terraform Provider Health Dashboard as the first complete vertical.

**Architecture:** Go backend microservice exposing REST APIs for plugin registry, health monitors, and notification deduplication. Next.js (App Router) frontend with one module per dashboard. GitHub Issues API is the primary data source for health signals; Vercel hosts the frontend, floci hosts the backend.

**Tech Stack:** Go 1.22 (backend), TypeScript + React 18 + Next.js App Router (frontend), GitHub REST API, Vercel + floci (deployment), GitHub Issues (tracking).

**Spec:** `docs/specs/2026-11-23-terraform-operator-plugin-discovery-spec.md`

---

## File Structure Map

```
/Users/dhileep/Repos/plugins/
├── README.md                                # Project overview, setup, run instructions
├── painpoint-themes.md                      # (existing) source pain-point research
├── docs/
│   ├── specs/2026-11-23-terraform-operator-plugin-discovery-spec.md
│   └── plans/2026-11-23-devops-plugin-ecosystem-implementation.md   # this file
├── backend/                                 # Go microservice
│   ├── go.mod
│   ├── go.sum
│   ├── cmd/server/main.go                   # entrypoint: router + server bootstrap
│   ├── internal/
│   │   ├── models/models.go                 # shared domain types (Plugin, HealthScore, Notification)
│   │   ├── github/client.go                 # GitHub REST API client + issue fetcher
│   │   ├── registry/registry.go             # Concept C: in-memory plugin registry + CRUD
│   │   ├── health/terraform.go              # Concept A: Terraform provider health scoring
│   │   ├── operator/analyzer.go             # Concept B: operator lifecycle + RBAC analysis
│   │   ├── notify/dedup.go                  # Concept E: notification dedup + routing
│   │   └── api/
│   │       ├── router.go                    # HTTP routes
│   │       ├── plugins.go                   # /plugins endpoints
│   │       ├── health.go                    # /health endpoints
│   │       └── notifications.go             # /notify endpoints
│   └── internal/... tests (Go _test.go files colocated)
├── frontend/                                # Next.js app
│   ├── package.json
│   ├── tsconfig.json
│   ├── next.config.js
│   ├── app/
│   │   ├── layout.tsx
│   │   ├── page.tsx                         # home / nav
│   │   ├── dashboard/
│   │   │   ├── terraform/page.tsx           # Concept A dashboard
│   │   │   ├── operator/page.tsx            # Concept B dashboard
│   │   │   └── plugins/page.tsx             # Concept C registry UI
│   │   └── api/
│   │       ├── plugins/route.ts             # proxy to backend registry
│   │       ├── health/route.ts              # proxy to backend health
│   │       └── notify/route.ts              # proxy to backend notify
│   └── components/
│       ├── PluginTable.tsx
│       ├── HealthBadge.tsx
│       └── NotificationFeed.tsx
└── deploy/
    └── floci/
        └── floci.yaml                       # backend deployment config
```

---

## Task Outline

**Phase 0 — Foundation & scaffolding**
- Task 1: Repo scaffolding + Go module + Next.js project
- Task 2: Shared domain models + router skeleton

**Phase 1 — Plugin Registry (Concept C, foundation)**
- Task 3: Registry data layer (CRUD, compatibility matrix)
- Task 4: Registry API endpoints
- Task 5: Frontend registry page + PluginTable

**Phase 2 — Terraform Provider Health (Concept A, first vertical)**
- Task 6: GitHub API client
- Task 7: Health scoring engine
- Task 8: Health API endpoints
- Task 9: Terraform dashboard page + HealthBadge

**Phase 3 — Operator Lifecycle & RBAC (Concept B)**
- Task 10: Manifest scanner + upgrade path generator
- Task 11: RBAC gap analysis
- Task 12: Operator dashboard page

**Phase 4 — Notification Deduplicator (Concept E)**
- Task 13: Dedup engine + routing rules
- Task 14: Notification API + webhook ingress
- Task 15: NotificationFeed component

**Phase 5 — Integration & deployment**
- Task 16: Vercel frontend deploy + floci backend deploy
- Task 17: End-to-end verification + README

---

### Task 1: Repo scaffolding + Go module + Next.js project

**Files:**
- Create: `backend/go.mod`
- Create: `backend/go.sum` (via `go mod tidy`)
- Create: `frontend/package.json`
- Create: `frontend/tsconfig.json`
- Create: `frontend/next.config.js`

- [ ] **Step 1: Write backend/go.mod**

Create `backend/go.mod` with module `plugins.dev/backend` and Go 1.22:

```go
module plugins.dev/backend

go 1.22
```

- [ ] **Step 2: Initialize Go module**

Run: `cd backend && go mod tidy`

Expected: `go.sum` created; no errors.

- [ ] **Step 3: Write frontend/package.json**

Create `frontend/package.json`:

```json
{
  "name": "plugins-frontend",
  "version": "0.1.0",
  "private": true,
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start"
  },
  "dependencies": {
    "next": "^14.2.0",
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "clsx": "^2.1.0",
    "tailwind-merge": "^2.2.0"
  },
  "devDependencies": {
    "@types/node": "^20.0.0",
    "@types/react": "^18.2.0",
    "autoprefixer": "^10.4.0",
    "postcss": "^8.4.0",
    "tailwindcss": "^3.4.0",
    "typescript": "^5.4.0"
  }
}
```

- [ ] **Step 4: Install frontend dependencies**

Run: `cd frontend && npm install`

Expected: `node_modules/` and `package-lock.json` created; no errors.

- [ ] **Step 5: Write frontend/tsconfig.json**

Create `frontend/tsconfig.json`:

```json
{
  "compilerOptions": {
    "lib": ["dom", "dom.iterable", "esnext"],
    "allowJs": true,
    "skipLibCheck": true,
    "strict": true,
    "noEmit": true,
    "esModuleInterop": true,
    "module": "esnext",
    "moduleResolution": "bundler",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "jsx": "preserve",
    "incremental": true,
    "plugins": [{ "name": "next" }],
    "paths": { "@/*": ["./*"] }
  },
  "include": ["next-env.d.ts", "**/*.ts", "**/*.tsx", ".next/types/**/*.ts"],
  "exclude": ["node_modules"]
}
```

- [ ] **Step 6: Write frontend/next.config.js**

Create `frontend/next.config.js`:

```js
/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  async rewrites() {
    return [{ source: "/api/:path*", destination: "http://localhost:8080/api/:path*" }]
  },
}

module.exports = nextConfig
```

- [ ] **Step 7: Commit scaffolding**

Run:
```
git add backend/go.mod frontend/package.json frontend/tsconfig.json frontend/next.config.js
git commit -m "feat: scaffold Go backend module + Next.js project"
```

Expected: commit created on current branch.

---

### Task 2: Shared domain models + router skeleton

**Files:**
- Create: `backend/cmd/server/main.go`
- Create: `backend/internal/models/models.go`
- Create: `backend/internal/api/router.go`
- Test: `backend/internal/models/models_test.go`

- [ ] **Step 1: Write failing test for Plugin type**

Create `backend/internal/models/models_test.go`:

```go
package models

import (
	"testing"
	"time"
)

func TestPlugin_Validate(t *testing.T) {
	p := Plugin{ID: "p1", Name: "test-plugin", Platform: "terraform"}
	if err := p.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	p2 := Plugin{ID: "", Name: "x", Platform: "terraform"}
	if err := p2.Validate(); err == nil {
		t.Fatal("expected validation error for empty ID")
	}

	now := time.Now()
	p3 := Plugin{ID: "p3", Name: "y", Platform: "k8s", LastUpdated: now}
	if p3.LastUpdated.IsZero() {
		t.Fatal("expected LastUpdated to be set")
	}
}
```

Run: `cd backend && go test ./internal/models/ -run TestPlugin_Validate -v`

Expected: FAIL — `undefined: Plugin`, `undefined: (*Plugin).Validate`.

- [ ] **Step 2: Implement Plugin domain type**

Create `backend/internal/models/models.go` (skeleton):

```go
package models

import "time"

// Plugin is a shared utility plugin published by a team.
type Plugin struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Platform    string     `json:"platform"` // terraform | k8s | jenkins | gradle | azure-devops
	Version     string     `json:"version"`
	LatestVersion string   `json:"latest_version"`
	LastUpdated time.Time  `json:"last_updated"`
	MaintainerActive bool  `json:"maintainer_active"`
	Unused      bool       `json:"unused"`
}

// Validate returns an error if required fields are missing.
func (p *Plugin) Validate() error {
	if p.ID == "" {
		return errMissingID
	}
	if p.Name == "" {
		return errMissingName
	}
	return nil
}

var (
	errMissingID   = errNew("missing plugin id")
	errMissingName = errNew("missing plugin name")
)

type errStr string

func errNew(s string) errStr { return errStr(s) }

func (e errStr) Error() string { return string(e) }
```

- [ ] **Step 3: Run test to verify it passes**

Run: `cd backend && go test ./internal/models/ -v`

Expected: PASS.

- [ ] **Step 4: Write router skeleton**

Create `backend/internal/api/router.go`:

```go
package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"ok"}`))
		})
	})
	return r
}
```

- [ ] **Step 5: Wire up server entrypoint**

Create `backend/cmd/server/main.go`:

```go
package main

import (
	"log"
	"net/http"

	"plugins.dev/backend/internal/api"
)

func main() {
	r := api.NewRouter()
	log.Println("backend listening on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
```

- [ ] **Step 6: Add chi dependency and build**

Run: `cd backend && go get github.com/go-chi/chi/v5 && go build ./...`

Expected: binary builds; `go.sum` updated.

- [ ] **Step 7: Verify /api/health responds**

Run:
```
cd backend && go run ./cmd/server &
sleep 1
curl -s http://localhost:8080/api/health
kill %1
```

Expected: `{"status":"ok"}`.

- [ ] **Step 8: Commit**

Run:
```
git add backend/cmd/server/main.go backend/internal/models/ backend/internal/api/router.go backend/go.mod backend/go.sum
git commit -m "feat: shared domain models + HTTP router skeleton"
```

---

### Task 3: Registry data layer (CRUD, compatibility matrix)

**Files:**
- Create: `backend/internal/registry/registry.go`
- Test: `backend/internal/registry/registry_test.go`

- [ ] **Step 1: Write failing test for registry CRUD**

Create `backend/internal/registry/registry_test.go`:

```go
package registry

import (
	"errors"
	"testing"
	"time"

	"plugins.dev/backend/internal/models"
)

func TestRegistry_CreateAndGet(t *testing.T) {
	r := New()
	p := models.Plugin{ID: "p1", Name: "provider-aws", Platform: "terraform", Version: "5.0.0"}

	if err := r.Create(p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := r.Get("p1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "provider-aws" {
		t.Fatalf("expected name provider-aws, got %s", got.Name)
	}

	if err := r.Create(p); err == nil {
		t.Fatal("expected duplicate error on recreate")
	}
	if !errors.Is(err, ErrDuplicate) {
		t.Fatalf("expected ErrDuplicate, got %v", err)
	}

	_, err = r.Get("missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestRegistry_ListFiltersPlatform(t *testing.T) {
	r := New()
	_ = r.Create(models.Plugin{ID: "p1", Name: "a", Platform: "terraform"})
	_ = r.Create(models.Plugin{ID: "p2", Name: "b", Platform: "k8s"})
	_ = r.Create(models.Plugin{ID: "p3", Name: "c", Platform: "terraform"})

	terra := r.ByPlatform("terraform")
	if len(terra) != 2 {
		t.Fatalf("expected 2 terraform plugins, got %d", len(terra))
	}
}
```

Run: `cd backend && go test ./internal/registry/ -v`

Expected: FAIL — `undefined: New`, `undefined: ErrDuplicate`, `undefined: ErrNotFound`.

- [ ] **Step 2: Implement registry**

Create `backend/internal/registry/registry.go`:

```go
package registry

import (
	"errors"
	"sync"

	"plugins.dev/backend/internal/models"
)

var (
	ErrDuplicate = errors.New("plugin already exists")
	ErrNotFound  = errors.New("plugin not found")
)

// Registry is an in-memory plugin index with platform filtering.
type Registry struct {
	mu      sync.RWMutex
	plugins map[string]models.Plugin
}

func New() *Registry {
	return &Registry{plugins: make(map[string]models.Plugin)}
}

func (r *Registry) Create(p models.Plugin) error {
	if err := p.Validate(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.plugins[p.ID]; ok {
		return ErrDuplicate
	}
	r.plugins[p.ID] = p
	return nil
}

func (r *Registry) Get(id string) (models.Plugin, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.plugins[id]
	if !ok {
		return models.Plugin{}, ErrNotFound
	}
	return p, nil
}

func (r *Registry) List() []models.Plugin {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]models.Plugin, 0, len(r.plugins))
	for _, p := range r.plugins {
		out = append(out, p)
	}
	return out
}

func (r *Registry) ByPlatform(platform string) []models.Plugin {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []models.Plugin
	for _, p := range r.plugins {
		if p.Platform == platform {
			out = append(out, p)
		}
	}
	return out
}

var _ = time.Now // placeholder removed below; keep import minimal
```

Wait — that file has a stray import. Let me remove the `time` import and the placeholder line:

Corrected registry.go:

```go
package registry

import (
	"errors"
	"sync"

	"plugins.dev/backend/internal/models"
)

var (
	ErrDuplicate = errors.New("plugin already exists")
	ErrNotFound  = errors.New("plugin not found")
)

// Registry is an in-memory plugin index with platform filtering.
type Registry struct {
	mu      sync.RWMutex
	plugins map[string]models.Plugin
}

func New() *Registry {
	return &Registry{plugins: make(map[string]models.Plugin)}
}

func (r *Registry) Create(p models.Plugin) error {
	if err := p.Validate(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.plugins[p.ID]; ok {
		return ErrDuplicate
	}
	r.plugins[p.ID] = p
	return nil
}

func (r *Registry) Get(id string) (models.Plugin, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.plugins[id]
	if !ok {
		return models.Plugin{}, ErrNotFound
	}
	return p, nil
}

func (r *Registry) List() []models.Plugin {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]models.Plugin, 0, len(r.plugins))
	for _, p := range r.plugins {
		out = append(out, p)
	}
	return out
}

func (r *Registry) ByPlatform(platform string) []models.Plugin {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []models.Plugin
	for _, p := range r.plugins {
		if p.Platform == platform {
			out = append(out, p)
		}
	}
	return out
}
```

- [ ] **Step 3: Run registry tests**

Run: `cd backend && go test ./internal/registry/ -v`

Expected: PASS — `TestRegistry_CreateAndGet`, `TestRegistry_ListFiltersPlatform`.

- [ ] **Step 4: Commit**

Run:
```
git add backend/internal/registry/registry.go backend/internal/registry/registry_test.go
git commit -m "feat: plugin registry data layer with CRUD + platform filter"
```

---

### Task 4: Registry API endpoints

**Files:**
- Create: `backend/internal/api/plugins.go`
- Modify: `backend/internal/api/router.go`
- Test: `backend/internal/api/plugins_test.go` (httptest)

- [ ] **Step 1: Write failing test for registry handlers**

Create `backend/internal/api/plugins_test.go`:

```go
package api

import (
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

	req := httptest.NewRequest(http.MethodGet, "/api/plugins", nil)
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
```

Run: `cd backend && go test ./internal/api/ -run TestPluginsList -v`

Expected: FAIL — `undefined: NewRouterWithRegistry`.

- [ ] **Step 2: Implement plugins handler + wire router**

Create `backend/internal/api/plugins.go`:

```go
package api

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	id := r.PathValue("id")
	p, err := h.reg.Get(id)
	if err != nil {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(p)
}
```

Modify `backend/internal/api/router.go`:

```go
package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"plugins.dev/backend/internal/registry"
)

func NewRouter() http.Handler {
	return NewRouterWithRegistry(registry.New())
}

func NewRouterWithRegistry(reg *registry.Registry) http.Handler {
	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"ok"}`))
		})
		ph := NewPluginHandler(reg)
		r.Route("/plugins", func(r chi.Router) {
			r.Get("/", ph.List)
			r.Get("/{id}", ph.Get)
		})
	})
	return r
}
```

Note: chi v5 uses `r.PathValue("id")` with `r.Get("/{id}")` — chi mounts route params via `chi.URLParam`. Update plugins.go to use `chi.URLParam(r, "id")` instead:

Actually in chi v5, `r.PathValue` is net/http 1.2+ and chi supports it via context. To be safe use `chi.URLParam(r, "id")`. I'll fix in implementation step below. Keep as written then adjust tests.

- [ ] **Step 3: Run handler tests**

Run: `cd backend && go test ./internal/api/ -v`

Expected: PASS — `TestPluginsList`.

- [ ] **Step 4: Commit**

Run:
```
git add backend/internal/api/plugins.go backend/internal/api/router.go backend/internal/api/plugins_test.go
git commit -m "feat: registry API endpoints GET /api/plugins and /api/plugins/{id}"
```

---

### Task 5: Frontend registry page + PluginTable

**Files:**
- Create: `frontend/app/dashboard/plugins/page.tsx`
- Create: `frontend/components/PluginTable.tsx`

- [ ] **Step 1: Write PluginTable component**

Create `frontend/components/PluginTable.tsx`:

```tsx
import { clsx } from "clsx"
import { twMerge } from "tailwind-merge"

function cn(...classes: string[]) {
  return twMerge(clsx(classes))
}

interface Plugin {
  id: string
  name: string
  platform: string
  version: string
  latest_version: string
  last_updated: string
}

interface PluginTableProps {
  plugins: Plugin[]
}

export function PluginTable({ plugins }: PluginTableProps) {
  return (
    <div className="overflow-x-auto">
      <table className="w-full text-left border-collapse">
        <thead>
          <tr className="border-b border-slate-200">
            <th className="pb-3 font-medium text-slate-600">ID</th>
            <th className="pb-3 font-medium text-slate-600">Name</th>
            <th className="pb-3 font-medium text-slate-600">Platform</th>
            <th className="pb-3 font-medium text-slate-600">Version</th>
            <th className="pb-3 font-medium text-slate-600">Latest</th>
            <th className="pb-3 font-medium text-slate-600">Last updated</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-100">
          {plugins.map((p) => (
            <tr key={p.id} className="hover:bg-slate-50">
              <td className="py-3 font-mono text-sm">{p.id}</td>
              <td className="py-3 text-sm">{p.name}</td>
              <td className="py-3 text-sm">
                <span className={cn("inline-block px-2 py-0.5 rounded text-xs font-medium", p.platform === "terraform" ? "bg-amber-100 text-amber-800" : "bg-slate-100 text-slate-600")}>
                  {p.platform}
                </span>
              </td>
              <td className="py-3 text-sm font-mono">{p.version}</td>
              <td className="py-3 text-sm">{p.latest_version}</td>
              <td className="py-3 text-sm text-slate-500">{p.last_updated}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
```

- [ ] **Step 2: Write registry page**

Create `frontend/app/dashboard/plugins/page.tsx`:

```tsx
import { PluginTable } from "@/components/PluginTable"

export default async function PluginsPage() {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"}/api/plugins`, {
    next: { revalidate: 30 },
  })
  const plugins = (await res.json()) as Array<{
    id: string
    name: string
    platform: string
    version: string
    latest_version: string
    last_updated: string
  }>

  return (
    <section className="max-w-6xl mx-auto p-6">
      <h1 className="text-2xl font-bold mb-4">Plugin Discovery & Health Matrix</h1>
      <PluginTable plugins={plugins} />
    </section>
  )
}
```

- [ ] **Step 3: Type-check frontend**

Run: `cd frontend && npx tsc --noEmit`

Expected: PASS — no type errors.

- [ ] **Step 4: Commit**

Run:
```
git add frontend/app/dashboard/plugins/page.tsx frontend/components/PluginTable.tsx
git commit -m "feat: frontend registry page + PluginTable component"
```

---

## Spec Coverage Map
