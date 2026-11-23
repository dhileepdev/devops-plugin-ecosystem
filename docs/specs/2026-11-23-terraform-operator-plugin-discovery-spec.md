# DevOps Plugin Ecosystem Spec — Single Integrated Platform

**Date:** 2026-11-23  
**Status:** Approved — design confirmed; implementation plan next
**Scope:** Three high-impact plugin concepts unified into one integrated platform.  
**Target Audience:** DevOps engineers managing Jenkins, Terraform, Kubernetes operators, and multi-team tooling.

---

## 1. Purpose & Vision

Build a **single devops plugin discovery & health management layer** that addresses recurring community pain points across multiple platforms. Instead of teams maintaining separate plugins or fighting conflicting versions, this platform provides:

- A **centralized plugin registry** where teams publish shared utilities and request new ones.
- **Health monitoring** for each installed plugin (compatibility, last updated, crash history) so teams avoid breaking changes.
- **Terraform provider SDK quality gates** — surfacing bugs, panics, performance regressions from GitHub Issues.
- **Kubernetes operator lifecycle/RBAC audit** — detecting upgrade requirements, port-forward issues, RBAC drift.
- **Notification deduplication** across Slack, Discord, GitHub Actions.

This reduces tool sprawl, eliminates copy-paste duplication, and lowers maintenance burden by consolidating plugin management into one ecosystem.

---

## 2. Concept Overview

### Concept C (Plugin Discovery & Health Matrix) — Foundation Layer
Aggregates all known plugins (Jenkins, Terraform, K8s, etc.), checks compatibility matrices, reports version status, and surfaces unused/abandoned plugins. Serves as the central index.

### Concept B (Operator Lifecycle & RBAC Analyzer)
Scans running Kubernetes operators for:
- Compatibility with latest kubectl/k8s API versions.
- Required upgrade frequency (per KEP 1498: ≥9 months).
- RBAC rule gaps / least-privilege violations.
Generates an upgrade path and alternative port-forward setup.

### Concept A (Terraform Provider Health Dashboard)
Pulls provider issue data from `hashicorp/terraform-plugin-sdk` GitHub Issues, monitors schema generation panics, detects O(n²) performance regressions, and flags change-detection failures. Shows live plan consistency.

### Concept E (Smart Notification Deduplicator)
Centralizes notifications from GitHub Actions → Slack, Discord, email; deduplicates by event type + severity; routes to channels based on user preferences and escalation policies.

---

## 3. Architecture

```
┌─────────────────────────────────────────────────────┐
│                    User / Team                       │
│   (Jenkins, Terraform, K8s, custom CLI, etc.)      │
└──────────────┬──────────────────────────────────────┘
               │ publishes plugins
               ▼
┌─────────────────────────────────────────────────────┐
│           Plugin Registry (Concept C)                │
│   • Index of all installed plugins                     │
│   • Compatibility matrix (plugin ↔ platform)         │
│   • Last update date, version, author                 │
│   • Crash/history links                                │
└──────────────┬──────────────────────────────────────┘
               │ queries
               ▼
┌─────────────────────────────────────────────────────┐
│        Health Monitor (Concepts A+B+E)              │
│   ├─ Terraform Provider Health (A)                   │
│   ├─ Operator Lifecycle/RBAC Audit (B)              │
│   └─ Notification Deduplicator (E)                  │
└──────────────┬──────────────────────────────────────┘
               │
    ┌──────────┼──────────┐
    ▼          ▼          ▼
  Dashboards  Alerts   Notifications
  (Terraform  RBAC  (Slack/Discord/GitHub)
  Health)     Upgrade   Routing
```

### Data Flow (Terraform Example)

1. User runs `terraform init` → triggers Terraform Provider SDK.
2. Plugin Registry detects `hashicorp/terraform-plugin-sdk` vX.Y.Z installed.
3. Health Monitor queries GitHub Issues for bugs tagged `terraform-plugin-sdk`:
   - Panic on `(*GRPCProviderServer).GenerateResourceConfig`
   - Write-only args produce non-empty plans
   - ValuesSDKEquivalent O(n²) complexity
4. Dashboard shows: current plan stability score, recent bug count, performance regression flag.
5. If threshold breached → alert via Smart Notification Deduplicator.

---

## 4. Detailed Concepts

### 4.1 Plugin Discovery & Health Matrix (Concept C)

**Problem**: No centralized view of plugins across Jenkins, Terraform, K8s, Gradle, etc.; manual copy-paste of build/deploy utilities; risk of conflicts when teams request different plugins in a shared instance.

**Solution**:
- **Plugin Registry API**: REST endpoint `GET /plugins`, `POST /plugins/{id}`, `GET /plugins/{id}/versions`.
- **Compatibility Matrix**: Predefined tables mapping plugin → supported platforms (Jenkins, Terraform, K8s, Gradle, Azure DevOps).
- **Health Fields per Plugin**:
  - `last_updated`: ISO timestamp
  - `latest_version`: semver
  - `known_bugs`: list of issue IDs (from GitHub)
  - `panics`: runtime panic events (e.g., `ValidateProviderConfig`)
  - `performance_regression`: e.g., O(n²) vs expected linear
  - `unused`: boolean (no active usage in target org)
  - `maintainer_active`: boolean
