// CLI for individually using each devops plugin.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"plugins.dev/backend/pkg/plugin"

	notificationdedup "plugins.dev/backend/plugins/notification-dedup"
	operatoranalyzer "plugins.dev/backend/plugins/operator-analyzer"
	registrymatrix "plugins.dev/backend/plugins/registry-matrix"
	terraformhealth "plugins.dev/backend/plugins/terraform-health"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	cfg := plugin.Config{GitHubToken: os.Getenv("GITHUB_TOKEN")}

	reg := plugin.NewRegistry()
	reg.Register(terraformhealth.New())
	reg.Register(operatoranalyzer.New())
	reg.Register(notificationdedup.New())
	reg.Register(registrymatrix.New())

	for _, p := range reg.All() {
		if err := p.Init(cfg); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "init %s: %v\n", p.Name(), err)
			os.Exit(1)
		}
	}

	ctx := context.Background()

	switch cmd {
	case "list":
		_, _ = fmt.Println("Available plugins:")
		for _, name := range reg.List() {
			p, _ := reg.Get(name)
			_, _ = fmt.Printf("  %-32s %s (%s) - %s\n", p.Name(), p.Version(), p.Category(), p.Description())
		}
	case "health":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: cli health <plugin-name>")
			os.Exit(1)
		}
		p, ok := reg.Get(os.Args[2])
		if !ok {
			fmt.Fprintf(os.Stderr, "unknown plugin: %s\n", os.Args[2])
			os.Exit(1)
		}
		st, err := p.HealthCheck(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "health check failed: %v\n", err)
			os.Exit(1)
		}
		b, _ := json.MarshalIndent(st, "", "  ")
		fmt.Println(string(b))
	case "run":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: cli run <plugin-name> [json-input]")
			os.Exit(1)
		}
		p, ok := reg.Get(os.Args[2])
		if !ok {
			fmt.Fprintf(os.Stderr, "unknown plugin: %s\n", os.Args[2])
			os.Exit(1)
		}
		input := map[string]interface{}{}
		if len(os.Args) >= 4 {
			if err := json.Unmarshal([]byte(os.Args[3]), &input); err != nil {
				fmt.Fprintf(os.Stderr, "bad json input: %v\n", err)
				os.Exit(1)
			}
		}
		out, err := p.Execute(ctx, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "execute failed: %v\n", err)
			os.Exit(1)
		}
		b, _ := json.MarshalIndent(out, "", "  ")
		fmt.Println(string(b))
	case "manifest":
		printManifests()
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: devops-plugins <command> [args]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  list                      List all available plugins")
	fmt.Println("  health <plugin>           Health check a plugin")
	fmt.Println("  run <plugin> [json]       Execute a plugin with optional JSON input")
	fmt.Println("  manifest                  Print plugin manifests")
}

func printManifests() {
	manifests := []plugin.Manifest{
		{
			Name:         "terraform-provider-health",
			Version:      "0.1.0",
			Description:  "Surfaces bugs, panics, and performance regressions from hashicorp/terraform-plugin-sdk GitHub Issues.",
			Category:     "health",
			Platforms:    []string{"terraform"},
			Entrypoint:   "plugins/terraform-health",
			APIEndpoints: []string{"/api/health/terraform"},
			CLICommands:  []string{"run terraform-provider-health"},
		},
		{
			Name:         "operator-lifecycle-rbac",
			Version:      "0.1.0",
			Description:  "Scans Kubernetes operator manifests for upgrade requirements, RBAC gaps, and port-forward config.",
			Category:     "analyzer",
			Platforms:    []string{"kubernetes"},
			Entrypoint:   "plugins/operator-analyzer",
			APIEndpoints: []string{"/api/operator/analyze"},
			CLICommands:  []string{"run operator-lifecycle-rbac"},
		},
		{
			Name:         "smart-notification-dedup",
			Version:      "0.1.0",
			Description:  "Deduplicates and routes notifications across Slack, Discord, and email.",
			Category:     "dedup",
			Platforms:    []string{"github-actions", "slack", "discord"},
			Entrypoint:   "plugins/notification-dedup",
			APIEndpoints: []string{"/api/notifications"},
			CLICommands:  []string{"run smart-notification-dedup"},
		},
		{
			Name:         "plugin-discovery-health-matrix",
			Version:      "0.1.0",
			Description:  "Centralized index of all plugins across platforms with compatibility matrix.",
			Category:     "registry",
			Platforms:    []string{"jenkins", "terraform", "kubernetes", "gradle", "azure-devops"},
			Entrypoint:   "plugins/registry-matrix",
			APIEndpoints: []string{"/api/plugins"},
			CLICommands:  []string{"run plugin-discovery-health-matrix"},
		},
	}
	b, _ := json.MarshalIndent(manifests, "", "  ")
	fmt.Println(string(b))
}
