lint:
	golangci-lint run
	npx eslint . --ext .js
	npx stylelint "**/*.css"
	npx htmlhint "**/*.html"
	npx prettier --check "**/*.{js,css,html}"

format:
	npx prettier --write "**/*.{js,css,html}"

