# TOEFL Prep — Project Rules & Conventions

**Status:** Draft v1.0
**Date:** 2026-08-19
**Related docs:** [PRD.md](PRD.md) · [SRS.md](SRS.md) · [architecture.md](architecture.md) · [design.md](design.md)

These rules exist so a single developer (you) can drop back into the repo weeks later and know exactly how it works, where things live, and what "done" means. Treat them as the contract for every PR, commit, and page.

---

## 1. Repository layout (monorepo)

```
toefl-prep/
├── PRD.md
├── SRS.md
├── rules.md
├── architecture.md
├── design.md
├── backend/                # Go + Gin API (Railway)
│   ├── cmd/api/main.go
│   ├── internal/
│   │   ├── config/         # env loading, validation
│   │   ├── database/       # pgx pool, migrations runner
│   │   ├── http/           # router, middleware, handlers
│   │   ├── auth/           # sessions, password hashing
│   │   ├── questions/      # bank: repository + service
│   │   ├── exams/          # templates: repository + service
│   │   ├── attempts/       # attempt lifecycle, snapshots
│   │   ├── grading/        # PURE grading engine (no I/O)
│   │   ├── reporting/      # dashboard aggregates
│   │   └── seed/           # seed loader
│   ├── migrations/         # goose migrations (001_.., 002_..)
│   ├── seed/               # questions seed data (JSON)
│   ├── go.mod
│   └── go.sum
├── frontend/               # React + Vite + TS SPA (Vercel)
│   ├── src/
│   │   ├── api/            # fetch client + typed endpoints
│   │   ├── components/     # shared UI (Button, Card, Timer, Badge…)
│   │   ├── features/       # feature-scoped code
│   │   │   ├── auth/
│   │   │   ├── dashboard/
│   │   │   ├── quiz/       # per-question + overall runners
│   │   │   ├── review/
│   │   │   └── admin/      # question & exam editors
│   │   ├── hooks/
│   │   ├── lib/            # formatting, pdf/print helpers
│   │   ├── styles/         # Tailwind + design tokens
│   │   ├── types/          # shared TS types (mirror Go DTOs)
│   │   └── routes.tsx
│   ├── index.html
│   ├── vite.config.ts
│   ├── tailwind.config.ts
│   ├── tsconfig.json
│   └── package.json
├── prompts/                # AI prompt contracts (see §4b)
│   └── question-generator.md
├── design-system/          # ui-ux-pro-max persisted tokens
│   └── toefl-prep/MASTER.md (+ pages/)
└── README.md
```

**Rule:** UI copy stays in English (product language). **Question explanations are written in Bahasa Indonesia** (teaching language) — this is a hard content rule, enforced by the AI prompt contract and admin UX.

---

## 2. Tooling & versions

- **Backend:** Go ≥ 1.22, Gin (router), pgx/v5 (Postgres driver), goose (migrations), golangci-lint.
- **Frontend:** Node ≥ 20, Vite ≥ 5, React ≥ 18, TypeScript strict, Tailwind ≥ 3, React Router ≥ 6, TanStack Query ≥ 5, Recharts ≥ 2, lucide-react.
- **Formatting:** `gofmt` / `goimports` on Go; Prettier on TS/TSX (single config at repo root).
- **Migrations:** goose, forward-only files named `NNNN_name.sql`; never edit an applied migration — add a new one.

---

## 3. Git workflow

- **Commits:** Conventional Commits — `feat:`, `fix:`, `refactor:`, `docs:`, `test:`, `chore:`, `style:`.
- **Branches:** `main` is always deployable. Feature branches `feat/<slug>`, fix branches `fix/<slug>`.
- **PR size:** one logical change per PR; if a PR touches both `frontend/` and `backend/`, it must be one vertical slice (e.g. "add highlight editor end-to-end").
- Never commit secrets, `.env`, or local state. `.gitignore` covers `.env*`, `node_modules`, `dist`, `bin`.

---

## 4. Go conventions

- **Layout:** `cmd/` for entrypoints, `internal/` for everything else (never import `internal/` from outside the module).
- **Package layout per domain:** `repo.go` (data access), `service.go` (use cases), `http.go` (handlers + DTOs), `model.go` (structs). One package per domain from §1.
- **Grading module is sacred:** `internal/grading` contains **zero I/O** — pure functions only. It takes plain structs, returns plain structs. It is the only module with a hard no-dependency rule.
- **Handlers are thin:** parse → validate → call service → map to DTO. No business logic in handlers.
- **Errors:** sentinel errors + `errors.Is`; map to HTTP statuses in a single middleware/helper. Use stable error codes (see SRS §8).
- **Validation:** validate at the edge (handler) AND in services that persist. Reject unknown fields in JSON (`DisallowUnknownFields`).
- **JSON:** fields tagged snake_case to match the DB and frontend. Timestamps as RFC3339.
- **DB:** always parameterized queries (pgx). Wrap multi-step writes in transactions (attempt creation, seed upserts).
- **Offsets for text:** use Go `rune` offsets when validating highlight regions; never byte offsets (UTF-8).
- **Testing:** table-driven tests. `internal/grading` tests live next to the module. Use `testify` sparingly; std-lib `testing` preferred.

