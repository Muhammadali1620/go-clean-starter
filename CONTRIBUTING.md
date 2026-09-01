***
# Backend Architecture & Development Guidelines

This document defines the architectural standards, tech stack, and coding conventions for the project. Every developer must strictly adhere to these rules to maintain code quality, readability, and consistency.

## 1. Tech Stack
* **Language**: Go (Golang)
* **Web Framework**: Echo (`github.com/labstack/echo/v4`)
* **Database**: PostgreSQL
* **ORM**: Bun (`github.com/uptrace/bun`)
* **Decimals**: Shopspring Decimal (`github.com/shopspring/decimal`) for all financial/money fields.

## 2. Project Directory Structure
We follow a modular, layer-based structure. To make file searching easier, **every file must include its layer/purpose in its name**.

```text
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── models/                     # Database entities (e.g., user_model.go, club_model.go)
│   ├── handlers/                   # HTTP/Echo controllers (e.g., user_handler.go)
│   ├── services/                   # Business logic (e.g., user_service.go)
│   ├── repositories/               # Database access (e.g., user_repository.go)
│   ├── dto/                        # Data Transfer Objects for Request/Response (e.g., user_dto.go)
│   └── core/                       # Core configurations, middlewares, utils
│       ├── config/
│       ├── middleware/
│       └── database/
├── pkg/                            # Reusable external packages (e.g., logger, sms_client)
└── docs/                           # Documentation and AI rules
```
## 3. File and Folder Naming Conventions
Folders: All lowercase, single words preferred (e.g., handlers, services).

Files: snake_case. Always append the component's role to the filename.

✅ DO: user_model.go, club_service.go, booking_repository.go, payment_handler.go.

❌ DON'T: user.go, club.go, general.go (Too ambiguous for IDE file search).

## 4. Architectural Layers & Strict Separation of Concerns
# A. Handlers Layer (internal/handlers/)
Purpose: Handle HTTP requests, parse inputs, and return HTTP responses.

Rules:

Must use echo.Context.

Must parse requests into DTOs (Data Transfer Objects) and validate them.

NEVER contain business logic.

NEVER import or call bun.DB or models directly. Only call the Service layer.

# B. Service Layer (internal/services/)
Purpose: Contain 100% of the application's business logic.

Rules:

Must orchestrate calls between different repositories.

NEVER import echo or know anything about HTTP (no echo.Context, no HTTP status codes).

NEVER write SQL queries here. Use Repositories for data fetching.

# C. Repository Layer (internal/repositories/)
Purpose: Direct interaction with the PostgreSQL database using Bun ORM.

Rules:

Contains all CRUD operations, SQL queries, and DB transactions.

NEVER contain business logic or validation (e.g., do not check if a user has enough balance here; check it in the Service, then update in the Repository).

## 5. Coding Style & Conventions
Receiver Naming
Use short, 1-2 letter variable names that reflect the struct type.

✅ DO: func (h *UserHandler) Create(c echo.Context)

✅ DO: func (s *UserService) GetByID(ctx context.Context, id int64)

✅ DO: func (r *UserRepository) FindAll(ctx context.Context)

❌ DON'T: func (this *UserService)... or func (self *UserService)... (Anti-pattern in Go).

Variable and Struct Naming
Use camelCase for private variables/functions and PascalCase for public ones.

Avoid stuttering.

✅ DO: user.Service (when importing the package).

❌ DON'T: user.UserService.

Error Handling
Never ignore errors (_).

Wrap errors with context before returning them up the chain using fmt.Errorf("failed to fetch user: %w", err).

Handlers should map specific domain errors to proper HTTP status codes (e.g., 400, 403, 404, 500).

Comments and Documentation
Use standard Godoc format.

All code, comments, and docstrings MUST be written in English. No Cyrillic or other languages in the source code.
***
