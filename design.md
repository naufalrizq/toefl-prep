# TOEFL Prep — Design (UI/UX)

**Status:** Draft v1.0
**Date:** 2026-08-19
**Related docs:** [PRD.md](PRD.md) · [SRS.md](SRS.md) · [architecture.md](architecture.md) · [rules.md](rules.md)
**Token source:** `design-system/toefl-prep/MASTER.md` + `pages/*.md`. **This file is the human-authored authority** — where it differs from the machine-generated MASTER.md, this file wins and the overrides are listed in §1.3.

---

## 1. Design direction

### 1.1 The concept: "The Exam Paper"

This is a **personal TOEFL practice room**, not a gamified quiz app and not a marketing site. The one job: sit down, face a timed exam, understand your mistakes, get measurably better.

The visual metaphor is **a clean exam paper on a quiet desk**:

- **Warm, quiet surface** — paper-toned background, no gradients, no glassmorphism, no decorative blobs.
- **Serif question text** — English sentences are set in a serif face so they read like a printed test, not a dashboard.
- **The highlight is the hero** — the only "colorful" thing in the whole app is the part-of-speech highlight system in review. It is the product's signature and it stays loud because everything else stays quiet.
- **Calm under pressure** — during an exam the screen is stripped to the question, the timer, and the options. Nothing competes.

### 1.2 What this is not (anti-patterns we refuse)

- ❌ Playful claymorphism/kid-education styling — this is an adult, self-directed study tool.
- ❌ "Exaggerated minimalism" — oversized type and loud contrast belong to fashion portfolios, not a reading-focused exam screen.
- ❌ Dashboard-with-everything — the quiz screen shows only the quiz.
- ❌ Emoji as icons, glassmorphism cards, gradient buttons, confetti on results (results are informative, not celebratory).
- ❌ Color as the *only* signal — every color (correct/wrong/POS) is paired with a shape, label, or text.

### 1.2b Language
- **UI copy (nav, buttons, labels): English** — exam-day immersion.
- **Question text & options: English** (they are TOEFL content).
- **Explanations (`explanation`): Bahasa Indonesia** — the teaching voice speaks the learner's language so grammar points land fast. This is enforced by the AI prompt contract (`prompts/question-generator.md`) and by the admin editor's guidance. Never write an English explanation.

### 1.3 Overrides vs `design-system/MASTER.md`

| Token | MASTER.md | This design | Why |
|---|---|---|---|
| Body font | Raleway | **Inter** | Raleway is thin and stylized; Inter is neutral, highly legible, ships tabular numerals for the timer. |
| Question/exam text | — (none) | **Lora** (serif) | Reading rhythm for English exam sentences; keeps MASTER's serif as the *reading* face. |
| Background | `#EFF6FF` (blue tint) | `#FAF9F6` (warm paper) | Paper metaphor; less screen-glow during long sessions. |
| Primary | `#2563EB` | `#2563EB` (kept) | Blue-600 = calm, trustworthy, exam-appropriate. |
| Accent | `#F59E0B` | `#F59E0B` (timer only) | Amber is reserved for time pressure — one meaning, one place. |

---

## 2. Visual identity

### 2.1 Color tokens

```css
:root {
  /* surfaces */
  --paper:        #FAF9F6;   /* app background */
  --card:         #FFFFFF;
  --card-muted:   #F4F2ED;   /* secondary surfaces, zebra, wells */
  --ink:          #17181C;   /* primary text (warm near-black) */
  --ink-muted:    #5C5F66;   /* secondary text */
  --ink-faint:    #8B8E95;   /* placeholders, disabled */
  --border:       #E4E1DA;   /* hairlines */
  --border-strong:#C9C4BA;

  /* brand + action */
  --primary:      #2563EB;   /* links, primary button, selected */
  --on-primary:   #FFFFFF;
  --primary-soft: #DBEAFE;   /* selected option fill */

  /* status */
  --success:      #16A34A;   /* correct */
  --success-soft: #DCFCE7;
  --danger:       #DC2626;   /* wrong, expired, destructive */
  --danger-soft:  #FEE2E2;
  --warning:      #D97706;   /* timer < 10s */

  /* POS highlights (see §3) */
  --pos-verb:        #6D28D9;  /* violet-700 */
  --pos-noun:        #1D4ED8;  /* blue-700 */
  --pos-pronoun:     #BE123C;  /* rose-700 */
  --pos-adjective:   #15803D;  /* green-700 */
  --pos-adverb:      #B45309;  /* amber-700 */
  --pos-preposition: #0F766E;  /* teal-700 */
  --pos-conjunction: #4338CA;  /* indigo-700 */
  --pos-determiner:  #475569;  /* slate-600 */
  --pos-other:       #6B7280;  /* gray-600 */
}
```

