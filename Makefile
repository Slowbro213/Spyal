# ────────────────────────────────────────────────────────────────
# 🛠️ Makefile for ALSpy Project (Go + Bun Frontend)
# ────────────────────────────────────────────────────────────────

# 🔍 Lint all source files
lint:
	cd backend && golangci-lint run ./...
	cd frontend && bunx eslint "src/**/*.ts" --max-warnings=0
	cd frontend && bunx stylelint "public/**/*.css"
	cd frontend && bunx htmlhint "views/**/*.html"
	cd frontend && bunx prettier --check .
	hadolint Dockerfile

# Format frontend files
format:
	cd frontend && bunx prettier --write . --ignore-path .gitignore

# 🔬 Run all tests
test:
	cd backend && go test ./... -v -cover -race -shuffle=on
	cd frontend && bun test --coverage

# 🔐 Run security scans
secure:
	cd backend && gosec ./...
	cd frontend && bun audit

# ✅ Full backend check (fmt, vet, lint, test)
check:
	cd backend && go fmt ./...
	cd backend && go vet ./...
	cd backend && golangci-lint run ./...
	cd backend && go test ./... -v -cover -race -shuffle=on

# 🏗️ Build frontend
build:
	cd frontend && bun run build

# 🚀 Run backend dev server
dev:
	cd backend/cmd/server && go run main.go

# 🪝 Run Lefthook manually
precommit:
	lefthook run pre-commit

prepush:
	lefthook run pre-push
