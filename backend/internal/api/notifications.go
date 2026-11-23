package api

import (
	"encoding/json"
	"net/http"
	"time"

	"plugins.dev/backend/internal/notify"
)

type NotificationHandler struct {
	engine *notify.Engine
}

func NewNotificationHandler(engine *notify.Engine) *NotificationHandler {
	if engine == nil {
		engine = notify.NewEngine(5 * time.Minute)
	}
	return &NotificationHandler{engine: engine}
}

func (h *NotificationHandler) Ingest(w http.ResponseWriter, r *http.Request) {
	var n notify.Notification
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		http.Error(w, `{"error":"bad request"}`, http.StatusBadRequest)
		return
	}
	rn := h.engine.Ingest(n)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(rn)
}

func (h *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {
	log := h.engine.AuditLog()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(log)
}

func (h *NotificationHandler) Stats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(h.engine.Stats())
}
