# ────────────────────────────────────────────────────────────────
# 🛠️ Makefile for ALSpy Project (Go + Bun Frontend)
# ────────────────────────────────────────────────────────────────

# 🔍 Lint all source files
lint:
	cd backend && golangci-lint run ./...
	cd frontend && bunx eslint . --ext .ts
	cd frontend && bunx stylelint "public/**/*.css"
	cd frontend && bunx htmlhint "views/**/*.html"
	cd frontend && bunx prettier --check "src/**/*.{ts,css,html,json,md}"

# Format frontend files
format:
	bunx --cwd frontend prettier --write "src/**/*.{ts,css,html,json,md}"

# Run Go and frontend tests
test:
	cd backend && go test ./...
	cd frontend && bun test

# Full backend check
check:
	cd backend && go fmt ./...
	cd backend && go vet ./...
	cd backend && golangci-lint run ./...
	cd backend && go test ./... -cover

# Build frontend
build:
	cd frontend && bun run build

# Run backend dev server
dev:
	cd backend/cmd/server && go run main.go

# Lefthook
precommit:
	lefthook run pre-commit
prepush:
	lefthook run pre-push
