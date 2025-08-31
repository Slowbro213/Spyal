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
	bunx markdownlint-cli2 README.md

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
	cd backend && go run cmd/register/main.go

# 🚀 Run backend dev server
dev:
	trap 'kill -TERM $$BACKEND_PID $$FRONTEND_PID 2>/dev/null; wait; exit 0' INT TERM EXIT; \
	cd backend && air & BACKEND_PID=$$!; \
	cd frontend && bun run dev & FRONTEND_PID=$$!; \
	wait

# 🐳 Build Docker image
docker-build:
	docker build -t alspy .

# 🐳 Run Docker container with env file
docker-run:
	docker run --rm\
		--env-file .env.production \
		-p 8080:8080 \
		-v /var/log/alspy.log:/var/log/alspy.log \
		alspy

# 🪝 Run Lefthook manually
precommit:
	lefthook run pre-commit

prepush:
	lefthook run pre-push
