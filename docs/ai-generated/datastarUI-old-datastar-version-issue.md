# datastarUI uses pre-v1 Datastar attribute syntax

## What happened

Submitting the create-property form did nothing in Datastar terms — the form fell back to a native browser GET (`/property/create?name=...&rent_amount=...`) instead of firing Datastar's `@post`. No console errors, no failed network requests, Datastar itself was loaded and working (a `data-on:click` counter on the same page reacted fine).

## Root cause

The `form.Form` component copied from the [datastarUI repo](https://datastar-ui.com/docs) emits **pre-v1 Datastar attribute syntax**. Specifically, in `internal/web/component/form/form.templ`:

```go
formAttrs["data-on-submit"] = submitHandler          // ← pre-v1
formAttrs["data-indicator-fetching"] = ""            // ← pre-v1
```

Datastar v1 uses **colon-separated** event attributes:

```go
formAttrs["data-on:submit"] = submitHandler          // ← v1 correct
```

And `data-indicator` in v1 takes a signal name as its **value**, not a hyphen-suffixed attribute name.

Because the syntax doesn't match, Datastar v1 silently skips the attribute — no binding is attached, the form's default native submit takes over, and you get a GET to the current URL with form fields as query params.

## Why this is a recurring risk

- We're copy-pasting components from the datastarUI repo as needed (see `specs/main.md` → "UI library: DatastarUI").
- That repo appears to track an older Datastar release. Its components mix v1-valid syntax (e.g. `data-bind` in `input.templ`) with pre-v1 syntax (e.g. `data-on-submit` in `form.templ`).
- The Datastar runtime we load is **v1.0.1** (see `internal/web/page/base.templ`), so any pre-v1 attribute in a copied component will fail silently.

## How to spot it

Symptoms when a component has stale syntax:
- The component's interactive behavior just doesn't happen.
- No console errors.
- Network tab is empty (or shows a native browser navigation instead of an `@post`/`@get` fetch).
- DevTools → Elements shows the `data-*` attribute is present in the DOM.

## What to do when copying a new component

1. Skim the component's `.templ` source for `data-` attributes.
2. Convert any `data-on-{event}` → `data-on:{event}`.
3. Convert any `data-indicator-{name}=""` → `data-indicator="${name}"` (or whatever the v1 spec actually wants — verify against [data-star.dev/reference/attributes](https://data-star.dev/reference/attributes)).
4. Watch out for other v0-era patterns: `data-action-*`, `data-store`, etc. When in doubt, check the v1 attribute reference.
5. Smoke-test the component (click/submit/etc.) before assuming it works.

## Components confirmed patched

- `internal/web/component/form/form.templ` — `data-on-submit` → `data-on:submit`. (The `data-indicator-fetching` line may also need fixing — flag if a "fetching" indicator doesn't render.)

## Components not yet audited

The rest of `internal/web/component/` was copied from the same source and may have the same issue. See `docs/ai-generated/todo.md` for the audit task.
