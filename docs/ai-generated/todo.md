# TODO — deferred concerns

## Remaining for the auth + user feature

### Register UI flow not finished
POST /register handler is wired and the service is solid, but the UI side is half-done:
- `http.Error` placeholders still in `auth-handler.go` Register branches (lines 61, 69) — replace with `sse.PatchSignals` for inline error display, matching the Login pattern that's already working.
- No dedicated register form. The login page has a "Register" button but it doesn't submit to `/register` (it's not even `type="submit"`). Options: build a separate `register-page.templ` with its own form posting to `/register`, **or** make the same login form switch action depending on which button was clicked.
- Add `registerError` signal + `<p data-show data-text>` display element in whichever templ ends up holding the register form (same shape as the working `loginError` pattern).

### Logout UI button
Logout works (POST `/logout` clears the cookie and redirects), but there's no way for a user to trigger it from the UI yet — they currently have to delete the cookie in DevTools. Add a logout button somewhere visible after login (dashboard for now, or a shared nav component when one exists). Trigger via Datastar `@post` to `/logout`.

### Property handler — pull `user_id` from authenticated context
`CreateProperty` in `internal/handler/property-handler.go` still has TODOs at lines 39-41. Once the property service/sqlc work lands (see "Property: sqlc query" below), the user_id integration is straightforward:
- `user, _ := ctxkeys.UserFrom(r.Context())` (route is RequireAuth-protected, so the user will be present)
- Pass `user.ID` as the `user_id` field to the service.

### Google OAuth flow (LoginWithGoogle)
Service stub exists in `auth-service.go`; everything else is pending. See `docs/ai-generated/dual-auth-plan.md` for the full plan. Concrete tasks:
- Implement `LoginWithGoogle` in `auth-service.go` — find-or-create with `emailVerified`-gated account linking.
- New handler `internal/handler/google-auth-handler.go` with `GoogleStart` (generate state, redirect to Google) and `GoogleCallback` (verify state, exchange code, call service, set cookie, redirect).
- Routes: `GET /auth/google`, `GET /auth/google/callback` (both unprotected).
- New deps: `golang.org/x/oauth2`, `golang.org/x/oauth2/google`.
- Env vars: `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `GOOGLE_REDIRECT_URL` (e.g. `http://localhost:8080/auth/google/callback` in dev).
- Google Cloud Console: create OAuth Web Application client, configure consent screen (External + Testing while iterating), add own gmail as a test user.
- "Sign in with Google" button on the login page — a simple `<a href="/auth/google">` works.

### Small auth cleanups
- **`RequireAuth` signature inconsistency.** Currently `RequireAuth` takes/returns `http.HandlerFunc`, while `AttachUser` uses the standard `http.Handler`. Pick one: either revert `RequireAuth` to `func(http.Handler) http.Handler` and add a small `protect()` helper at call sites, or flip `AttachUser` to match. Right now the package has two middleware shapes.
- **Extract `middleware.AttachUser(authService)` to a named variable** in `main.go` for readability — the `(svc)(mux)` factory-call double-paren is hard to scan.
- **Stale comments in `auth-handler.go`**: `// TODO: show inline err message`, `// TODO: prompt email already taken`, etc. — clear these out as the Register UI lands.

### (Stretch) Constant-time defense for Google
Login already does a dummy bcrypt compare to avoid user-enumeration timing attacks. The Google path's "no such user, no auto-link possible" branch should similarly avoid leaking timing differences. Probably negligible in practice (OAuth flow is much longer than DB-side timing differences) — defer until Google flow is functional.

---

## Other deferred concerns (non-auth)

### pgx type conversion lives in the handler
Form input is a string; the repo expects `pgtype.Numeric` for rent, `pgtype.UUID` for IDs. Conversion currently happens in the handler. Future move: introduce a domain input struct (e.g. `service.CreatePropertyInput{Name string, RentAmount decimal.Decimal}`) so the service is decoupled from `pgtype.*`. Conversion happens at the service ↔ repo boundary.

### Routes still live in `cmd/server/main.go`
Architecture doc (`specs/architecture.md`) calls for `internal/routes/routes.go`. Acceptable while route count is small. Refactor trigger: ≥3-4 feature routes or main.go's route registration crosses ~15 lines — we're close.

### Datastar response strategy — half decided
Auth flows landed on signal-patch (errors) + SSE redirect (success). The `CreateProperty` handler still doesn't write a response. Decide: render a success partial that swaps into the form container, or redirect to a property list once that view exists.

### Validation layer — barely used
`internal/validation/` now exists with `RegisterInput` only. Extract more helpers as additional handlers need validation (`CreateProperty` is the obvious next caller). Architecture doc anticipates this directory.

### Property: sqlc query `CreateProperty` not written yet
Repository only has `ListProperties`. Need to add `CreateProperty` to `internal/repository/property.query.sql` and run `sqlc generate` before the service does anything useful. Decisions still open:
- UUID generation: SQL-side via `gen_random_uuid()` (no Go dep) vs Go-side via a uuid lib.
- `:one` (return inserted row, useful for success UI) vs `:exec`.

### Audit copied datastarUI components for pre-v1 syntax
Already hit the form-submit issue once (`data-on-submit` vs `data-on:submit` — see `datastarUI-old-datastar-version-issue.md`). The datastarUI source we're copy-pasting from is older than our Datastar runtime. One-time sweep of `internal/web/component/` for `data-on-`, `data-indicator-`, `data-action-`, `data-store`. Going forward: smoke-test every newly-copied component's interactive behavior — failures are silent.
