# ─── Backend Stage ───────────────────────────────────────────────
FROM golang:1.24.4 AS backend
WORKDIR /app
COPY backend ./backend
WORKDIR /app/backend
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/server ./cmd/server

# ─── Frontend Stage ──────────────────────────────────────────────
FROM oven/bun AS bun
WORKDIR /frontend
COPY frontend/package.json ./package.json
COPY frontend/bun.lock ./bun.lock
COPY frontend/public ./public
COPY frontend/src ./src
COPY frontend/views ./views
RUN bun i -p && bun run build

# ─── Final Stage ─────────────────────────────────────────────────
FROM scratch
WORKDIR /app

# Copy backend binary
COPY --from=backend /app/server ./server

# Copy static frontend assets
COPY --from=bun /frontend/public ./public
COPY --from=bun /frontend/views ./views

# Expose port (optional, if server binds to 8080)
EXPOSE 8080

# Run server
CMD ["./server"]
