# TokTik

TokTik is a distributed streaming platform built for a Scalable Systems class at MUIC. It combines ideas and features from platforms like YouTube and TikTok to demonstrate a scalable microservices architecture, streaming (HLS) workflows, and deployment using Docker Compose for local development and Kubernetes for production-like multi-node deployments.

This repository contains backend microservices, video processing workers, and a Next.js frontend. The project explores security, scalability, integrity, and consistency while practicing containerization and orchestration.

## Key features

- Microservices architecture (separate services for authentication, user, video processing, thumbnailing, websocket, etc.).
- HLS streaming using ffmpeg for creating .m3u8/.ts segments.
- Background workers for media processing (original project used Celery; this repository uses task worker tooling—see `machinery.yaml` files).
- Files stored in S3-compatible object storage.
- JWT-based authentication and authorization.
- Redis used as Pub/Sub and session/real-time support for WebSocket communication.
- Docker Compose config for local development and Kubernetes manifests for multi-node deployment.

## Repository layout (high level)

- `Backend/` — Main backend service (auth, API, etc.).
- `UserService/` — User management service.
- `VideoService/` — Core video metadata, publishing, and orchestration.
- `ThumbnailMaker/` — Worker service to generate thumbnails.
- `VideoChunker/` — Worker to chunk videos into HLS segments.
- `VideoConvertor/` — Worker to transcode videos (ffmpeg used here).
- `WebsocketService/` / `WebsocketServer/` — WebSocket services for real-time notifications.
- `frontend/` — Frontend code (Next.js) under `frontend/toktik_frontend`.
- `k8s/` — Kubernetes manifests, `kustomization.yaml`, and related configs.
- `docker-compose.yml` — Compose file for local development.

Note: For exact module-level details check each service's `go.mod` and internal folders.

## Prerequisites

- Docker & Docker Compose
- kubectl (and a Kubernetes cluster or local Kubernetes like Minikube / Kind)
- Node.js (for the frontend) — check `frontend/toktik_frontend/package.json` for recommended version
- Go toolchain to build Go services (check top-level `go.mod` and each service `go.mod`)
- Access to S3-compatible storage (AWS S3, MinIO, etc.)
- RabbitMQ or the configured message broker and Redis

## Quick start — Local development with Docker Compose

1. Copy or create environment files required by services. The repository expects several env variables (DB, S3, broker, redis, JWT secrets). For local testing you can use local Postgres, MinIO, RabbitMQ, and Redis. The `k8s/secret.yaml` contains keys used for Kubernetes — adapt them to `.env` files for Compose.

2. Start all services with Docker Compose (run from repository root):

```bash
# build images and start services
docker-compose up --build -d

# follow logs
docker-compose logs -f
```

3. Frontend (Next.js):

```bash
cd frontend/toktik_frontend
npm install
npm run dev
# or use Dockerfile in that folder to build an image used by compose
```

4. Verify health endpoints and APIs. Many services expose HTTP endpoints under ports declared in `docker-compose.yml`.

## Deploy to Kubernetes

This repo includes manifests in `k8s/` and a `kustomization.yaml` to aggregate resources.

1. Prepare cluster and apply secrets/configs. Replace placeholders with real secrets (S3 keys, DB password, JWT secret, etc.).
   - You can use `kubectl create secret` or edit `k8s/secret.yaml`.

2. Apply the kustomization:

```bash
kubectl apply -k k8s/
```

3. Check deployments and pods:

```bash
kubectl get pods,deploy,svc -n default
kubectl logs -l app=video-service
```

4. The `k8s/` folder contains per-service deployment YAMLs (video service, thumbnail maker, user service, websocket service, etc.). Tune resource requests/limits and replicas for horizontal scaling.

## Streaming pipeline (HLS)

- Uploaded videos are processed by background workers which run ffmpeg to create HLS segments (.m3u8 + .ts files).
- These files are stored in S3-compatible storage and served by the video service or a CDN.
- A task queue (the repo contains `machinery.yaml` and worker implementations) coordinates work. In earlier versions of this project Celery (Python) was used; this repository uses the Machinery-style worker setup—see each service's internal `machinery`/`machineryutil`/`machineryUtil` references.

## Configuration

- Look in each service under `internal/env_config` or `internal/config` for required environment variables and defaults.
- Secrets for Kubernetes are under `k8s/secret.yaml`.
- Common configuration items:
  - DATABASE_URL / Postgres connection
  - S3 endpoint, access key, secret key, bucket
  - MESSAGE_BROKER (RabbitMQ) connection
  - REDIS_URL for Pub/Sub and caching
  - JWT_SECRET and token settings

## Observability & scaling notes

- The app is designed for horizontal scaling — scale worker deployments (chunker/convertor/thumbnail) and API services independently.
- Use Kubernetes Horizontal Pod Autoscaler (HPA) to autoscale pods based on CPU or custom metrics.
- Use a CDN in front of S3 for better video delivery performance.

## Troubleshooting

- If workers are not processing tasks: ensure the message broker (RabbitMQ) is reachable and credentials match.
- If video segments are missing: check worker logs for ffmpeg errors and verify S3 upload permissions.
- Redis Pub/Sub issues: verify Redis is reachable and the correct channel names are used.
- When deploying to k8s, if pods stay in CrashLoopBackOff, run `kubectl describe pod <pod>` and `kubectl logs <pod>` to inspect.

## Development tips & next steps

- Add automated tests (unit + integration) for services. Tests will help with CI/CD.
- Add a GitHub Actions or other CI pipeline to build images and run basic integration tests before deploying.
- Add monitoring (Prometheus + Grafana) and structured logging to improve observability.
- Consider using a managed S3/CDN for production-level performance.

## Contributing

Contributions are welcome. If you're adding features or fixing bugs:

- Fork the repo and open a PR with a clear description of changes.
- Update or add tests where appropriate.
- Document any new environment variables or Kubernetes changes in this README or `k8s/` docs.

## License

This project is the personal project of the author. Add your preferred license file if you want to open-source it publicly.

---

If you want, I can:

- Make the README more specific to the exact configuration values used in this repo (I can scan service folders and list env vars).
- Add a short `try-it` section with exact commands for the most minimal local development setup (e.g., minimal compose subset to run frontend + backend + redis + postgres).

Which of those would you like next?
