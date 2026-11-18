# Image of the application for a production server.
#
# The build happens in two stages, so the runtime image holds no compiler and
# no node modules.

# --- Frontend -------------------------------------------------------------
FROM node:20-alpine AS frontend
WORKDIR /src/frontend

# The dependencies are installed before the sources are copied, so a change
# of the code does not install them again.
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci

COPY frontend/ ./
RUN npm run build

# --- Backend --------------------------------------------------------------
FROM golang:1.24-alpine AS backend
WORKDIR /src/backend

COPY backend/ ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/rekr-server ./cmd/server
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/rekr-backup ./cmd/backup

# --- Runtime --------------------------------------------------------------
FROM alpine:3.20

# The certificates are needed to reach the Google Sheets API.
RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -u 10001 -h /app rekrutmen

WORKDIR /app
COPY --from=backend /out/rekr-server /app/rekr-server
COPY --from=backend /out/rekr-backup /app/rekr-backup
COPY --from=frontend /src/frontend/dist /app/frontend

# The data directory holds the accounts and the session secret, so it belongs
# to the user of the application only.
RUN mkdir -p /app/data && chown -R rekrutmen:rekrutmen /app

USER rekrutmen

ENV ADDR=:8080 \
    APP_ROOT=/app \
    DATA_DIR=/app/data \
    FRONTEND_DIR=/app/frontend \
    COOKIE_SECURE=false

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://127.0.0.1:8080/api/health > /dev/null || exit 1

CMD ["/app/rekr-server"]
