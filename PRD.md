# TOEFL Prep — Product Requirements Document (PRD)

**Status:** Draft v1.0
**Date:** 2026-08-19
**Owner:** Personal project (single owner = student + admin)
**Related docs:** [SRS.md](SRS.md) · [architecture.md](architecture.md) · [design.md](design.md) · [rules.md](rules.md)

---

## 1. Overview

TOEFL Prep is a **personal, single-owner English learning web app** for TOEFL preparation. The owner uses it to practice TOEFL-style questions under a time constraint, then reviews their results with explanations, exports the exam as a PDF, and tracks whether their scores improve over time.

The product is deliberately narrow: **timed quizzes → review → export → track progress**. There is no social layer, no course content, no game mechanics. The app exists to make *one* habit frictionless: do a timed exam, understand why you got answers wrong, and see the trend.

### Product principles

1. **The exam is the center of attention.** During a quiz there is nothing else on screen. No nav, no chrome, no distractions.
2. **Every wrong answer is a lesson.** The review shows *where* in the sentence the tested concept lives (highlighted in the question) and *why* the correct answer is correct.
3. **Progress is visible at a glance.** One dashboard tells you if you are improving, stagnating, or declining — per test and per section.
4. **Simple beats clever.** Two roles, two quiz modes, one honest scoring model.
5. **The bank grows with AI, generated outside the app.** Questions, answer keys, and Indonesian explanations are drafted by AI (Claude/Gemini/ChatGPT) using the prompt contract in `prompts/question-generator.md` — *in the chat tool, not in the app*. The admin pastes the output into the app's import box; the app validates it and turns it into a reviewable draft. The admin keeps final control — drafts are never silently published.

---

## 2. Goals & Non-Goals

### Goals (v1)

- G1. Take TOEFL-style quizzes (Structure & Written Expression + Vocabulary) with a timer.
- G2. Support **two configurable quiz modes**: per-question countdown (auto-next) and whole-exam countdown (navigable question set).
- G3. Score every attempt and show an instant result summary.
- G4. Provide a **review screen** where question text is **highlighted by part-of-speech** (verb, noun, pronoun, etc.) alongside the explanation, the chosen answer, and the correct answer.
- G5. Export an attempt (questions + your answers + review/explanation) as a **PDF**.
- G6. Track progress on a **dashboard** with a score-over-time chart showing improvement or decline, plus per-section performance.
- G7. Let an **admin** add/edit questions with answer keys, explanations, and highlight annotations, and compose exams from the bank.
- G8. Single-sign-in for exactly two accounts: one student, one admin.

### Non-Goals (explicitly out of scope for v1)

- Multi-user / multi-student support, leaderboards, sharing.
- TOEFL Listening section (audio), Reading Comprehension passages, Speaking.
- Content delivery (lessons, video, courses) — the app is a *quiz* tool only.
- Payments, onboarding flows, email notifications, password reset.
- Spaced-repetition scheduling, streak gamification, badges.
- Mobile native apps (responsive web only).

---

## 3. Personas

| Persona | Role | Needs |
|---|---|---|
| **The Student** (owner) | `student` | Take timed quizzes, see results & explanations, export PDF, track score trend and per-section weaknesses. |
| **The Admin** (same owner) | `admin` | Maintain the question bank (add/edit/delete questions with key + explanation + highlights), define exams, seed the bank, and see all student data. |

Because both roles belong to the same person, the admin surface must be a **low-friction side of the same app** — not a separate "back office" product. Sharing components, tokens, and layout keeps the mental load small.

---

## 4. User Stories

### Student

- As the student, I can log in with my seeded account.
- As the student, I can see a list of published exams and start one.
- As the student, I can choose to run any exam in per-question mode or overall-exam mode (whichever the exam template allows).
- As the student, I get a visible countdown and know when time is about to run out.
- As the student, when time expires I can't keep answering — the attempt is submitted for what I have.
- As the student, immediately after submitting I see my score, correct/wrong counts, and time spent.
- As the student, I can review each question: the question text with **highlighted POS regions**, my answer, the correct answer, and the explanation.
- As the student, I can export the full attempt + review to PDF from my browser.
- As the student, I can open a dashboard showing my score trend (improvement/decline), averages, best score, per-section accuracy, and recent attempts.

