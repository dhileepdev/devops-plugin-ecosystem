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
