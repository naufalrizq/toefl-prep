# TOEFL Prep — Backend

Go + Gin API for the TOEFL Prep quiz app. Cookie-based auth (session + CSRF
double-submit), server-side grading, review with POS highlights, and a
progress dashboard. See `../SRS.md` §8 for the endpoint contract.

## Layout

```
cmd/api/                 main (config -> migrate -> seed -> router -> serve)
cmd/genusers/            load-test user generator (N students per VU)
internal/config/         env config (DATABASE_URL, SESSION_SECRET, ...)
internal/database/       pgxpool + embedded goose migrations
internal/http/           envelope {data|error}, errors, BindJSON, constants
internal/http/middleware CORS, AuthRequired, RequireRole, CSRF
internal/httpapi/        router assembly + integration/contract tests
internal/auth/           login/verify/logout, session cookies, login rate limit
internal/questions/      bank: validate, import (AI drafts), CRUD, normalize highlights
internal/exams/          exam templates: validate, publish (bank check)
internal/attempts/       quiz lifecycle: start/snapshot, answer, flag, submit, review
internal/grading/        pure scoring module (unit-tested, 98% coverage)
internal/reporting/      dashboard aggregates (trend, sections, worst POS)
internal/seed/           embedded seed questions + bootstrap users
migrations/              001 users/sessions, 002 questions, 003 exams, 004 attempts
api/openapi.yaml         machine-readable API contract
```

## Run

```bash
# 1. create a Postgres database
./scripts/db-test-setup.sh create

# 2. run (migrates + seeds users student / admin)
DATABASE_URL='postgres://toefl:toefl@localhost:5432/toefl_test' \
SESSION_SECRET='change-me' \
go run ./cmd/api
```

Env vars: `PORT` (8080), `CORS_ORIGINS` (comma list; empty = reflect any),
`SEED_USERS` (true), `LOGIN_RATE_PER_MIN` (10; raise for load testing),
`SESSION_TTL` is fixed at 30 days.

Bootstrap accounts (dev only, rotate in prod):

| username | password | role    |
|----------|----------|---------|
| `student` | `123`   | student |
| `admin`   | `123`   | admin   |

After boot, seed the question bank (8 questions) and publish an exam through
the admin API, or via the integration test:

```bash
# login as admin, then:
curl -c /tmp/c.txt -X POST :8080/api/v1/auth/login -d '{"email":"admin","password":"123"}'
CSRF=$(awk '$6=="toefl_csrf"{print $7}' /tmp/c.txt)
curl -b /tmp/c.txt -H "X-CSRF-Token: $CSRF" -X POST :8080/api/v1/seed
curl -b /tmp/c.txt -H "X-CSRF-Token: $CSRF" -X POST :8080/api/v1/exams \
     -d '{"title":"Structure Basics","section_filters":{"structure":2,"vocabulary":2},"mode":"both","seconds_per_question":20,"total_minutes":15}'
```

## Tests

```bash
go vet ./...
go test ./...                         # unit tests (no DB needed)
go test -race ./...                   # race detector
go test -cover ./...                  # coverage (grading 98%, reporting 83%)

# integration + API contract tests (need a DB):
./scripts/db-test-setup.sh reset
DATABASE_URL='postgres://toefl:toefl@localhost:5432/toefl_test' \
  go test -tags integration ./...
```

The integration tests boot the real router against Postgres and validate
every response body against the JSON schemas embedded in
`internal/httpapi/schemas/` (envelope shape, required fields, enums) — the
executable form of the contract in `api/openapi.yaml`.

## Load testing with k6

```bash
go install go.k6.io/k6@latest

# 1. server up with a high login rate limit
./scripts/db-test-setup.sh reset
DATABASE_URL='postgres://toefl:toefl@localhost:5432/toefl_test' \
SESSION_SECRET='k6' PORT=8080 LOGIN_RATE_PER_MIN=100000 \
go run ./cmd/api

# 2. seed questions + publish an exam (see curl example above)

# 3. one student account per VU (one-active-attempt rule, FR-4.10)
DATABASE_URL='postgres://toefl:toefl@localhost:5432/toefl_test' \
  go run ./cmd/genusers -n 40

# 4. smoke: single-VU full quiz flow (login -> exams -> start -> answer ->
#    submit -> review -> dashboard)
k6 run --vus 1 --iterations 1 k6/scripts/smoke-quiz-flow.js

# 5. load: ramping 0 -> 10 -> 30 -> 0 VUs, 50s, thresholds on p95/p99/errors
k6 run k6/scripts/load-quiz-flow.js
```

Baseline (M1-class machine, local Postgres): ~14 full attempts/s, 137 req/s
at 30 VUs, p(95) ≈ 258ms, p(99) ≈ 762ms, 0% failed, all thresholds green.