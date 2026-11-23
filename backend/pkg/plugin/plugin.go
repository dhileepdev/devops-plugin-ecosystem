// Package plugin defines the common contract every devops plugin implements.
// Plugins can be used standalone (via CLI) or composed into the platform.
package plugin

import "context"

// Plugin is the interface every plugin must satisfy.
type Plugin interface {
	// Name returns the plugin's unique identifier.
	Name() string
	// Version returns the semver version.
	Version() string
	// Description returns a human-readable summary.
	Description() string
	// Category returns the plugin category (health, analyzer, dedup, registry).
	Category() string
	// Init prepares the plugin with configuration.
	Init(cfg Config) error
	// HealthCheck reports whether the plugin is operational.
	HealthCheck(ctx context.Context) (HealthStatus, error)
	// Execute runs the plugin's primary action with the given input.
	Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error)
}

// Config holds plugin configuration.
type Config struct {
	GitHubToken string            `json:"github_token"`
	Extra       map[string]string `json:"extra"`
}

// HealthStatus is the runtime health of a plugin.
type HealthStatus struct {
	Healthy bool   `json:"healthy"`
	Message string `json:"message"`
}

// Manifest describes a plugin for discovery.
type Manifest struct {
	Name         string   `json:"name"`
	Version      string   `json:"version"`
	Description  string   `json:"description"`
	Category     string   `json:"category"`
	Platforms    []string `json:"platforms"`
	Entrypoint   string   `json:"entrypoint"`
	Dependencies []string `json:"dependencies"`
	APIEndpoints []string `json:"api_endpoints"`
	CLICommands  []string `json:"cli_commands"`
}

// Registry is a collection of loaded plugins.
type Registry struct {
	plugins map[string]Plugin
}

// NewRegistry creates an empty plugin registry.
func NewRegistry() *Registry {
	return &Registry{plugins: make(map[string]Plugin)}
}

// Register adds a plugin.
func (r *Registry) Register(p Plugin) {
	r.plugins[p.Name()] = p
}

// Get retrieves a plugin by name.
func (r *Registry) Get(name string) (Plugin, bool) {
	p, ok := r.plugins[name]
	return p, ok
}

// List returns all registered plugin names.
func (r *Registry) List() []string {
	out := make([]string, 0, len(r.plugins))
	for n := range r.plugins {
		out = append(out, n)
	}
	return out
}

// All returns all registered plugins.
func (r *Registry) All() []Plugin {
	out := make([]Plugin, 0, len(r.plugins))
	for _, p := range r.plugins {
		out = append(out, p)
	}
	return out
}
