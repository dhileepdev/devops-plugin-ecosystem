package notify

import (
	"sync"
	"time"
)

// Severity levels.
type Severity string

const (
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityCritical Severity = "critical"
)

// Notification is an inbound alert.
type Notification struct {
	ID        string    `json:"id"`
	EventType string    `json:"event_type"`
	Source    string    `json:"source"`
	Severity  Severity  `json:"severity"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Channels  []string  `json:"channels"`
}

// RoutedNotification is a notification after dedup + routing.
type RoutedNotification struct {
	Notification
	Deduplicated bool     `json:"deduplicated"`
	RoutedTo     []string `json:"routed_to"`
}

// Engine deduplicates and routes notifications.
type Engine struct {
	mu     sync.Mutex
	seen   map[string]time.Time // dedup key -> last seen
	window time.Duration
	prefs  map[string][]string // event type -> channels
	audit  []RoutedNotification
}

// NewEngine creates a dedup engine with the given window.
func NewEngine(window time.Duration) *Engine {
	if window == 0 {
		window = 5 * time.Minute
	}
	return &Engine{
		seen:   make(map[string]time.Time),
		window: window,
		prefs:  defaultPrefs(),
	}
}

func defaultPrefs() map[string][]string {
	return map[string][]string{
		"build.failure": {"slack", "discord"},
		"build.success": {"slack"},
		"deploy":        {"slack", "discord", "email"},
		"rbac.change":   {"slack", "email"},
		"pr.merge":      {"slack"},
	}
}

// SetPreference overrides routing for an event type.
func (e *Engine) SetPreference(eventType string, channels []string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.prefs[eventType] = channels
}

// Ingest processes a notification: dedup + route.
//
//nolint:gocritic // passing by value is acceptable
func (e *Engine) Ingest(n Notification) RoutedNotification {
	e.mu.Lock()
	defer e.mu.Unlock()

	if n.Timestamp.IsZero() {
		n.Timestamp = time.Now()
	}
	key := n.EventType + "|" + n.Source
	// bucket timestamp into window
	bucket := n.Timestamp.Truncate(e.window)
	key = key + "|" + bucket.Format(time.RFC3339)

	rn := RoutedNotification{Notification: n}
	if last, ok := e.seen[key]; ok && n.Timestamp.Sub(last) < e.window {
		rn.Deduplicated = true
		rn.RoutedTo = []string{}
		e.audit = append(e.audit, rn)
		return rn
	}
	e.seen[key] = n.Timestamp

	channels := e.prefs[n.EventType]
	if len(channels) == 0 {
		channels = []string{"slack"}
	}
	if len(n.Channels) > 0 {
		channels = n.Channels
	}
	rn.RoutedTo = channels
	e.audit = append(e.audit, rn)
	return rn
}

// AuditLog returns all routed notifications.
func (e *Engine) AuditLog() []RoutedNotification {
	e.mu.Lock()
	defer e.mu.Unlock()
	out := make([]RoutedNotification, len(e.audit))
	copy(out, e.audit)
	return out
}

// Stats returns counts.
func (e *Engine) Stats() map[string]int {
	e.mu.Lock()
	defer e.mu.Unlock()
	total := len(e.audit)
	deduped := 0
	for _, r := range e.audit { //nolint:gocritic
		if r.Deduplicated {
			deduped++
		}
	}
	return map[string]int{
		"total":        total,
		"deduplicated": deduped,
		"sent":         total - deduped,
	}
}
