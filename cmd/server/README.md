# Server - Clean Architecture

This is the **main production server** using Clean Architecture / Hexagonal Architecture.

## Architecture

- **Domain Layer**: Core business logic with no external dependencies
- **Application Layer**: Use cases and DTOs
- **Infrastructure Layer**: External implementations (database, JWT, Casbin)
- **Adapter Layer**: gRPC handlers

## Run

```bash
# Development
go run ./cmd/server

# Build
go build -o server ./cmd/server

# Run binary
./server
```

## Features

✅ Clean Architecture with proper dependency inversion
✅ Domain-driven design
✅ Testable with mocked interfaces
✅ Scalable and maintainable
✅ New infrastructure implementations

## Environment Variables

See `.env.example` or `internal/infrastructure/config/config_loader.go` for configuration.

## Documentation

- Architecture: `../../ARCHITECTURE_NEW.md`
- Migration Guide: `../../MIGRATION_GUIDE.md`
- Refactoring Summary: `../../REFACTORING_SUMMARY.md`

---

**Note:** This is the recommended version for production use.
For the legacy version, see `../legacy/`