### Admin

- As the admin, I can log in with my seeded account.
- As the admin, I can create, edit, delete, and soft-disable questions.
- As the admin, I can specify section (Structure / Vocabulary), question type, question text, options, correct option, difficulty, and explanation.
- As the admin, I can mark **highlight regions** on the question text (select a span of text → assign a part-of-speech) that will appear in the student review.
- As the admin, I can compose an exam: title, quiz mode (per-question and/or overall), time limits, section mix, and how many questions to draw.
- As the admin, I can publish/unpublish exams so only published exams are visible to the student.
- As the admin, I can re-seed the question bank from the repository seed files.

---

## 5. Features

### F1 — Authentication
Two seeded accounts (`student@…`, `admin@…`) with password login. Server issues an httpOnly session cookie. Role is checked on protected endpoints. No registration, no recovery.

### F2 — Question bank (admin)
CRUD over questions. Each question has:

- `section`: `structure` | `vocabulary`
- `type`: `sentence-completion` (Structure) | `vocab-multiple-choice` (Vocabulary) *(extensible)*
- `question_text`: the sentence to complete (Structure) or the prompt (Vocabulary)
- `options`: 4 options (A–D)
- `correct_index`: 0–3
- `explanation`: **in Bahasa Indonesia** — why the answer is correct
- `highlight_regions`: array of `{start, end, pos}` spans on `question_text`
- `difficulty`, `active`, `created_at`, `updated_at`

Deletion is soft (`active = false`) so historical attempts keep their references.

### F2a — AI question import
AI question generation happens **outside the app** — the admin runs the prompt contract in `prompts/question-generator.md` in Claude/Gemini/ChatGPT. The app provides an **import** capability:

- "Import JSON" box in the question editor accepts the AI output (a single question or an array).
- The backend validates every item (schema, bounds, overlap) and normalizes highlight phrases into character offsets, then shows each as a **draft preview** for review.
- The admin reviews/edits drafts, then saves them to the bank. Nothing is auto-published.
- Per-item status on import: valid items become drafts, invalid ones are listed with reasons and don't block the rest.
- Invalid output is rejected with a clear message; the admin can regenerate outside the app and re-import.

### F3 — Exam templates
An exam template defines: `title`, `section_filters`, `question_count` (per section), `shuffle` (default true), `mode`: `per_question` | `overall` | `both` (student can pick), `seconds_per_question` (per-question mode), `total_minutes` (overall mode), `published`.

Starting an attempt **snapshots** the selected questions into the attempt so later edits to the bank never corrupt a finished exam.

### F4 — Quiz taking
Two modes, both timer-driven:

- **Per-question mode:** one question per screen. Timer counts down per question (default 60–90 s). Selecting an option advances to the next question immediately. When the per-question timer hits 0, that question is submitted as unanswered/wrong and the app advances.
- **Overall mode:** a fixed set of questions with a total countdown. Question palette allows jumping between questions, marking for review, and seeing answered/unanswered state. Submitting (manual or auto at timeout) ends the attempt.

The student must confirm submission in overall mode before the attempt closes. Time is tracked client-side for display; the **server records `started_at` / `finished_at`** as the source of truth.

### F5 — Scoring
- Overall score = **percentage of correct answers** (rounded to integer).
- Per-section accuracy (Structure %, Vocabulary %).
- Correct / wrong / unanswered counts.
- Optional future: per-question time analysis.

### F6 — Review (result screen)
After submit, the student sees:

1. **Result summary** — score, correct/wrong/unanswered, time spent, per-section bars.
2. **Per-question review** — question text with **highlighted POS regions**, selected answer vs. correct answer (marked), and the explanation. Wrong answers are visually distinct from correct ones.
3. **Export PDF** button.

