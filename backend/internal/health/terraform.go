package health

import (
	"fmt"
	"strings"

	"plugins.dev/backend/internal/github"
)

// Status represents a provider health level.
type Status string

const (
	StatusGreen  Status = "green"
	StatusYellow Status = "yellow"
	StatusRed    Status = "red"
)

// TerraformHealth is the health report for a Terraform provider/SDK.
type TerraformHealth struct {
	Provider               string         `json:"provider"`
	Status                 Status         `json:"status"`
	OpenCriticalBugs       int            `json:"open_critical_bugs"`
	PerformanceRegressions []string       `json:"performance_regressions"`
	PlanStability          float64        `json:"plan_stability"`
	KnownUnusedResources   int            `json:"known_unused_resources"`
	OpenIssues             int            `json:"open_issues"`
	Issues                 []github.Issue `json:"issues"`
}

// isCritical reports whether an issue is a critical bug.
//
//nolint:gocritic // passing by value is acceptable for this struct
func isCritical(i github.Issue) bool {
	for _, l := range i.Labels {
		n := strings.ToLower(l.Name)
		if n == "bug" || n == "critical" || n == "crash" || n == "panic" {
			return true
		}
	}
	lower := strings.ToLower(i.Title)
	return strings.Contains(lower, "panic") || strings.Contains(lower, "crash")
}

// isPerformanceRegression reports whether an issue describes a perf regression.
//
//nolint:gocritic // passing by value is acceptable for this struct
func isPerformanceRegression(i github.Issue) bool {
	lower := strings.ToLower(i.Title + " " + i.Body)
	keywords := []string{"o(n", "o(n²)", "performance", "regression", "slow", "complexity"}
	for _, k := range keywords {
		if strings.Contains(lower, k) {
			return true
		}
	}
	for _, l := range i.Labels {
		if strings.Contains(strings.ToLower(l.Name), "performance") {
			return true
		}
	}
	return false
}

// ComputeHealth builds a TerraformHealth report from a list of GitHub issues.
func ComputeHealth(provider string, issues []github.Issue) TerraformHealth {
	h := TerraformHealth{
		Provider: provider,
		Issues:   issues,
	}
	for _, i := range issues { //nolint:gocritic
		if i.State != "open" && i.State != "" {
			continue
		}
		h.OpenIssues++
		if isCritical(i) {
			h.OpenCriticalBugs++
		}
		if isPerformanceRegression(i) {
			h.PerformanceRegressions = append(h.PerformanceRegressions, fmtIssue(i))
		}
	}
	// plan stability: simple heuristic, 100% minus 5% per critical bug, floored at 0
	stability := 100.0 - float64(h.OpenCriticalBugs)*5.0
	if stability < 0 {
		stability = 0
	}
	h.PlanStability = stability

	switch {
	case h.OpenCriticalBugs >= 8 || len(h.PerformanceRegressions) >= 3:
		h.Status = StatusRed
	case h.OpenCriticalBugs >= 3 || len(h.PerformanceRegressions) >= 1:
		h.Status = StatusYellow
	default:
		h.Status = StatusGreen
	}
	return h
}

//nolint:gocritic // passing by value is acceptable for this struct
func fmtIssue(i github.Issue) string {
	return fmt.Sprintf("#%d %s", i.Number, i.Title)
}
