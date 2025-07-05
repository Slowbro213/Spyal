# ALSpy 🇦🇱🔍

A modern, minimalist, Albanian-language version of the Spyfall game — built with
Go, TypeScript, HTML, CSS, and powered by Bun for the frontend.

---

## 🔧 Tech Stack

- **Backend**: Go (`net/http`, server-side templating)
- **Frontend**: Vanilla JS, TypeScript, HTML, CSS (no framework)
- **Build Tool**: [bun][bun-url]
- **Linting**: eslint, prettier, stylelint
- **Testing**: bun test runner
- **Git Hooks**: [lefthook][lefthook-url] (for linting and tests pre-commit/push)

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

The app will be available at:  
`http://localhost:8080`

---

## 🛡️ Pre-commit and Pre-push Hooks with Lefthook

This project uses [lefthook][lefthook-url]
to automate tasks before committing and pushing.

### Features

- Run **eslint**, **stylelint**, and **go vet/test** before each commit
- Run **full test suite** before pushing

### Setup

1. Install Lefthook

   ```bash
   go install github.com/evilmartians/lefthook@latest
   ```

2. Enable git hooks

   ```bash
   lefthook install
   ```

3. Lefthook will now:

   1. Lint your code before commits  
   2. Run tests before pushes  
   3. Ensure quality checks are enforced consistently

---

## 📁 Folder Structure

```text
.
├── assets/         # Static assets (CSS, images, TS/JS)
├── cmd/server/     # Go HTTP server entrypoint
├── handlers/       # Go HTTP handlers
├── renderer/       # Server-side template rendering logic
├── views/          # HTML templates and components
├── middleware/     # Middleware logic for Go
├── tsconfig.json   # TypeScript config
├── eslint.config.mjs
└── .lefthook.yml   # Git hook task definitions
```

---

## 🧼 Available Commands

```bash
make lint
make test
make build
make precommit
```

---

## 🧠 Contributing

Pull requests are welcome!  
This project is a learning ground for building lightweight,
maintainable web apps with Go and Bun.

---

[bun-url]: https://bun.sh/  
[lefthook-url]: https://github.com/evilmartians/lefthook
