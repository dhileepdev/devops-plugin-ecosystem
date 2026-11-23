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

//nolint:gocritic // passing by value is acceptable
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
