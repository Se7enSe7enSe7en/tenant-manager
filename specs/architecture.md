# Architecture

## Project Structure

```
project/
├── cmd/server/          # Main web server entry point
├── internal/            # Private application code
│   ├── handler/         # HTTP handlers (presentation layer)
│   ├── service/         # Business logic
│   ├── repository/      # Data access layer (sqlc will generated files based on .query.sql files)
│   ├── middleware/      # HTTP middleware
│   ├── routes/          # Route definitions
│   ├── database/        # Database setup & migrations (model and repository layer is handled by sqlc)
│   │   ├── generated/   # sqlc generated files from .query.sql files
│   │   └── migrations/  # migration files using goose
│   ├── validation/      # Input validators
│   ├── config/          # Configuration management
│   ├── constants/       # Constant variables
│   ├── ctxkeys/         # Context key definitions
│   ├── logger/          # Structured logging
│   └── utils/           # Shared utilities
└── web/                 # Front end (templ templates)
    ├── static/          # Public files
    │   └── assets/      # Static files (CSS, JS, images)
    ├── components/      # Reusable components
    ├── layout/          # Page layouts
    └── pages/           # Full page templates
```

## 3 Layer architecture

```
HTTP Request
    ↓
Handler (presentation)
    ↓
Service (business logic)
    ↓
Repository (data access)
    ↓
Database
```

### 1. Handler Layer

Handles HTTP concerns:

- Parse request
- Validate input
- Call service
- Render response

```go
func (h *BlogHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
    posts, err := h.blogService.Posts()
    if err != nil {
        http.Error(w, "Failed to load posts", 500)
        return
    }
    ui.Render(w, r, pages.BlogList(posts))
}
```

### 2. Service Layer

Contains business logic:

- Coordinate operations
- Enforce business rules
- Manage transactions

```go
type AuthService struct {
    userRepo  UserRepository
    jwtSecret string
}

func (s *AuthService) Login(email, password string) (string, error) {
    user, err := s.userRepo.GetByEmail(email)
    // Authentication logic...
}
```

### 3. Repository Layer

Handles data access:

- Database queries
- Data mapping
- No business logic

```sql
-- name: ListTenants :many
SELECT * FROM tenant;
```

The repository layer will be handled by sqlc, generated from \*.query.sql files, the generated files are in internal/database/generated/

## Dependency Injection

Constructor-based injection. Services get dependencies via New\*() functions:

```go
// Handler needs service
func NewBlogHandler(service *BlogService) *BlogHandler {
    return &BlogHandler{service: service}
}

// Service needs repository
func NewAuthService(repo UserRepository, jwtSecret string) *AuthService {
    return &AuthService{
        userRepo:  repo,
        jwtSecret: jwtSecret,
    }
}
```

- Why do we need Dependency Injection? couldn't we just do imports?
  - it provides clear requirements for what is needed to run the thing (service, handler, repository etc.)
  - testability, for example, we can provide a mock repo for the service in our unit tests

## Routing

Routes are centralized in /internal/routes:

```go
func SetupRoutes(app *app.App) http.Handler {
    mux := http.NewServeMux()

    // Blog routes
    blog := handler.NewBlogHandler(service.NewBlogService("content"))
    mux.HandleFunc("GET /blog", blog.ListPosts)
    mux.HandleFunc("GET /blog/{slug}", blog.ShowPost)

    return mux
}
```

Using stdlib net/http router - no third-party packages.

## Adding New Features

Example: Adding a “products” feature

1. Create model in /internal/model/product.go
2. Create repository in /internal/repository/product.go
3. Create service in /internal/service/product.go
4. Create handler in /internal/handler/product.go
5. Create UI in /web/pages/products.templ
6. Add routes in /internal/routes/routes.go

```
internal/
├── model/product.go
├── repository/product.go
├── service/product.go
└── handler/product.go

web/
└── pages/products.templ
```

## Best Practices

- Keep /internal private - All application code goes here
- Separate concerns - Handler → Service → Repository (no shortcuts)
- Tests alongside code - user.go + user_test.go in same directory
- No circular dependencies - Handlers depend on services, services depend on repositories (not the other way)

### When to Use Interfaces

- Use interfaces where change is likely - Database layer (SQLite → Postgres), external APIs (Polar → Stripe), storage (MinIO → S3)
- Use concrete types where stable - Services and handlers rarely need swapping
  Producer-side interfaces - Define at repository layer, not at every consumer (avoids duplication, easier testing)

### Database Layer

- Minimal abstraction - We use sqlc and write raw SQL, we then use the sqlc generated go code for type safety and as our query builder
- Repository pattern - Clean interface makes switching databases trivial

### Security

- Keep sensitive logic in /internal (cannot be imported externally)
- Use environment variables for secrets (see Configuration)
- Validate all inputs in handlers
- Never commit .env files

## References and Credits

- this architecture document is heavily inspired by Axel Adrian's goilerplate doc: https://goilerplate.com/docs/architecture
