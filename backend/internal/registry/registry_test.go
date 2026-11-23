package registry

import (
	"errors"
	"testing"

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
	} else if !errors.Is(err, ErrDuplicate) {
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
