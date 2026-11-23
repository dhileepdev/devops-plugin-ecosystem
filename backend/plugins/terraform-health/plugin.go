// Package terraformhealth implements the Terraform Provider Health plugin.
package terraformhealth

import (
	"context"

	"plugins.dev/backend/internal/github"
	"plugins.dev/backend/internal/health"
	"plugins.dev/backend/pkg/plugin"
)

const (
	PluginName    = "terraform-provider-health"
	PluginVersion = "0.1.0"
)

// Plugin wraps the Terraform health monitor as a standalone plugin.
type Plugin struct {
	gh     *github.Client
	report *health.TerraformHealth
}

// New creates a new Terraform health plugin.
func New() *Plugin {
	return &Plugin{}
}

func (p *Plugin) Name() string    { return PluginName }
func (p *Plugin) Version() string { return PluginVersion }
func (p *Plugin) Description() string {
	return "Surfaces bugs, panics, and performance regressions from hashicorp/terraform-plugin-sdk GitHub Issues and computes a provider health score."
}
func (p *Plugin) Category() string { return "health" }

func (p *Plugin) Init(cfg plugin.Config) error {
	p.gh = github.NewClient(cfg.GitHubToken)
	return nil
}

func (p *Plugin) HealthCheck(ctx context.Context) (plugin.HealthStatus, error) {
	if p.gh == nil {
		return plugin.HealthStatus{Healthy: false, Message: "not initialized"}, nil
	}
	_, err := p.gh.FetchTerraformSDKIssues()
	if err != nil {
		return plugin.HealthStatus{Healthy: false, Message: err.Error()}, nil
	}
	return plugin.HealthStatus{Healthy: true, Message: "GitHub API reachable"}, nil
}

func (p *Plugin) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	provider, _ := input["provider"].(string)
	if provider == "" {
		provider = "hashicorp/terraform-plugin-sdk"
	}
	issues, err := p.gh.FetchTerraformSDKIssues()
	if err != nil {
		return nil, err
	}
	report := health.ComputeHealth(provider, issues)
	p.report = &report
	return map[string]interface{}{
		"provider":                report.Provider,
		"status":                  string(report.Status),
		"open_critical_bugs":      report.OpenCriticalBugs,
		"performance_regressions": report.PerformanceRegressions,
		"plan_stability":          report.PlanStability,
		"open_issues":             report.OpenIssues,
	}, nil
}

// Report returns the last computed health report.
func (p *Plugin) Report() *health.TerraformHealth {
	return p.report
}
