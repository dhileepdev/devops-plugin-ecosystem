package operator

import (
	"strings"
	"testing"
)

const sampleManifest = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-operator
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: op-role
rules:
- verbs: ["get", "list"]
  resources: ["pods"]
`

func TestScanManifest(t *testing.T) {
	s := ScanManifest(sampleManifest)
	if s.APIVersion != "apps/v1" {
		t.Fatalf("expected apps/v1, got %s", s.APIVersion)
	}
	if s.UpgradePolicy != "missing" {
		t.Fatalf("expected missing upgrade policy, got %s", s.UpgradePolicy)
	}
	if len(s.RBACRoles) != 1 || s.RBACRoles[0] != "op-role" {
		t.Fatalf("unexpected rbac roles: %v", s.RBACRoles)
	}
	if s.HasPortForward {
		t.Fatal("expected no port-forward")
	}
}

func TestGenerateUpgradePath(t *testing.T) {
	s := ScanManifest(sampleManifest)
	up := GenerateUpgradePath(s, "v2.0.0")
	if up.TargetVersion != "v2.0.0" {
		t.Fatalf("unexpected target: %s", up.TargetVersion)
	}
	if len(up.Steps) == 0 {
		t.Fatal("expected steps")
	}
}

func TestAnalyzeRBAC_NoGaps(t *testing.T) {
	report := AnalyzeRBAC(sampleManifest)
	if len(report.Gaps) != 0 {
		t.Fatalf("expected no gaps, got %v", report.Gaps)
	}
}

func TestAnalyzeRBAC_ClusterAdmin(t *testing.T) {
	m := sampleManifest + "\ncluster-admin\n"
	report := AnalyzeRBAC(m)
	found := false
	for _, g := range report.Gaps {
		if strings.Contains(strings.Join(g.OverlyBroad, " "), "cluster-admin") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected cluster-admin gap, got %v", report.Gaps)
	}
}
