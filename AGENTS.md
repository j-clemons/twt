# Agent Guidelines for twt

## Build/Test Commands
- `go build` - Build the project
- `go test ./...` - Run all tests
- `go test ./internal/command` - Run tests for specific package
- `go test -run TestValidate` - Run specific test function
- `go run main.go` - Run the CLI tool
- `gofmt -w .` - Format code

## Code Style Guidelines
- Use tabs for indentation (Go standard)
- Package names: lowercase, single word (e.g., `command`, `git`, `tmux`)
- Constants: UPPER_SNAKE_CASE (e.g., `NEW_DIR_PERM`)
- Variables: camelCase (e.g., `sessionName`, `hasCommonSession`)
- Functions: PascalCase for exported, camelCase for unexported
- Import grouping: standard library, then third-party, then local packages
- Error handling: explicit checks with early returns
- Use `github.com/fatih/color` for colored output (Red, Green, Yellow, Cyan)
- Use `github.com/spf13/cobra` for CLI commands
- Test files: `*_test.go` with `package_test` suffix for external tests

## Project Structure
- `cmd/` - CLI command definitions using Cobra
- `internal/` - Private packages (command, git, tmux, utils, checks)
- Main entry point: `main.go` calls `cmd.Execute()`
- Module: `github.com/j-clemons/twt`
