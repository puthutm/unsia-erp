# Dockerfile Update Plan

## Overview
Update all 10 Go service Dockerfiles to follow the user's provided template with Go 1.23.0.

## Current vs Updated Configuration

| Aspect | Current | Updated |
|--------|---------|---------|
| Go Version | golang:1.22-alpine | golang:1.23.0-alpine |
| Multi-stage | Yes | Yes |
| .env copy | Not copied | COPY .env . |
| *.pem copy | Not copied | COPY *.pem . |
| Ports | 8001-8010 | Keep current ports (8001-8010) |
| Binary paths | cmd/*-service | Keep actual paths |
| NETRC handling | Not implemented | ARG NETRC_FILE + chmod 600 |
| Builder cleanup | None | RUN rm -f /root/.netrc |
| bash in runtime | Not installed | apk add bash |

## Services and Their Configuration

| Service | Port | Binary Path | Dockerfile Path |
|---------|------|------------|-----------------|
| unsia-core-service | 8001 | cmd/core-service | services/unsia-core-service/Dockerfile |
| unsia-reference-service | 8002 | cmd/reference-service | services/unsia-reference-service/Dockerfile |
| unsia-pmb-service | 8003 | cmd/pmb-service | services/unsia-pmb-service/Dockerfile |
| unsia-academic-service | 8004 | cmd/academic-service | services/unsia-academic-service/Dockerfile |
| unsia-finance-service | 8005 | cmd/finance-service | services/unsia-finance-service/Dockerfile |
| unsia-lms-service | 8006 | cmd/lms-service | services/unsia-lms-service/Dockerfile |
| unsia-assessment-service | 8007 | cmd/assessment-service | services/unsia-assessment-service/Dockerfile |
| unsia-hris-service | 8008 | cmd/hris-service | services/unsia-hris-service/Dockerfile |
| unsia-crm-service | 8009 | cmd/crm-service | services/unsia-crm-service/Dockerfile |
| unsia-portal-service | 8010 | cmd/portal-service | services/unsia-portal-service/Dockerfile |

## Update Strategy

Each Dockerfile will be updated to follow this pattern:

### Stage 1: Builder
```dockerfile
FROM golang:1.23.0-alpine as builder
RUN apk add --no-cache git gcc musl-dev
WORKDIR /app
ARG NETRC_FILE=.netrc
COPY ${NETRC_FILE} /root/.netrc
RUN chmod 600 /root/.netrc
COPY services/<service-name>/go.mod services/<service-name>/go.sum* ./services/<service-name>/
COPY packages/ ./packages/
WORKDIR /app/services/<service-name>
COPY services/<service-name>/ ./
RUN go mod tidy
RUN CGO_ENABLED=1 go build -o /<binary-name> ./cmd/<cmd-path>
RUN rm -f /root/.netrc
```

### Stage 2: Runner
```dockerfile
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /<binary-name> /app/<binary-name>
COPY .env .
COPY *.pem .
EXPOSE <port>
CMD ["/bin/sh", "-c", "apk update && apk add bash"]
CMD ["/app/<binary-name>"]
```

## Implementation Order
1. unsia-core-service (8001)
2. unsia-reference-service (8002)
3. unsia-pmb-service (8003)
4. unsia-academic-service (8004)
5. unsia-finance-service (8005)
6. unsia-lms-service (8006)
7. unsia-assessment-service (8007)
8. unsia-hris-service (8008)
9. unsia-crm-service (8009)
10. unsia-portal-service (8010)

## Notes
- Each service already has its own .env file (verified for finance, reference, pmb)
- Keep actual ports to avoid breaking docker-compose.yml
- Keep actual binary paths from current Dockerfiles
- The user template had `cmd/web/app_name/main.go` but actual services use cmd/*-service paths
