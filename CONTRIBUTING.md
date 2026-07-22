# Contributing to `terraform-provider-multicloud`

Thank you for your interest in contributing to `terraform-provider-multicloud`! This document outlines developer setup, testing procedures, resource generator usage, and pull request standards.

---

## 1. Local Development Setup

### Prerequisites
- **Go 1.26+**
- **Terraform CLI 1.8+**
- **Git**

### Clone & Build
```bash
# Clone repository
git clone https://github.com/abmarcum/multi-cloud-provider.git
cd multi-cloud-provider

# Build provider executable binary
go build -o terraform-provider-multicloud .
```

---

## 2. Running Unit & Benchmark Tests

All pull requests must pass 100% of automated unit tests and benchmarks:

```bash
# Execute unit test suite
go test -v ./...

# Execute fuzz testing
go test -fuzz=FuzzSanitizeResourceName ./internal/cloud/... -fuzztime=10s

# Execute benchmark profiling
go test -bench=. ./internal/cloud/...
```

---

## 3. Developer CLI Utilities

### 3.1 API Sync Inspector (`tools/cmd/api-sync`)
Before adding new resources, run the AST API inspector to check for upstream AWS, GCP, and Azure SDK changes:

```bash
go run ./tools/cmd/api-sync --dry-run
```

### 3.2 Resource Code Generator (`tools/cmd/gen-resources`)
Generate standard resource struct boilerplate for new `multicloud_*` primitives:

```bash
go run ./tools/cmd/gen-resources --name my_new_resource
```

### 3.3 Provider HCL & State Converter CLI (`tools/cmd/tf-migrate`)
Convert legacy AWS, GCP, and Azure HCL manifests into unified `multicloud_*` definitions and generate state migration scripts:

```bash
go run ./tools/cmd/tf-migrate --input-dir ./legacy_tf --out-file multicloud_migrated.tf --migrate-state
```

### 3.4 Interactive Terminal TUI (`tools/cmd/tui`)
Inspect local infrastructure schemas and cost estimates interactively:

```bash
go run ./tools/cmd/tui
```

---

## 4. Pull Request Standards

1. **Test Coverage:** Ensure new resources or functions include unit tests (`*_test.go`).
2. **Sensitive Attributes:** Mark secrets, passwords, and tokens as `Sensitive: true` in resource schemas.
3. **Naming Sanitization:** Pass resource names through `cloud.SanitizeResourceName()`.
4. **Documentation:** Update `docs/RESOURCES_REFERENCE.md` and run `go run ./tools/cmd/gen-docs`.
