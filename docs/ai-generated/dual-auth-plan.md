# Plan — dual auth (email/password + Google OAuth) with linked identities

## Context

The app has a stub login page (email + password fields) and a `user` table with no auth columns. The user wants to support **both** email/password **and** Google OAuth, with a single user account potentially having multiple linked login methods (e.g. someone who signed up with Google can later attach a password, or vice versa). They've also flagged that property creation is currently blocked on missing `user_id` (see `docs/ai-generated/todo.md`).

This plan only covers the design/schema/flows. Implementation comes after approval.

---

## Schema changes (one new migration)

New file: `internal/database/migrations/<timestamp>_add_auth.sql` (use `goose create add_auth sql`).

### `user` table — small adjustments
- Add `UNIQUE` constraint on `email` (currently missing — your ERD says UK but the migration doesn't enforce it).
- No new columns needed.

### `identity` table (new)
One row per "way this user can log in." Lets one user link Google + password to the same account.

```sql
CREATE TABLE identity (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    provider TEXT NOT NULL,           -- 'local' | 'google'
    provider_user_id TEXT NOT NULL,   -- email for local, Google's `sub` for google
    password_hash TEXT,                -- NULL unless provider='local'
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE(provider, provider_user_id)
);
CREATE INDEX idx_identity_user_id ON identity(user_id);
```

### `session` table (new)
Server-side sessions (cookie holds an opaque session id only — simpler & safer than JWT for a server-rendered app).

```sql
CREATE TABLE session (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE INDEX idx_session_user_id ON session(user_id);
```

Update `specs/erd.md` with `Identity` and `Session` entities.

---

## New dependencies (go.mod)

- `golang.org/x/crypto/bcrypt` — password hashing
- `golang.org/x/oauth2` + `golang.org/x/oauth2/google` — Google OAuth client

## New env vars (`.env`)

- `GOOGLE_CLIENT_ID`
- `GOOGLE_CLIENT_SECRET`
- `GOOGLE_REDIRECT_URL` (e.g. `http://localhost:8080/auth/google/callback` in dev)
- (Optional later) `SESSION_COOKIE_NAME`, default `tm_session`

You'll create OAuth credentials in Google Cloud Console → APIs & Services → Credentials → "OAuth client ID" → Web application, with the redirect URL above.

---

## Code structure (following the 3-layer pattern in `specs/architecture.md`)

### Repository (sqlc)
- `internal/repository/user.query.sql` — add `CreateUser`, `GetUserByEmail`, `GetUserByID`.
- `internal/repository/identity.query.sql` (new) — `CreateIdentity`, `GetIdentityByProvider(provider, provider_user_id)`, `GetLocalIdentityByUserID(user_id)`.
- `internal/repository/session.query.sql` (new) — `CreateSession`, `GetSession`, `DeleteSession`, `DeleteExpiredSessions`.

After editing, run `sqlc generate`.

### Service — `internal/service/auth-service.go` (new)

Single service exposing:
- `Register(ctx, email, password, name) (User, Session, error)` — creates user + local identity + session; rejects if email already exists.
- `Login(ctx, email, password) (User, Session, error)` — looks up local identity, `bcrypt.CompareHashAndPassword`, creates session.
- `Logout(ctx, sessionID) error`.
- `LoginWithGoogle(ctx, googleSub, email, name, emailVerified bool) (User, Session, error)` — find-or-create flow:
  1. If identity (`google`, sub) exists → use that user.
  2. Else if user with this email exists AND `emailVerified` → **auto-link**: attach a new google identity to that user.
  3. Else → create new user + google identity.
  4. Create session.
- `UserFromSession(ctx, sessionID) (User, error)` — validates session, returns user.

### Handlers
- `internal/handler/login-handler.go` — extend with:
  - `Login(w,r)` — POST /login; on success sets `tm_session` cookie and redirects (Datastar `Datastar-Redirect` header) to `/dashboard`.
  - `Logout(w,r)` — POST /logout; deletes session row + clears cookie.
- `internal/handler/register-handler.go` (new) — `RegisterPage`, `Register` (POST /register).
- `internal/handler/google-auth-handler.go` (new):
  - `GoogleStart(w,r)` — GET /auth/google: generate random state, set short-lived `oauth_state` cookie, redirect to Google's auth URL.
  - `GoogleCallback(w,r)` — GET /auth/google/callback: verify state cookie, exchange code → token, call `https://www.googleapis.com/oauth2/v3/userinfo` for `sub`, `email`, `email_verified`, `name`. Call `service.LoginWithGoogle(...)`. Set session cookie, redirect to `/dashboard`.

### Middleware — `internal/middleware/auth.go` (new)
- `AttachUser(next http.Handler)` — reads `tm_session` cookie, calls `service.UserFromSession`, attaches user to context. Always passes through; missing/invalid session = no user attached.
- `RequireAuth(next http.Handler)` — checks for user in context, redirects to `/login` if missing.

### Context keys — `internal/ctxkeys/keys.go` (new)
- Define `type contextKey int` and `UserKey contextKey = iota`. Provide typed `WithUser(ctx, *User)` / `UserFrom(ctx) (*User, bool)` helpers (avoids naked string keys, matches the directory the architecture doc already anticipates).

### Pages
- `internal/web/page/login-page.templ` — add a "Sign in with Google" button (just an `<a href="/auth/google">`) and a "Register" link.
- `internal/web/page/register-page.templ` (new) — email, password, name fields.

### Routes — `cmd/server/main.go`
Add the auth service + handlers to the wiring. Updated route table:

| Method | Path | Handler | Middleware |
|---|---|---|---|
| GET | /login | `loginHandler.LoginPage` | — |
| POST | /login | `loginHandler.Login` | — |
| POST | /logout | `loginHandler.Logout` | RequireAuth |
| GET | /register | `registerHandler.RegisterPage` | — |
| POST | /register | `registerHandler.Register` | — |
| GET | /auth/google | `googleHandler.GoogleStart` | — |
| GET | /auth/google/callback | `googleHandler.GoogleCallback` | — |
| GET | /dashboard | `tenantHandler.ListTenantPage` | RequireAuth |
| GET | /property/create | `propertyHandler.CreatePropertyPage` | RequireAuth |
| POST | /property/create | `propertyHandler.CreateProperty` | RequireAuth |

`AttachUser` wraps the whole mux; `RequireAuth` wraps individual protected handlers.

### Property handler integration
`internal/handler/property-handler.go` — replace the `user_id` TODO: pull current user from `ctxkeys.UserFrom(r.Context())`. Closes the matching item in `docs/ai-generated/todo.md`.

---

## Security details to keep in mind during implementation

- **Cookie flags**: `HttpOnly`, `SameSite=Lax`, `Path=/`. `Secure` in prod (gate on env).
- **Session id**: server-generated UUID v4 stored in DB. Client never sees a JWT.
- **Session lifetime**: e.g. 30 days; renew on each request or on login only (your call).
- **bcrypt cost**: 10–12.
- **OAuth state**: random 32-byte token, base64-encoded, stored in a short-lived (5 min) HttpOnly cookie; compare-and-delete on callback.
- **Account linking**: only auto-link a Google identity to an existing local user if Google says `email_verified == true`. Otherwise reject with "this email is already registered, please log in with your password first."
- **Password rules**: minimum length only (e.g. 8 chars); skip complexity rules.

---

## Verification

After implementation, end-to-end manual test:

1. `goose up` — schema migrations apply cleanly.
2. `sqlc generate` — generated code compiles.
3. Set `GOOGLE_CLIENT_ID`/`SECRET` in `.env`; create OAuth credentials in Google Cloud.
4. `go run ./cmd/server` (with templ proxy).
5. Browser walkthrough:
   - Visit `/dashboard` while logged out → redirected to `/login`.
   - Register with email+password → redirected to `/dashboard`.
   - Log out → cookie cleared, `/dashboard` again redirects.
   - Log in with same email+password → success.
   - Log out, click "Sign in with Google" → Google consent → back at `/dashboard`. Verify a *single* `user` row with *two* `identity` rows (one `local`, one `google`) — i.e., account linking worked.
   - Create a property → row gets the logged-in user's `user_id`.
6. Tail Postgres: `select * from "user", identity, session;` to sanity check.

---

## Files modified or created (summary)

**Created**
- `internal/database/migrations/<timestamp>_add_auth.sql`
- `internal/repository/identity.query.sql`
- `internal/repository/session.query.sql`
- `internal/service/auth-service.go`
- `internal/handler/register-handler.go`
- `internal/handler/google-auth-handler.go`
- `internal/middleware/auth.go`
- `internal/ctxkeys/keys.go`
- `internal/web/page/register-page.templ`

**Modified**
- `internal/repository/user.query.sql`
- `internal/handler/login-handler.go`
- `internal/handler/property-handler.go`
- `internal/web/page/login-page.templ`
- `cmd/server/main.go`
- `specs/erd.md`
- `go.mod` / `go.sum`
- `.env` (local only — don't commit)
- `docs/ai-generated/todo.md` (close the auth + user_id items)