Contrast: all POS text colors are 700-level shades — ≥ 7:1 on paper/card, ≥ 4.5:1 even on `--card-muted`. Status colors only used as text/border on light fills (never as bare background with dark text).

### 2.2 Typography

| Role | Face | Weight / Size |
|---|---|---|
| App UI, nav, buttons, labels | **Inter** | 400/500/600 · 13–16px |
| Question text (exam & review) | **Lora** | 400/500 · 18–22px, line-height 1.6 |
| Result score | **Inter** (tabular) | 800 · 44–64px |
| Timer | **Inter** tabular-nums | 600 · 20–24px |

Google Fonts: `Inter:wght@400;500;600;700;800` + `Lora:ital,wght@0,400;0,500;1,400`.

Rules:
- Timer **always** uses `font-variant-numeric: tabular-nums` in a fixed-width container → no layout shift while ticking (CLS guard).
- Body text never below 14px; forms at 16px (prevents iOS zoom-on-focus).
- Question text min 18px — this app's core job is *reading English carefully*.

### 2.3 Shape, depth, spacing

- Radius: cards 12px, buttons/inputs 8px, chips 999px.
- Shadows: only `--shadow-sm` on cards (0 1px 2px rgba(0,0,0,0.05)); hover lifts to `--shadow-md` with `translateY(-1px)` over 150–200ms. No blur-heavy effects.
- Spacing: 8px grid; section padding 24–48px; page gutter 16px mobile → 32px desktop; content max-width 720px for reading surfaces, 1200px for dashboard/admin.

### 2.4 Icons & motion

- Icons: **lucide-react** only. Never emoji. Icon-only buttons carry `aria-label`.
- Motion: 150–200ms ease transitions for hover/focus/press; timer warning uses a gentle pulse **only** if `prefers-reduced-motion` is off. No entrance choreography, no parallax, no confetti.

---

## 3. The signature: POS highlight system

The review experience is anchored by highlighting the tested words **inside the question sentence**. This is the product's visual identity.

### 3.1 Rendering rules

- Each `highlight_regions` span renders as: bold text in its POS color + a distinct underline style + a small label chip on hover/focus.
- Color AND underline style together (never color alone — colorblind-safe):

| POS | Color | Underline |
|---|---|---|
| Verb | `--pos-verb` | dashed |
| Noun | `--pos-noun` | solid |
| Pronoun | `--pos-pronoun` | dotted |
| Adjective | `--pos-adjective` | double |
| Adverb | `--pos-adverb` | solid-thick |
| Preposition | `--pos-preposition` | wavy |
| Conjunction | `--pos-conjunction` | dash-dot |
| Determiner | `--pos-determiner` | thin-solid |
| Other | `--pos-other` | dotted-thin |

- A compact **legend** always visible at the top of the review ("Verb · Noun · Pronoun …"), each item clickable to filter/highlight that category across all questions (nice-to-have).
- Regions are non-interactive otherwise; they carry `aria-label="verb"` and a `title` tooltip.
- The **explanation** paragraph sits beside/below the question; the words it discusses are the ones highlighted above.

### 3.2 Example (mock)

```
The committee has decided to postpone the meeting.

   [Verb: dashed]     [Noun: solid]      [Verb: dashed]
   has decided        committee           decided

Explanation: "The verb 'has decided' agrees with the singular
collective noun 'committee'."
```

---

## 4. Pages

### 4.1 Login

```
┌──────────────────────────────────────────────┐
│                                              │
│              (logo)  TOEFL Prep              │
│                                              │
│          ┌────────────────────────┐          │
│          │  Email                 │          │
│          │  ┌──────────────────┐  │          │
│          │  │ student@…         │  │          │
│          │  └──────────────────┘  │          │
│          │  Password             │          │
│          │  ┌──────────────────┐  │          │
│          │  │ ••••••••          │  │          │
│          │  └──────────────────┘  │          │
│          │  [ Sign in ]          │          │
│          │                        │          │
│          │  Seeded accounts only  │          │
│          └────────────────────────┘          │
│                                              │
└──────────────────────────────────────────────┘
```

