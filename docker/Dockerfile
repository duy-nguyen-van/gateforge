# syntax=docker/dockerfile:1

# Stage 1: frontend
FROM node:26-alpine3.24 AS frontend-builder

WORKDIR /app/frontend

COPY frontend/package.json frontend/package-lock.json ./
RUN --mount=type=cache,target=/root/.npm \
    npm ci

COPY frontend/ ./
RUN npm run build

# Stage 2: backend
FROM golang:1.26-alpine3.24 AS backend-builder

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY backend/go.mod backend/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY backend/ ./
COPY --from=frontend-builder /app/frontend/dist ./internal/static/dist

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build \
    -buildvcs=false \
    -tags embedfrontend \
    -ldflags="-s -w" \
    -o server ./cmd/server

# Stage 3: runtime
FROM alpine:3.23

RUN apk --no-cache add ca-certificates wget && \
    adduser -D -s /bin/sh appuser

WORKDIR /app

COPY --from=backend-builder /app/server .

RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 3000

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000/api/v1/ || exit 1

CMD ["./server"]
