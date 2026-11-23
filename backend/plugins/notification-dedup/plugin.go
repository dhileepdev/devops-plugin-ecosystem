// Package notificationdedup implements the Smart Notification Deduplicator plugin.
package notificationdedup

import (
	"context"
	"time"

	"plugins.dev/backend/internal/notify"
	"plugins.dev/backend/pkg/plugin"
)

const (
	PluginName    = "smart-notification-dedup"
	PluginVersion = "0.1.0"
)

// Plugin wraps the notification dedup engine as a standalone plugin.
type Plugin struct {
	engine *notify.Engine
}

// New creates a new notification dedup plugin.
func New() *Plugin { return &Plugin{} }

func (p *Plugin) Name() string    { return PluginName }
func (p *Plugin) Version() string { return PluginVersion }
func (p *Plugin) Description() string {
	return "Deduplicates and routes notifications across Slack, Discord, and email by (event_type, source, timestamp_window)."
}
func (p *Plugin) Category() string { return "dedup" }

func (p *Plugin) Init(cfg plugin.Config) error {
	p.engine = notify.NewEngine(5 * time.Minute)
	return nil
}

func (p *Plugin) HealthCheck(ctx context.Context) (plugin.HealthStatus, error) {
	if p.engine == nil {
		return plugin.HealthStatus{Healthy: false, Message: "not initialized"}, nil
	}
	return plugin.HealthStatus{Healthy: true, Message: "engine ready"}, nil
}

// Engine exposes the underlying dedup engine.
func (p *Plugin) Engine() *notify.Engine { return p.engine }

func (p *Plugin) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	n := notify.Notification{}
	if v, ok := input["id"].(string); ok {
		n.ID = v
	}
	if v, ok := input["event_type"].(string); ok {
		n.EventType = v
	}
	if v, ok := input["source"].(string); ok {
		n.Source = v
	}
	if v, ok := input["severity"].(string); ok {
		n.Severity = notify.Severity(v)
	}
	if v, ok := input["message"].(string); ok {
		n.Message = v
	}
	rn := p.engine.Ingest(n)
	return map[string]interface{}{
		"deduplicated": rn.Deduplicated,
		"routed_to":    rn.RoutedTo,
		"notification": rn.Notification,
	}, nil
}
