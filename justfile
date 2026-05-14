BINARY := "lazycosmos"
GO := env_var_or_default("GO", "go")
VERSION := env_var_or_default("VERSION", "dev")
COMMIT := `git rev-parse --short HEAD 2>/dev/null || echo none`
DATE := `date -u +%Y-%m-%dT%H:%M:%SZ`
LDFLAGS := "-X main.version=" + VERSION + " -X main.commit=" + COMMIT + " -X main.date=" + DATE
GOLANGCI_LINT_VERSION := "v2.12.2"
GOLANGCI_LINT := env_var_or_default("GOLANGCI_LINT", "go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@" + GOLANGCI_LINT_VERSION)

run:
    {{GO}} run ./cmd/lazycosmos

test:
    {{GO}} test ./...

build:
    mkdir -p bin
    {{GO}} build -ldflags "{{LDFLAGS}}" -o bin/{{BINARY}} ./cmd/lazycosmos

lint:
    {{GOLANGCI_LINT}} run ./...

fix:
    {{GOLANGCI_LINT}} run --fix ./...

clean:
    rm -rf bin dist build coverage.out