- Centered single card on paper; brand wordmark; single CTA.
- Errors inline under the failing field (not only at top).
- Hint line "Seeded accounts only" (honest for a 2-account tool).

### 4.2 Dashboard

```
┌────────────────────────────────────────────────────────┐
│ TOEFL Prep      [Dashboard] [Start exam] [Admin]  [me] │
├────────────────────────────────────────────────────────┤
│  Exams taken   Avg score   Best score    Trend         │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐   │
│  │    12     │ │   71%    │ │   85%    │ │  ↑ +4.2  │   │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘   │
│                                                        │
│  Score over time              Section accuracy         │
│  ┌────────────────────┐      ┌────────────────────┐    │
│  │      ╱╲    ╱        │      │ Structure ████░ 76%│    │
│  │    ╱╱  ╲╱╱  ╲      │      │ Vocabulary ███░░ 64%│    │
│  │  ╱          ╲     │      └────────────────────┘    │
│  └────────────────────┘      Weakest: Pronoun (41%)    │
│  Recent attempts                                        │
│  • Structure #3 · 82% · Aug 18  → review                │
│  • Vocabulary #2 · 64% · Aug 16 → review                │
└────────────────────────────────────────────────────────┘
```

- **Stat cards** (exams taken, average, best, trend) — 4 across on desktop, 2×2 on mobile.
- **Score over time** — Recharts **line/area chart** (per the chart guidance: trend-over-time → line/area). Points = attempts; trend line overlay shows direction. Tooltips show date + score + section split. Colorblind-safe: series differentiated by shape/marker + line style, not color alone.
- **Section accuracy** — two horizontal progress bars (waffle/bar), values as text beside them.
- **Weakest POS** card — computed from wrong answers' highlight regions; actionable signal.
- **Recent attempts** list, each row links to its review (deep-linkable, back works).
- **Empty state** (0 attempts): friendly one-liner + "Start your first exam" CTA instead of empty charts.

### 4.3 Quiz — per-question mode

```
┌──────────────────────────────────────────────┐
│ Quiz: Structure #3          Q 3 / 15   01:12 │
├──────────────────────────────────────────────┤
│  3. Choose the option that best completes    │
│     the sentence.                            │
│                                             │
│  The manager asked the staff _______ the     │
│  new procedure before Friday.                │
│                                             │
│  A  to review                               │
│  B  reviewing                               │
│  C  review                                  │
│  D  reviewed                                │
│                                             │
│              [ 1 ] [ 2 ] [ 3 ] [ 4 ]        │
└──────────────────────────────────────────────┘
```

- **Nothing but the exam.** Header row: exam title · progress ("Q 3 / 15") · timer (fixed width, tabular).
- Options as full-width cards; keyboard shortcuts **1–4 / A–D**; hover lifts; selected = blue fill.
- **Selecting an option immediately advances** (auto-next). No confirm — it's a drill.
- Timer per question (default 60s). **Last 10s: amber + gentle pulse** (reduced-motion safe). **0s: the question is marked unanswered and we advance.**
- Progress dots (1–15) small, subtle, bottom; answered = filled, current = ring.

### 4.4 Quiz — overall mode

```
┌──────────────────────────────────────────────┐
│ Mock TOEFL A              12:34 left    Submit│
├──────────────────────────────────────────────┤
│  Question palette        Question 7 of 20    │
│  ┌─┬─┬─┬─┬─┐                                │
│  │1│2│3│4│5│    The board of directors has   │
│  ├─┼─┼─┼─┼─┤    announced that the merger    │
│  │6│7│8│9│10│   will ____ next quarter.      │
│  ├─┼─┼─┼─┼─┤                                │
│  │11│12│13│14│15                            │
│  └─┴─┴─┴─┴─┘     A  occurs                 │
│                   B  occur                  │
│   Legend:          C  occurring             │
│   ■ answered       D  occurred              │
│   ◌ unanswered                             │
│   ◆ flagged                                │
└──────────────────────────────────────────────┘
```

