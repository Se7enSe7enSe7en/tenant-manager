# TODO — deferred concerns

Things noted while wiring up the property creation flow that we agreed to handle later.

## Auth / `user_id` FK

The `property` table has `user_id UUID NOT NULL` with an FK to `"user"`. We have no auth yet, so the `CreateProperty` handler currently can't supply a real user.

- Short-term plan: simple auth implementation comes next (after the create-property flow is logging end-to-end).
- Until then: handler/service calls into `CreateProperty` should be left with a TODO comment where `user_id` would be populated.
- Once auth lands: source `user_id` from session / request context (likely via `internal/ctxkeys`).

## pgx type conversion lives in the handler (for now)

Form input is a string; the repo expects `pgtype.Numeric` (and `pgtype.UUID` for `user_id`). We're doing the string → pgtype conversion inside the handler for now.

- Tradeoff: keeps the service signature tied to repo types, which leaks the data-access layer up into the handler.
- Future move: introduce a domain input struct (e.g., `service.CreatePropertyInput{Name string, RentAmount decimal.Decimal}`) so the service is decoupled from `pgtype.*`. Conversion to `pgtype.*` then happens at the service ↔ repo boundary.

## Routes still live in `cmd/server/main.go`

Architecture doc (`specs/architecture.md`) calls for centralizing routes in `internal/routes/routes.go`, but they're currently inline in `main.go`.

- Acceptable while there are only a handful of routes.
- Refactor trigger: when we have ~3–4 features wired up, or when `main.go` route registration crosses ~15 lines, move to `internal/routes/SetupRoutes(app)`.

## Datastar response strategy not yet decided

`CreateProperty` handler currently only logs — no response written. Datastar expects an HTML fragment, an SSE stream, or a `Datastar-*` header (e.g., `Datastar-Redirect`) rather than a 302.

- Decide once we move past logging: probably render a small success partial that swaps into the form container, or redirect to a property list page once that exists.

## Validation layer unused

`internal/validation/` is referenced by the architecture doc but doesn't exist in the tree yet. Inline validation in the handler is fine for now; extract to `internal/validation/` once we have ≥2 handlers doing similar checks.

## Form data shape from Datastar — verify on first run

With `contentType: 'form'` (the default in `form.Form`), Datastar should submit the native `<form>` as `FormData`, so `r.FormValue("name")` and `r.FormValue("rent_amount")` will work. If the logs come back empty, check the browser Network tab — Datastar may be sending signal-pathed keys (e.g., `property_form.name`) instead, in which case the input components need `FormID` set so `data-bind` is applied.

## sqlc query for `CreateProperty` not written yet

Repository only has `ListProperties`. Need to add `CreateProperty` to `internal/repository/property.query.sql` and run `sqlc generate` before the service can do real work.

- Decision pending: generate UUIDs in SQL via `gen_random_uuid()` (no Go dep) or in Go via `github.com/google/uuid`.
- Decision pending: `:one` (return inserted row, useful for success UI) vs `:exec`.

## Audit copied datastarUI components for pre-v1 Datastar syntax

We hit one of these already on the property form (see `datastarUI-old-datastar-version-issue.md`). The datastarUI repo we're copy-pasting from uses an older Datastar release, so any component we pull in could silently misbehave under v1.0.1.

- One-time sweep: grep `internal/web/component/` for `data-on-`, `data-indicator-`, `data-action-`, `data-store` and convert to v1 syntax.
- Going forward: smoke-test every newly copied component's interactive behavior (click, submit, bind) before assuming it works — failures are silent.

## Pre-existing: dev user seed

Once auth is deferred-but-not-blocking, we'll likely need a seed migration that inserts a default dev user so FK-bound inserts work in local dev. Open question whether that lives in `migrations/` (goose) or a separate seed script.
