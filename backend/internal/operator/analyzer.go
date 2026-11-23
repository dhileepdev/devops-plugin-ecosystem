package operator

import (
	"regexp"
	"strings"
)

// ManifestScan holds extracted facts from a Kubernetes operator manifest.
type ManifestScan struct {
	KubectlVersion    string   `json:"kubectl_version"`
	APIVersion        string   `json:"api_version"`
	UpgradePolicy     string   `json:"upgrade_policy"`
	RBACRoles         []string `json:"rbac_roles"`
	HasPortForward    bool     `json:"has_port_forward"`
	PortForwardConfig string   `json:"port_forward_config"`
}

// UpgradePath is a recommended upgrade sequence.
type UpgradePath struct {
	CurrentVersion    string   `json:"current_version"`
	TargetVersion     string   `json:"target_version"`
	Steps             []string `json:"steps"`
	EstimatedDowntime string   `json:"estimated_downtime"`
}

// RBACGap is a single RBAC finding.
type RBACGap struct {
	Role        string   `json:"role"`
	Missing     []string `json:"missing"`
	OverlyBroad []string `json:"overly_broad"`
}

// RBACReport aggregates RBAC findings.
type RBACReport struct {
	Gaps []RBACGap `json:"gaps"`
}

// ScanManifest extracts facts from raw manifest text.
func ScanManifest(manifest string) ManifestScan {
	s := ManifestScan{}
	lower := strings.ToLower(manifest)

	// kubectl version
	if m := regexp.MustCompile(`kubectl[:\s]+v?(\d+\.\d+\.\d+)`).FindStringSubmatch(manifest); len(m) > 1 {
		s.KubectlVersion = m[1]
	}
	// apiVersion
	if m := regexp.MustCompile(`apiVersion:\s*(\S+)`).FindStringSubmatch(manifest); len(m) > 1 {
		s.APIVersion = m[1]
	}
	// upgradePolicy
	if strings.Contains(lower, "upgradepolicy") {
		if m := regexp.MustCompile(`upgradePolicy:\s*(\S+)`).FindStringSubmatch(manifest); len(m) > 1 {
			s.UpgradePolicy = m[1]
		} else {
			s.UpgradePolicy = "present"
		}
	} else {
		s.UpgradePolicy = "missing"
	}
	// RBAC roles
	roleRe := regexp.MustCompile(`kind:\s*(ClusterRole|Role|ClusterRoleBinding|RoleBinding)\s*\n\s*metadata:\s*\n\s*name:\s*(\S+)`)
	for _, m := range roleRe.FindAllStringSubmatch(manifest, -1) {
		s.RBACRoles = append(s.RBACRoles, m[2])
	}
	// port-forward
	if strings.Contains(lower, "port-forward") || strings.Contains(lower, "portforward") {
		s.HasPortForward = true
		s.PortForwardConfig = "present"
	} else {
		s.HasPortForward = false
		s.PortForwardConfig = "missing"
	}
	return s
}

// GenerateUpgradePath produces an upgrade path from a scan.
//
//nolint:gocritic // passing by value is acceptable
func GenerateUpgradePath(scan ManifestScan, targetVersion string) UpgradePath {
	if targetVersion == "" {
		targetVersion = "latest"
	}
	up := UpgradePath{
		CurrentVersion:    scan.KubectlVersion,
		TargetVersion:     targetVersion,
		EstimatedDowntime: "5-10 minutes",
	}
	steps := []string{}
	if scan.UpgradePolicy == "missing" {
		steps = append(steps, "Add upgradePolicy hook to manifest")
	}
	if !scan.HasPortForward {
		steps = append(steps, "Configure --port-forward alternative")
	}
	steps = append(steps, "Bump operator image to "+targetVersion) //nolint:gocritic
	steps = append(steps, "Apply manifest: kubectl apply -f operator.yaml")
	steps = append(steps, "Verify: kubectl rollout status deployment/operator")
	up.Steps = steps
	return up
}

// AnalyzeRBAC detects RBAC gaps from manifest text.
func AnalyzeRBAC(manifest string) RBACReport {
	report := RBACReport{Gaps: []RBACGap{}}
	// detect overly broad: ClusterRoleBinding with cluster-admin
	if strings.Contains(manifest, "cluster-admin") {
		report.Gaps = append(report.Gaps, RBACGap{
			Role:        "cluster-admin",
			OverlyBroad: []string{"cluster-admin binding grants full cluster access"},
		})
	}
	// detect missing permissions: look for verbs list
	if !strings.Contains(manifest, "verbs:") {
		report.Gaps = append(report.Gaps, RBACGap{
			Role:    "unknown",
			Missing: []string{"verbs not declared"},
		})
	}
	// if no roles at all
	if !strings.Contains(manifest, "kind: Role") && !strings.Contains(manifest, "kind: ClusterRole") {
		report.Gaps = append(report.Gaps, RBACGap{
			Role:    "none",
			Missing: []string{"no RBAC roles defined"},
		})
	}
	return report
}