- Palette: numbered grid; **answered** = filled blue, **unanswered** = outlined, **flagged** = amber corner marker. Click a number to jump.
- Flag icon per question to mark for review; flagged questions asked about in the submit dialog.
- Overall countdown top-right; **last 60s amber, last 10s red**.
- **Submit** opens a confirm dialog: "X answered, Y unanswered, Z flagged — submit now?" then finalizes. Auto-submit at 0:00.

### 4.5 Review

```
┌──────────────────────────────────────────────┐
│ Review: Structure #3              [Export PDF]│
│  Score   Correct   Wrong   Unanswered  Time   │
│   78%      12        2         1       6:40   │
│  Structure ████████░ 80% · Vocabulary ████░░ 70%│
│  Legend: [Verb][Noun][Pronoun][Adj]…          │
├──────────────────────────────────────────────┤
│ ✓ Q3 · Structure · 38s                        │
│   The manager asked the staff **to review**    │
│   the new procedure before Friday.            │
│   Your answer: A (correct)                    │
│   Explanation: "to review" is an infinitive    │
│   phrase; after "asked" the complement is      │
│   "to + base verb".                           │
├──────────────────────────────────────────────┤
│ ✗ Q7 · Structure · 45s                        │
│   The committee **has decided** to postpone…   │
│   Your answer: B ✗  Correct: C ✓              │
│   Explanation: collective noun "committee"     │
│   takes a singular verb, "has decided".       │
└──────────────────────────────────────────────┘
```

- **Summary hero**: big tabular score, four counts, per-section bars, POS legend.
- One card per question: correct (✓ green hairline) vs wrong (✗ red hairline); **highlighted regions rendered per §3**; your answer vs correct answer; per-item time; explanation.
- **Export PDF** → switches to print layout (§4.6) and calls `window.print()`.
- Wrong answers sort first (learn from mistakes), then correct — toggleable.

### 4.6 Print layout (PDF export)

- Quality is the requirement: **"when I print, the result must still look good."** The print stylesheet is treated as a first-class design surface, not an afterthought.
- `@media print`: hide nav/buttons/interactive chrome; expand all cards (no accordion); add a header block (app name, exam title, date, score summary); keep POS colors via `print-color-adjust: exact`; page margins ~12mm; question text never split awkwardly (`break-inside: avoid` on cards); serif question text + tabular score preserved so the PDF reads like a real marked paper.
- The print layout is **tested visually before ship** (a dev-only "Print preview" toggle). If the OS print dialog ever degrades rendering on a target browser, the fallback is a real `.pdf` via `react-pdf` — same layout, rendered to a downloadable file. Decision: ship print-first; adopt react-pdf only if a browser shows unacceptable output.

### 4.7 Admin — question editor

```
┌──────────────────────────────────────────────┐
│ Admin · Questions          [+ New]  [Seed]   │
│ Filters: [Section▾] [Type▾] [Search…]        │
│ ┌──────────────────────────────────────────┐ │
│ │ Structure · The committee has decided…    │ │
│ │ Type: sentence-completion · Medium · ✓   │ │
│ └──────────────────────────────────────────┘ │
│ (list rows)                                  │
├──────────────────────────────────────────────┤
│ Edit question                                │
│ Section  [Structure▾]  Type  [sentence…▾]    │
│ Question text (select text to highlight)     │
│ ┌──────────────────────────────────────────┐ │
│ │ The [committee] has decided to postpone   │ │
│ │ the meeting.                             │ │
│ └──────────────────────────────────────────┘ │
│ Highlight regions: [Noun: committee]  [x]    │
│   [+ add]  (select a span → pick POS)        │
│ Options   A [         ]  B [         ]       │
│           C [         ]  D [         ]       │
│ Correct   ( A ) ( B ) ( C ) ( D )            │
│ Explanation  (Bahasa Indonesia)              │
│ ┌──────────────────────────────────────────┐ │
│ │ Kata kerja "has decided" cocok dengan    │ │
│ │ kata benda tunggal "committee".          │ │
│ └──────────────────────────────────────────┘ │
│ Difficulty [Easy▾]   [Save] [Deactivate]     │
├──────────────────────────────────────────────┤
│ Import from AI        [Import JSON…]         │
│ (paste prompt output from question-          │
│  generator.md → validated → drafts)          │
└──────────────────────────────────────────────┘
```