### F7 — PDF export
Client-side PDF from a print-optimized review layout (`window.print()` → Save as PDF), so no backend file storage is needed. Fallback path: `react-pdf` if a download file is preferred. PDF must include: header (exam title, date, score), every question with the student's answer vs correct answer, explanation, and highlights.

### F8 — Dashboard (progress tracking)
- Stat cards: exams taken, average score, best score, current trend (vs previous N attempts).
- **Score trend chart** — line/area over time (per attempt) with trend direction; the "range chart" the owner asked for.
- Per-section accuracy chart.
- Recent attempts list with links back to each review.
- Weakness signal: worst-performing POS categories across attempts.

---

## 6. MVP Scope vs Future

### MVP (v1)
F1–F8 with: two sections (Structure + Vocabulary), four-option questions, two quiz modes, PDF via print, AI question import (F2a).

### Future candidates (explicitly parked)
- Listening section (needs audio hosting/CDN + different review UX).
- Reading Comprehension (passage text + multi-question sets).
- Multiple student accounts, sharing, leaderboards.
- Spaced repetition ("which questions do I keep getting wrong?").
- Automated question import from file upload.

---

## 7. Success Metrics

Personal metrics, reviewed weekly by the owner:

- **Volume:** ≥ 3 completed attempts per week.
- **Trend:** dashboard trend line slopes upward over a 30-day window (or at least not consistently down).
- **Consistency:** average score across attempts has low variance after an initial learning curve.
- **Weakness reduction:** per-section accuracy gap narrows; worst POS category improves.
- **Feature use:** ≥ 1 PDF export per attempt (or per week) — proof the review loop is being used.

---

## 8. Key Decisions (recorded)

| # | Decision | Rationale |
|---|---|---|
| D1 | Questions live in the **database** (Neon), seeded from committed **structured files** | Admin edits via UI need a DB; highlight annotations are structured data; seed files give a working bank + test fixtures without manual input. Hardcode-only would force a redeploy per new question. |
| D2 | **Two quiz modes** configurable per exam | Covers both "quick daily drill" (per-question) and "mock TOEFL" (overall). |
| D3 | **Highlight regions stored as JSONB** on the question | Simple, queryable, and exactly matches the review requirement. |
| D4 | **PDF generated client-side** via print stylesheet | Zero backend file storage; Vercel-friendly; the browser already renders the review perfectly. |
| D5 | **Seeded simple auth** (2 accounts) | Exactly two users; full JWT/refresh/SSO machinery is overkill. |
| D6 | **Attempt snapshots questions** | Bank edits never corrupt a finished attempt. |
| D7 | Scoring = **percentage** (0–100) + per-section split | Honest and comparable across attempts of different length. |
| D8 | **AI drafts questions outside the app** via a strict prompt contract (`prompts/question-generator.md`); the app only **imports and validates** the pasted output (no LLM call inside the app) | The LLM work happens in the chat tool the owner already uses; the app stays a pure quiz tool with one simple, well-defined import surface. |
| D9 | **UI copy in English; explanations in Bahasa Indonesia** | Immersion for the exam itself; clear grammar teaching in the student's native language. |

---

## 9. Glossary

- **Attempt** — one completed run of an exam by the student (questions + answers + score + timing).
- **Exam template** — admin-defined configuration for picking questions and timing; publishing makes it startable.
- **Highlight region** — a character span on `question_text` tagged with a part-of-speech category, shown during review.
- **POS category** — part-of-speech label (verb, noun, pronoun, adjective, adverb, preposition, conjunction, determiner, other).
- **Section** — TOEFL section represented in the app: Structure & Written Expression, Vocabulary.
- **Seed** — loading the initial question bank from committed JSON files into the database.
- **AI draft** — a question produced by the LLM (Claude/Gemini/ChatGPT) **outside the app** as structured JSON; imported into the app, validated, and only stored after the admin saves.
- **Highlight normalization** — converting AI-provided POS phrases into validated character-offset regions on `question_text` during import.
