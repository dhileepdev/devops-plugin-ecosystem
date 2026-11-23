package api

import (
	"encoding/json"
	"net/http"

	"plugins.dev/backend/internal/github"
	"plugins.dev/backend/internal/health"
	"plugins.dev/backend/internal/operator"
)

type HealthHandler struct {
	gh *github.Client
}

func NewHealthHandler(gh *github.Client) *HealthHandler {
	if gh == nil {
		gh = github.NewClient("")
	}
	return &HealthHandler{gh: gh}
}

func (h *HealthHandler) Terraform(w http.ResponseWriter, r *http.Request) {
	issues, err := h.gh.FetchTerraformSDKIssues()
	if err != nil {
		// fall back to empty list so dashboard still renders
		issues = []github.Issue{}
	}
	report := health.ComputeHealth("hashicorp/terraform-plugin-sdk", issues)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(report)
}

// Operator placeholder endpoint
func (h *HealthHandler) Operator(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "operator analyzer ready"})
}

// OperatorAnalyze scans a manifest and returns upgrade path + RBAC report.
func (h *HealthHandler) OperatorAnalyze(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Manifest      string `json:"manifest"`
		TargetVersion string `json:"target_version"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"bad request"}`, http.StatusBadRequest)
		return
	}
	scan := operator.ScanManifest(req.Manifest)
	up := operator.GenerateUpgradePath(scan, req.TargetVersion)
	rbac := operator.AnalyzeRBAC(req.Manifest)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"scan":         scan,
		"upgrade_path": up,
		"rbac":         rbac,
	})
}