- **UI**: Searchable index with filters (platform, last updated, severity, unused).

**Accessibility**: Public GitHub APIs for all listed sources; Reddit/HN searchable via websearch. Slack/Discord/LinkedIn gated but not part of core registry (notified externally).

---

### 4.2 Operator Lifecycle & RBAC Analyzer (Concept B)

**Problem**: Kubernetes operators require upgrades every ~9 months (KEP 1498); node status updates are clunky; RBAC rules drift over time; port-forward lifecycles must be re-implemented after each release.

**Solution**:
- Scan Helm chart manifests and operator containers for:
  - `kubectl version` / API version field.
  - `upgradePolicy` / `preStop` hooks.
  - RBAC roles (`Role`, `ClusterRole`, `ClusterRoleBinding`) and missing permissions.
  - Port-forward configuration blocks.
- Generate:
  - **Upgrade Path**: recommended version bump, required steps, estimated downtime.
  - **RBAC Gap Report**: missing permissions, overly permissive bindings.
  - **Port-Forward Fix**: alternative `--port-forward` config if none exists.
- Expose as a dashboard widget and CLI command `kubectl op operator describe <op> --show-last-upgrade`.

**Data Sources**: GitHub issues (#131208), DocumentDB operator e2e test suite, KEP 1498.

**Accessibility**: Public GitHub issues; KEP docs public. Requires operator repo access (gated for private repos, but public examples available).

---

### 4.3 Terraform Provider Health Dashboard (Concept A)

**Problem**: `hashicorp/terraform-plugin-sdk` has critical bugs (panic on `GenerateResourceConfig`, write-only arg plan issues, O(n²) nested objects, false bool evaluation). These cause non-deterministic plans and runtime crashes during apply.

**Solution**:
- Pull issue tracker data from `hashicorp/terraform-plugin-sdk` via GitHub API.
- Track issue resolution status per issue ID.
- Detect performance regressions by comparing plan output sizes for deeply nested objects.
- Flag change detection failures (nested list value known post-apply).
- Display a **health score** per provider: green/yellow/red based on:
  - Number of open critical/buggy issues.
  - Recent performance complaints.
  - Plan stability (consistency of `terraform plan`).
- Alert when score drops below threshold.

**Dashboard Widgets**:
| Metric | Description |
|--------|-------------|
| Provider Status | Green/Yellow/Red badge |
| Open Critical Bugs | Count + link |
| Performance Regressions | List of O(n²)/slow issues |
| Plan Stability | % consistent plans over last N runs |
| Known Unused Resources | Count |

**Accessibility**: Public GitHub Issues API; Stack Overflow Q&A. Web-based dashboard (public URL).

---

### 4.4 Smart Notification Deduplicator (Concept E)

**Problem**: GitHub Actions → Slack notifications, Discord outbound delivery broken, fragmented alerts across channels; missed critical events due to duplicates.

**Solution**:
- Centralize notification templates per event type (build success/failure, PR merge, deploy, RBAC change).
- Deduplicate by `(event_type, source, timestamp_window)`.
- Route to appropriate channel(s):
  - Slack: for team members subscribed.
  - Discord: for ops channel users.
  - Email: fallback.
- Escalate unread/duplicated notifications after N minutes.

**Features**:
- Configurable preference profile (which channels receive what).
- Integration hooks for GitHub Actions webhooks, Discord app tokens, Slack incoming webhook.
- Audit log of sent notifications.

**Accessibility**: Public webhook URLs; requires auth for Discord/Slack (gate behind login).

---

## 5. Security & Privacy

- All data sourced from public APIs (GitHub, StackOverflow, Reddit, Mastodon) or self-hosted registries.
- No secrets stored; only public issue IDs, issue titles, and issue numbers.
- RBAC within the plugin platform itself (role-based access to dashboards/alerts).
- GDPR-compliant: no personal data collected beyond what's already public.

---

## 6. Deployment Model

- **Frontend**: React (Next.js App Router) — provides SSR and SEO-friendly dashboards; confirmed by user choice.
- **Backend**: Go (preferred) or Python (FastAPI) — exposing REST API for plugin registry, health monitors, and notification deduplicator.
- **Hosting**: Vercel (frontend) + floci (backend microservices) — both selected by user.
- **Auth**: Optional role-based access for internal teams; public read-only for general consumers. Required for enterprise scaling (>100 users).

---

## 7. Success Criteria

| Metric | Target |
|--------|--------|
| Time to discover a plugin | < 30 seconds (index search) |
| Time to see provider health | < 10 seconds (dashboard refresh) |
| False positive rate on plugin health | < 5% |
| Reduction in duplicate notifications | > 60% reduction |
| Team satisfaction (survey) | Positive trend within 3 months |

---

## 8. Resolved Decisions (User Confirmed)

1. **Tech stack preference**: React (Next.js App Router) frontend.
2. **Deployment environment**: Vercel (frontend) + floci (backend).
3. **Team size / org scale**: Enterprise (>100 users) — requires scalable hosting and RBAC.
4. **Integration depth**: Separate modules per dashboard under one umbrella (Terraform Provider Health, Operator Lifecycle & RBAC, Plugin Discovery Matrix).

Spec approved for implementation plan.