- **Highlight editor is the key interaction**: admin selects a text span in the question textarea → a small popover lists POS categories → region created and rendered as a colored underline under the field in real time. Regions deletable inline. Validation errors (out-of-bounds, overlap) shown inline.
- **Import from AI**: a dedicated "Import JSON" panel accepts the pasted output of `prompts/question-generator.md` (run **outside the app** in Claude/Gemini/ChatGPT). Each valid item becomes a **draft** filling the form (question, options, correct, Indonesian explanation, highlights); invalid items are listed with reasons and don't block the rest. Admin reviews/edits drafts, then **Save** publishes to the bank. There is no AI button inside the app — generation happens in the chat tool.
- Bulk seed button reloads from `backend/seed/*.json` (idempotent).

### 4.8 Admin — exam editor

```
┌──────────────────────────────────────────────┐
│ Admin · Exams                     [+ New]    │
│ Title            Mock TOEFL A  [Published ✓] │
│ Sections  Structure [10]  Vocabulary [5]    │
│ Mode       ( ) Per-question  ( ) Overall     │
│            (•) Both (student chooses)        │
│ Per-question time   60 s   | Total time 15 m │
│ Shuffle questions   [✓]                      │
│ Bank check: ✓ 15 questions available         │
│                               [Save] [Delete]│
└──────────────────────────────────────────────┘
```

- Live "bank check" — confirms enough active questions per section before publish.
- Publish toggle is one clear switch.

---

## 5. Interaction & state rules

- **Every clickable element**: `cursor: pointer`, hover transition 150–200ms, visible focus ring (2px offset `--primary`).
- **Buttons**: primary (filled `--primary`), secondary (outline `--primary`), ghost (muted), destructive (red outline/fill for delete/submit-timeout).
- **Loading**: skeletons for lists/charts; button-level spinners for mutations; the quiz start button shows "Preparing…".
- **Empty states**: explicit, with a single next action (dashboard → start exam; question bank → add first question).
- **Errors**: inline field errors + a top-level toast for API failures; never a blank screen.
- **The timer is the one element allowed to animate**: amber pulse last 10s (reduced-motion: color-only change, no pulse).
- **Navigation**: student sees Dashboard / Start exam / History; admin additionally sees Questions / Exams; active route clearly marked; deep links work (dashboard rows, review URLs).

---

## 6. Accessibility & performance

- WCAG 2.1 AA: contrast ≥ 4.5:1 (POS/status colors all ≥ 7:1 on paper/card); focus visible; full keyboard nav; options selectable with 1–4/A–D; `aria-label` on icon buttons; timer has `aria-live="off"` (avoid screen-reader spam) but the question counter and expiry state are announced.
- `prefers-reduced-motion` respected (pulse/transitions degrade to color-only/instant).
- Touch targets ≥ 44×44px; forms ≥ 16px input text.
- CLS: timer in fixed-width tabular container; images none; charts reserve height.
- Performance: lazy-route code splitting; quiz bundle loaded eagerly (it's the core); fonts `display=swap`; Recharts rendered only on dashboard.
- Responsive: 375 / 768 / 1024 / 1440. Mobile: quiz = single column, options full-width, progress dots; overall-mode palette collapses to a horizontal scroll strip with sticky timer; dashboard 2×2 stat grid, charts full-width stacked.
- No horizontal scroll anywhere; nav never fixed-height blocking content.

---

## 7. Pre-delivery checklist (run before shipping any page)

- [ ] No emoji as icons — lucide-react SVG only.
- [ ] `cursor-pointer` on all clickable elements.
- [ ] Hover + focus states, 150–200ms transitions, visible focus rings.
- [ ] Contrast ≥ 4.5:1 everywhere; POS/status verified on paper AND card surfaces.
- [ ] `prefers-reduced-motion` respected.
- [ ] Timer uses tabular-nums in reserved-width container (no CLS).
- [ ] Forms: 16px inputs, inline errors, visible labels.
- [ ] Keyboard: full page navigable; quiz options via 1–4/A–D.
- [ ] Empty/loading/error states on every data page.
- [ ] Responsive at 375 / 768 / 1024 / 1440; no horizontal scroll.
- [ ] Colors used from token layer only — zero raw hex in components.