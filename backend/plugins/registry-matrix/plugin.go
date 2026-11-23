// Package registrymatrix implements the Plugin Discovery & Health Matrix plugin.
package registrymatrix

import (
	"context"

	"plugins.dev/backend/internal/models"
	"plugins.dev/backend/internal/registry"
	"plugins.dev/backend/pkg/plugin"
)

const (
	PluginName    = "plugin-discovery-health-matrix"
	PluginVersion = "0.1.0"
)

// Plugin wraps the plugin registry as a standalone plugin.
type Plugin struct {
	reg *registry.Registry
}

// New creates a new registry matrix plugin.
func New() *Plugin { return &Plugin{reg: registry.New()} }

func (p *Plugin) Name() string    { return PluginName }
func (p *Plugin) Version() string { return PluginVersion }
func (p *Plugin) Description() string {
	return "Centralized index of all plugins across platforms with compatibility matrix, version status, and health fields."
}
func (p *Plugin) Category() string { return "registry" }

func (p *Plugin) Init(cfg plugin.Config) error { return nil }

func (p *Plugin) HealthCheck(ctx context.Context) (plugin.HealthStatus, error) {
	return plugin.HealthStatus{Healthy: true, Message: "registry ready"}, nil
}

// Registry exposes the underlying registry.
func (p *Plugin) Registry() *registry.Registry { return p.reg }

func (p *Plugin) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	action, _ := input["action"].(string)
	switch action {
	case "create":
		pl := models.Plugin{}
		if v, ok := input["id"].(string); ok {
			pl.ID = v
		}
		if v, ok := input["name"].(string); ok {
			pl.Name = v
		}
		if v, ok := input["platform"].(string); ok {
			pl.Platform = v
		}
		if err := p.reg.Create(pl); err != nil {
			return nil, err
		}
		return map[string]interface{}{"created": pl.ID}, nil
	case "list":
		platform, _ := input["platform"].(string)
		var plugins []models.Plugin
		if platform != "" {
			plugins = p.reg.ByPlatform(platform)
		} else {
			plugins = p.reg.List()
		}
		return map[string]interface{}{"plugins": plugins, "count": len(plugins)}, nil
	case "get":
		id, _ := input["id"].(string)
		pl, err := p.reg.Get(id)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{"plugin": pl}, nil
	default:
		return map[string]interface{}{"plugins": p.reg.List(), "count": len(p.reg.List())}, nil
	}
}
