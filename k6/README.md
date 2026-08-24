# k6 load tests

Scripts hit the full student flow: login → exams → start → answer → submit → review → dashboard.

## Prerequisites

The scripts log in as dedicated per-VU accounts (`student1@toefl.dev`, `student2@toefl.dev`, …)
with password `student1234`. The seeded dev accounts (`student`/`123`) do **not** match —
create the load-test users first:

```sh
DATABASE_URL=postgres://toefl:toefl@localhost:5432/toefl_test \
  go run ./backend/cmd/genusers -n 20
```

Also required: a **published** exam template (the smoke script assumes id `1`; override
with `EXAM_ID`). Run migrations first so the DB is up to date.

## Smoke

```sh
k6 run k6/scripts/smoke-quiz-flow.js
```

## Load

Each VU uses its own account (one active attempt per user, FR-4.10):

```sh
BASE_URL=http://localhost:8080 k6 run k6/scripts/load-quiz-flow.js
```

## Overridable env vars

| Var | Default | Purpose |
|---|---|---|
| `BASE_URL` | `http://localhost:8080` | API base |
| `EMAIL` / `PASSWORD` | `student1@toefl.dev` / `student1234` | Smoke-test credentials |
| `PASSWORD` | `student1234` | Load test password for all VU accounts |
| `EXAM_ID` | `1` | Published exam template to run |
