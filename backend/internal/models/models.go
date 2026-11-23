package models

import "time"

// Plugin is a shared utility plugin published by a team.
type Plugin struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Platform         string    `json:"platform"` // terraform | k8s | jenkins | gradle | azure-devops
	Version          string    `json:"version"`
	LatestVersion    string    `json:"latest_version"`
	LastUpdated      time.Time `json:"last_updated"`
	MaintainerActive bool      `json:"maintainer_active"`
	Unused           bool      `json:"unused"`
}

// Validate returns an error if required fields are missing.
func (p *Plugin) Validate() error {
	if p.ID == "" {
		return errMissingID
	}
	if p.Name == "" {
		return errMissingName
	}
	return nil
}

var (
	errMissingID   = errNew("missing plugin id")
	errMissingName = errNew("missing plugin name")
)

type errStr string

func errNew(s string) errStr { return errStr(s) }

func (e errStr) Error() string { return string(e) }
