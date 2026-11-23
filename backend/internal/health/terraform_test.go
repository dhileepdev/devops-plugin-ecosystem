package health

import (
	"testing"

	"plugins.dev/backend/internal/github"
)

func TestComputeHealth_Green(t *testing.T) {
	issues := []github.Issue{
		{Number: 1, Title: "minor typo", State: "open", Labels: []github.Label{{Name: "documentation"}}},
	}
	h := ComputeHealth("test-provider", issues)
	if h.Status != StatusGreen {
		t.Fatalf("expected green, got %s", h.Status)
	}
	if h.OpenIssues != 1 {
		t.Fatalf("expected 1 open issue, got %d", h.OpenIssues)
	}
}

func TestComputeHealth_Red(t *testing.T) {
	var issues []github.Issue
	for i := 0; i < 9; i++ {
		issues = append(issues, github.Issue{Number: i + 1, Title: "panic on apply", State: "open", Labels: []github.Label{{Name: "bug"}}})
	}
	h := ComputeHealth("test-provider", issues)
	if h.Status != StatusRed {
		t.Fatalf("expected red, got %s", h.Status)
	}
	if h.OpenCriticalBugs != 9 {
		t.Fatalf("expected 9 critical, got %d", h.OpenCriticalBugs)
	}
}

func TestComputeHealth_PerfRegression(t *testing.T) {
	issues := []github.Issue{
		{Number: 5, Title: "O(n²) complexity in nested objects", State: "open"},
	}
	h := ComputeHealth("test-provider", issues)
	if len(h.PerformanceRegressions) != 1 {
		t.Fatalf("expected 1 perf regression, got %d", len(h.PerformanceRegressions))
	}
}
