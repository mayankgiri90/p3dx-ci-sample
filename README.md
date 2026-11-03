# p3dx-ci-sample

A small Go HTTP service that exists to verify the P3DX Jenkins CI pipeline. It uses the same
Jenkinsfile, stages and thresholds as `datakaveri/p3dx_gov_layer` and
`datakaveri/p3dx-kata-coco-orchestrator`, and is written so that every stage passes, including the
80% unit-test coverage gate.

## API

| Route | Response |
|---|---|
| `GET /healthz` | `{"status":"ok"}` |
| `GET /version` | build version and commit |
| `GET /greet?name=<name>` | `{"message":"hello, <name>"}`; 400 if the name is missing or over 64 characters |

## Run locally

```sh
go test -cover ./...
go run ./cmd/server          # listens on :8080 (ADDR to change)
```

## CI

Jenkins job `p3dx-ci-sample` builds `main` and publishes `ghcr.io/mayankgiri90/p3dx-ci-sample`.
Reports (coverage, Semgrep, gitleaks, Trivy) are in each build's artifacts under `reports/`.
