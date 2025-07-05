# ALSpy 🇦🇱🔍

A modern, minimalist, Albanian-language version of the Spyfall game — built with Go, TypeScript, HTML, CSS, and powered by Bun for the frontend.

---

## 🔧 Tech Stack

- **Backend**: Go (net/http, server-side templating)
- **Frontend**: Vanilla JS, TypeScript, HTML, CSS (no framework)
- **Build Tool**: [Bun](https://bun.sh/)
- **Linting**: ESLint, Prettier, Stylelint
- **Testing**: Bun test runner
- **Git Hooks**: Lefthook (for linting & tests pre-commit/push)

---

## 🧪 Development Setup

### 1. Install Dependencies

```bash
bun install
```

### 2. Build TypeScript

```bash
bun run build
```

### 3. Start Server

```bash
go run cmd/server/main.go
```

The app will be available at `http://localhost:8080`.

---

## 🛡️ Pre-commit & Pre-push Hooks with Lefthook

This project uses [Lefthook](https://github.com/evilmartians/lefthook) to automate tasks before committing and pushing.

### ✅ Features:

- Run **ESLint**, **Stylelint**, and **Go vet/test** before each commit.
- Run **full test suite** before pushing.

### 💻 Setup:

1. Install Lefthook (using Go):

```bash
go install github.com/evilmartians/lefthook@latest
```

2. Enable git hooks:

```bash
lefthook install
```

3. Now, Lefthook will automatically:
   - Lint your code before commits.
   - Run tests before pushes.

---

## 📁 Folder Structure

```
.
├── assets/         # Static assets (CSS, images, TS/JS)
├── cmd/server/     # Go HTTP server entrypoint
├── handlers/       # Go HTTP handlers
├── renderer/       # Server-side template rendering logic
├── views/          # HTML templates and components
├── middleware/     # (optional) Middleware logic for Go
├── tsconfig.json   # TypeScript config
├── eslint.config.mjs
└── .lefthook.yml   # Git hook task definitions
```

---

## 🧼 Available Commands

```bash
make lint        # Run all linters
make test        # Run all tests
make build       # Compile TypeScript
make precommit   # Run Lefthook pre-commit tasks manually
```

---

## 🧠 Contributing

Pull requests are welcome! This project is a learning ground for building lightweight, maintainable web apps with Go + Bun.

---
