# TOEFL Prep — Frontend

React + Vite + TypeScript SPA. Design system: `../design.md` (authoritative)
+ `../design-system/toefl-prep/`. Layout mirrors `../rules.md` §1.

## Stack

- React 18, Vite 5, TypeScript strict
- Tailwind 3 (design tokens as CSS variables in `src/styles/index.css`)
- React Router 6 (route-level code splitting), TanStack Query 5
- Recharts (dashboard only, lazy), lucide-react icons

## Run

```bash
npm install
npm run dev            # http://localhost:5173, proxies /api -> :8080
```

Backend must be running on :8080 (see `../backend/README.md`). Dev proxies
`/api` so cookies stay same-origin — no CORS work needed locally.

## Structure

```
src/
  api/              # fetch client: envelope unwrap, CSRF header, credentials
  components/       # shared primitives (Button, Card, Timer, HighlightedText…)
  features/
    auth/           # login
    dashboard/      # stats, area chart, section bars, recent attempts
    exams/          # published exam list + start
    quiz/           # per-question + overall runners (core)
    review/         # POS-highlight review + PDF print layout
    history/        # attempts list
    admin/          # question editor (+ highlight editor, AI import), exam editor
  hooks/useAuth.tsx # auth context + guards
  lib/              # api client, formatters
  styles/           # Tailwind + tokens + print stylesheet
  types/            # mirrors Go DTOs (see ../backend/api/openapi.yaml)
  routes.tsx
```

## Design notes

- **POS highlights are the signature**: `HighlightedText` renders each region
  bold + POS color + distinct underline style (colorblind-safe pairing), with
  `aria-label`/tooltip and an on-hover chip.
- **Timer** uses `tabular-nums` in a reserved-width container (no CLS); amber
  pulse only inside `prefers-reduced-motion` scope. Question expiry marks the
  item unanswered and advances (per-question mode).
- **Print/PDF**: review has a first-class print stylesheet — hidden chrome,
  header block, `break-inside: avoid` per card, `print-color-adjust: exact`.
- All colors come from the token layer only — no raw hex in components.
- UI copy in English; question explanations in Bahasa Indonesia (content rule).

## Tests

```bash
npm run typecheck     # tsc --noEmit
npm test              # vitest run (RTL)
npm run build         # typecheck + production build
```

Covered: `HighlightedText` region splitting/clamping/labels; `QuizPerQuestion`
auto-advance, A–D keyboard, timer-expiry → unanswered + advance, finish
behaviour, flag toggle.

## Vercel (deploy)

Set the API base at build time if not proxied: in production the app talks to
`/api/v1` on the same origin; Vercel rewrites `/api/*` to the backend, or set
`VITE_API_BASE` if you point at a separate API host.