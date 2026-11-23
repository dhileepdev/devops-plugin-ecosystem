package notify

import (
	"testing"
	"time"
)

func TestDedup_DeduplicatesWithinWindow(t *testing.T) {
	eng := NewEngine(5 * time.Minute)
	n := Notification{ID: "n1", EventType: "build.failure", Source: "ci", Severity: SeverityCritical, Message: "fail", Timestamp: time.Now()}
	r1 := eng.Ingest(n)
	if r1.Deduplicated {
		t.Fatal("first should not be deduped")
	}
	n2 := n
	n2.ID = "n2"
	r2 := eng.Ingest(n2)
	if !r2.Deduplicated {
		t.Fatal("second should be deduped")
	}
}

func TestDedup_RoutesToChannels(t *testing.T) {
	eng := NewEngine(5 * time.Minute)
	n := Notification{ID: "n1", EventType: "deploy", Source: "cd", Severity: SeverityInfo, Message: "deployed", Timestamp: time.Now()}
	r := eng.Ingest(n)
	if len(r.RoutedTo) == 0 {
		t.Fatal("expected routing")
	}
}

func TestDedup_Stats(t *testing.T) {
	eng := NewEngine(5 * time.Minute)
	now := time.Now()
	eng.Ingest(Notification{ID: "a", EventType: "pr.merge", Source: "gh", Timestamp: now})
	eng.Ingest(Notification{ID: "b", EventType: "pr.merge", Source: "gh", Timestamp: now})
	stats := eng.Stats()
	if stats["total"] != 2 {
		t.Fatalf("expected total 2, got %d", stats["total"])
	}
	if stats["deduplicated"] != 1 {
		t.Fatalf("expected 1 deduped, got %d", stats["deduplicated"])
	}
}
