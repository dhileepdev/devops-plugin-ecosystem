# DevOps Plugin Ecosystem

Unified plugin discovery & health management platform (Concepts A, B, C, E).

## Stack
- Backend: Go 1.24 microservice (chi router) — `backend/`
- Frontend: Next.js 14 App Router + TypeScript — `frontend/`

## Run locally

### Backend
```sh
cd backend
go mod tidy
go run ./cmd/server   # listens on :8080
```

### Frontend
```sh
cd frontend
npm install
npm run dev           # http://localhost:3000
```

Frontend proxies `/api/*` to `http://localhost:8080` via `next.config.js`.

## API endpoints

- `GET /api/health` — status
- `GET /api/health/terraform` — Terraform provider health report (GitHub issues)
- `GET /api/health/operator` — operator analyzer status
- `GET /api/plugins` — list plugins (`?platform=terraform`)
- `GET /api/plugins/{id}` — get plugin
- `POST /api/operator/analyze` — scan manifest, return upgrade path + RBAC gaps
- `POST /api/notifications` — ingest notification (dedup + route)
- `GET /api/notifications` — audit log
- `GET /api/notifications/stats` — counts

## Deploy
- Frontend: Vercel (`frontend/vercel.json`)
- Backend: floci (`deploy/floci/floci.yaml`)

## Tests
```sh
cd backend && go test ./...
cd frontend && npx tsc --noEmit
```
