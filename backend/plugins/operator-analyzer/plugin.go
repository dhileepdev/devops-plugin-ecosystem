// Package operatoranalyzer implements the Operator Lifecycle & RBAC Analyzer plugin.
package operatoranalyzer

import (
	"context"

	"plugins.dev/backend/internal/operator"
	"plugins.dev/backend/pkg/plugin"
)

const (
	PluginName    = "operator-lifecycle-rbac"
	PluginVersion = "0.1.0"
)

// Plugin wraps the operator analyzer as a standalone plugin.
type Plugin struct{}

// New creates a new operator analyzer plugin.
func New() *Plugin { return &Plugin{} }

func (p *Plugin) Name() string    { return PluginName }
func (p *Plugin) Version() string { return PluginVersion }
func (p *Plugin) Description() string {
	return "Scans Kubernetes operator manifests for upgrade requirements, RBAC gaps, and port-forward configuration. Generates an upgrade path."
}
func (p *Plugin) Category() string { return "analyzer" }

func (p *Plugin) Init(cfg plugin.Config) error { return nil }

func (p *Plugin) HealthCheck(ctx context.Context) (plugin.HealthStatus, error) {
	return plugin.HealthStatus{Healthy: true, Message: "stateless analyzer ready"}, nil
}

func (p *Plugin) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	manifest, _ := input["manifest"].(string)
	target, _ := input["target_version"].(string)
	if target == "" {
		target = "latest"
	}
	scan := operator.ScanManifest(manifest)
	up := operator.GenerateUpgradePath(scan, target)
	rbac := operator.AnalyzeRBAC(manifest)
	return map[string]interface{}{
		"scan":         scan,
		"upgrade_path": up,
		"rbac":         rbac,
	}, nil
}
