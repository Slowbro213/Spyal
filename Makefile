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
	watchexec -w frontend \
		--ignore frontend/public/out --ignore frontend/node_modules \
		"cd frontend && bun run build"
	wait

set-domain:
	@echo "🔍 Detecting current machine IP..."
	@IP=$$(if command -v ip &> /dev/null; then \
	        ip route get 1 | awk '{print $$7; exit}'; \
	    elif command -v ifconfig &> /dev/null; then \
	        ifconfig | grep 'inet ' | grep -v 127.0.0.1 | awk '{print $$2}' | head -n1; \
	    else \
	        echo ""; \
	    fi); \
	if [ -z "$$IP" ]; then \
	    echo "❌ Could not detect IP"; \
	    exit 1; \
	fi; \
	echo "🖧 Detected IP: $$IP"; \
	sed -i.bak -E "s|^DOMAIN=.*$$|DOMAIN=$$IP:8080|" .env.public; \
	echo "✅ Updated DOMAIN in .env.public to $$IP:8080"


# 🐳 Build Docker image
docker-build:
	docker build -t spyal .

# 🐳 Run Docker container with env file
docker-run:
	docker run --rm\
		--env-file .env.production \
		-p 8080:8080 \
		-v /var/log/alspy.log:/var/log/alspy.log \
		spyal

# 🪝 Run Lefthook manually
precommit:
	lefthook run pre-commit

prepush:
	lefthook run pre-push

# 🔑 Generate and set TOKEN_SECRET in .env.production
set-token-secret:
	@echo "🔑 Generating secure TOKEN_SECRET..."
	@SECRET=$$(openssl rand -hex 32); \
	if grep -q "^TOKEN_SECRET=" .env.production; then \
		sed -i.bak -E "s|^TOKEN_SECRET=.*$$|TOKEN_SECRET=$$SECRET|" .env.production; \
	else \
		echo "TOKEN_SECRET=$$SECRET" >> .env.production; \
	fi; \
	echo "✅ TOKEN_SECRET updated in .env.production"