---

## 4b. AI import conventions (`prompts/` + `internal/questions`)

- **AI generation happens outside the app.** The prompt contract is `prompts/question-generator.md` — single source of truth; the owner runs it in Claude/Gemini/ChatGPT. The app never calls an LLM and holds no AI API key.
- The import pipeline lives in `internal/questions` (`import.go`): parse pasted JSON → validate per FR-2.5/2.6 → normalize highlights (POS → phrases → offsets) → return per-item draft status.
- Imported JSON is **never trusted**: validate every field; reject unknown keys; normalize phrase-based highlights (never trust LLM offsets); on failure return `validation_failed` with reasons. Nothing is stored until the admin explicitly saves a draft.
- `prompts/question-generator.md` must stay in sync with the validation rules in SRS §3.2/§6.5 — if a rule changes, update the prompt file too.

---

## 5. React / TypeScript conventions

- **TypeScript strict mode** on. No `any` leaks across module boundaries; shared API types in `src/types/` mirror the Go DTOs.
- **Styling:** Tailwind only, using **design tokens** from `design.md`/MASTER.md. No inline `style=` props for colors/spacing. New colors MUST go through the token layer, not raw hex in a component.
- **Icons:** `lucide-react` SVG only. **Never emoji as icons.**
- **State:** server state via TanStack Query (queries + mutations + cache invalidation). Local UI state with `useState`/`useReducer` where needed. No global state store unless the feature demands it (it won't).
- **Routing:** React Router; route-level code splitting via `React.lazy`.
- **Components:** co-located in `features/<name>/`; shared primitives in `components/`. One component per file, PascalCase, default-exported.
- **Accessibility:** all interactive elements have `cursor-pointer`, visible focus styles, `aria-label` where icon-only; keyboard selectable options (A–D / 1–4).
- **Timer display:** use `tabular-nums` and a reserved-width container so the ticking clock never shifts layout (CLS guard).
- **Forms:** controlled inputs; validation messages inline next to fields (never only at the top).

---

## 6. API conventions (envelope)

```json
// success
{ "data": { ... } }
// failure
{ "error": { "code": "validation_failed", "message": "highlight region out of bounds" } }
```

- All endpoints under `/api/v1`.
- Mutating requests require CSRF-safe session (same-site cookie + CSRF token header for mutating calls).
- Every endpoint returns the envelope; no bare arrays or bare error strings.
- Pagination: `?page=1&limit=20` → `data: { items: [], page, limit, total }`.
- Versioning: break nothing; add `/api/v2` if the envelope changes.

---

## 7. Design token usage

- The source of truth is `design.md`; the persisted machine output is `design-system/toefl-prep/MASTER.md` + `pages/*.md` (page-level overrides win).
- Implement tokens as CSS custom properties / Tailwind theme extension (colors, spacing, radii, shadows, fonts).
- Before building a page: check `design-system/toefl-prep/pages/<page>.md` first; if absent, use MASTER.md + design.md.

---

## 8. Testing conventions

- **Backend:** every service has tests; `internal/grading` has ≥ 90% coverage using seed fixtures as inputs. Attempt lifecycle tested end-to-end at service level (start → answer → submit → review).
- **Frontend:** Vitest + React Testing Library. Cover: timer expiry behavior (fake timers), auto-advance, review highlight rendering, print layout presence. No snapshot tests of full pages (brittle).
- **Seed fixtures** double as tests: seeding the DB and grading a known attempt must produce a known score.
- **CI (optional but recommended):** `go test ./...` + `golangci-lint` on backend; `tsc --noEmit` + `vitest run` + `eslint` on frontend.

---

## 9. Environment & configuration

| Env var | Where | Purpose |
|---|---|---|
| `DATABASE_URL` | Railway | Neon Postgres connection string |
| `SESSION_SECRET` | Railway | Signs session cookies |
| `PORT` | Railway | API listen port (Railway provides) |
| `CORS_ORIGINS` | Railway | Comma-separated allowlist (Vercel origin) |
| `VITE_API_URL` | Vercel | API base URL for the SPA (use `/api` relative + Vercel rewrite proxy per ADR-11) |

- Config is validated at boot; the app fails fast if required vars are missing.
- Never put secrets in the frontend bundle. `VITE_*` vars are public by nature — keep only non-secrets there.

---

## 10. Definition of Done

A task is done when all of:

1. Code follows this file's conventions (lint + format clean).
2. Backend: services tested; grading pure + covered; migrations applied cleanly on a fresh DB.
3. Frontend: `tsc --noEmit` clean; feature works in Chrome + Firefox; responsive at 375 / 768 / 1024 / 1440.
4. Accessibility: keyboard navigable, focus visible, contrast ≥ 4.5:1, `prefers-reduced-motion` respected.
5. No emoji-as-icon, no inline raw colors, timer uses tabular numerals.
6. The docs that describe the feature (PRD/SRS) are updated if the behavior changed.
