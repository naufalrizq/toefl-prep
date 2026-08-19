# AI Question Generator — Prompt & Output Contract

**Purpose:** Generate TOEFL-style questions (Structure sentence-completion, Vocabulary, Reading comprehension, and Advanced Grammar error-identification) with answer keys, Indonesian explanations, and part-of-speech highlight annotations, in a JSON shape the TOEFL Prep system accepts.

**How it is used (generation happens OUTSIDE the app):**
1. **You** (the owner) copy the *system prompt* in §1, pick an instruction template in §2, and run it in **Claude, Gemini, or ChatGPT** to draft questions.
2. You copy the JSON output and paste it into the admin editor's **Import JSON** box in the app.
3. The app validates the JSON (rules in §3), turns valid items into **drafts**, and you review + save them to the bank. The app never calls an LLM and never holds an API key.

> Language rules: question text and options are **English**. Explanations are **Bahasa Indonesia**. This is a hard contract — see §3.

---

## 0. Style reference — read this before writing anything

Before drafting, read `examplequestion.md` (a curated set of approved items) and
match its style exactly.

### Question style (English)

- Sentence-completion items: ONE sentence with ONE blank `_____`; the four
  options (A–D) are short phrases that fill the blank. Test one classic
  Structure & Written Expression point per item:
  - **Subject–verb agreement**, including with interrupter phrases
    ("unlike many others", "as well as", "along with", "together with").
  - **Inversion** after negative openers ("Never before", "Under no
    circumstances", "Not until", "No sooner…than", "Hardly…when",
    "Scarcely…when", "So + adjective … that").
  - **Subjunctive** after requirement verbs/expressions ("suggest that",
    "insist that", "essential that").
  - **Fixed patterns** ("It was not until X that Y", "It was X who…",
    "the + comparative, the + comparative").
  - **Agreement pairs** ("A number of" vs "the number of", "Each of",
    "Neither…nor", "Both…and", "Whether…or not").
  - **Reduced relative clauses** (V3 as a post-modifier).
- Vocabulary items: a target word inside the sentence plus a clear task
  ("The word 'X' is closest in meaning to:").
- Reading items (`section: "reading"`): a short academic passage (`passage`,
  150–350 words) plus ONE question per item. The question can test main idea,
  detail, inference, or a vocabulary-in-context phrase ("The phrase 'X' most
  nearly means…"). Write several questions against the same passage so a part
  can reuse it, but each JSON item carries the full passage again (no
  references).
- Error-identification items (`section: "grammar_adv"`): `question_text` is ONE
  complete sentence containing the error. The four options (A–D) are the four
  SEGMENTS of that sentence; exactly one segment is grammatically wrong, and
  `correct_index` points to it. Test one advanced point per item (subject–verb
  agreement with "neither…nor" / interrupter phrases, parallel structure,
  pronoun case, sequence of tenses, dangling modifiers, comparative forms).
- Keep sentences academic but natural. No invented idioms. Exactly one
  defensible answer; distractors tempt a *plausible wrong rule* (e.g. a
  distractor that ignores an interrupter phrase).

### Explanation style (Bahasa Indonesia)

- Casual and plain ("bahasa simpel"), like a tutor explaining out loud. Short
  sentences, everyday words ("di tengah", "sisipan", "paten"), no stiff
  textbook phrasing.
- Structure each explanation as: **answer letter + the word**, then WHY — name
  the tested point and point out the trap that makes the distractors tempting
  (e.g. *"Ada sisipan 'as well as…' di tengah, tapi subjek utamanya tetap 'the
  findings' yang jamak, jadi 'were'."*).
- 1–3 sentences. If a pattern has a fixed pair ("no sooner…than",
  "hardly…when", "the more…the more"), say it plainly.
- Reference the words you mark in `highlights` when natural.

---

## 1. System prompt (send this as the system message)

```
You are an expert TOEFL (Structure and Written Expression + Vocabulary) item
writer. You produce multiple-choice questions that are grammatically precise,
unambiguous, and fair. You never invent idioms that change meaning, and every
question has exactly one defensible correct answer.

You must respond with VALID JSON ONLY. No markdown, no code fences, no
explanatory text outside the JSON. The JSON must conform exactly to the output
schema in the user's instruction.

Hard rules for every question:

1. Four options exactly (A, B, C, D). Exactly one is correct.
2. Distractors are plausible but clearly wrong for a defensible reason you can
   state in the explanation.
3. "section" is one of: "structure" | "vocabulary" | "reading" | "grammar_adv".
4. "type" is one of: "sentence-completion" | "vocab-multiple-choice" |
   "reading-comprehension" | "error-identification".
   - sentence-completion (section: structure): "question_text" is an English
     sentence with ONE blank written as "_____" (4 underscores). The options
     fill the blank so the sentence is grammatically correct.
   - vocab-multiple-choice (section: vocabulary): "question_text" is an English
     sentence containing a target word, followed by a clear task, e.g.:
     "The word 'mitigate' is closest in meaning to:"
   - reading-comprehension (section: reading): "passage" is the academic text
     (150–350 words, ≤ 4000 characters) and "question_text" is ONE question
     about it. Options are the four answer choices.
   - error-identification (section: grammar_adv): "question_text" is one full
     sentence; the four options are the four segments of that sentence (A–D),
     and "correct_index" is the segment that is grammatically wrong.
5. "difficulty" is one of: "easy" | "medium" | "hard". Calibrate to TOEFL iBT
   intermediate (easy), upper-intermediate (medium), advanced (hard).
6. "explanation" is written in plain, casual BAHASA INDONESIA following the
   style in examplequestion.md (§0): start with the answer letter + the word,
   then say WHY in simple terms, name the grammar point, and point out the
   trap that makes the distractors tempting. 1–3 short sentences. Reference
   the words you mark in "highlights" when natural.
7. "highlights" is an object mapping a part-of-speech key to an ARRAY OF
   PHRASES (exact substrings, verbatim from "question_text"). Allowed keys:
   verb, noun, pronoun, adjective, adverb, preposition, conjunction,
   determiner, other. Mark ONLY words that the explanation actually discusses.
   Keep it to 1–3 keys, 1–2 phrases each. The phrases must appear verbatim
   (case-insensitive) inside "question_text" — the system will locate them
   automatically; phrases it cannot find are discarded, so do not paraphrase.
8. Constraints: "question_text" ≤ 1000 characters; "passage" (reading only) ≤
   4000 characters and REQUIRED for reading-comprehension; each option ≤ 200
   characters; "correct_index" is 0-based (A=0, B=1, C=2, D=3).
9. Never include audio, images, tables, or anything the schema does not have.
```

---

## 2. User instruction template

Use one of these instructions (with the system prompt above) per call.

### 2a. Single question

```
Generate one TOEFL Structure question.

Section: structure
Type: sentence-completion
Difficulty: medium
Topic hint (optional): subject-verb agreement with collective nouns

Return exactly:
{"section":"structure","type":"sentence-completion","difficulty":"medium",
 "question_text":"...with _____ ...",
 "options":["...","...","...","..."],
 "correct_index":0,
 "explanation":"...",
 "highlights":{"noun":["..."],"verb":["..."]}}
```

### 2b. Batch (e.g. 5 questions)

```
Generate 5 TOEFL questions as a JSON array.

Section: vocabulary
Type: vocab-multiple-choice
Difficulty: mixed (1 easy, 3 medium, 1 hard)
Topic hints: academic word list, words from economics and science

Return a JSON array where each element is exactly:
{"section":"vocabulary","type":"vocab-multiple-choice","difficulty":"...",
 "question_text":"...",
 "options":["...","...","...","..."],
 "correct_index":0,
 "explanation":"...",
 "highlights":{"verb":["..."]}}

No other keys. No markdown. No code fences.
```

### 2c. Reading batch (one passage, several questions)

```
Generate a TOEFL Reading set: 1 passage and 3 questions about it, as a JSON
array of 3 items that all repeat the SAME "passage".

Section: reading
Type: reading-comprehension
Difficulty: mixed (easy, medium, hard)
Passage topic hint: short academic text (150-350 words), e.g. a science or
history topic.

Each element is exactly:
{"section":"reading","type":"reading-comprehension","difficulty":"...",
 "passage":"<the full passage, verbatim and identical in all 3 items>",
 "question_text":"<one question about the passage>",
 "options":["...","...","...","..."],
 "correct_index":0,
 "explanation":"<casual Indonesian, cite the sentence in the passage>",
 "highlights":{}}

No other keys. No markdown. No code fences.
```

### 2d. Error-identification batch

```
Generate 4 TOEFL Error-Identification questions as a JSON array.

Section: grammar_adv
Type: error-identification
Difficulty: mixed (1 easy, 2 medium, 1 hard)
Topic hints: subject-verb agreement with neither/nor, parallel structure,
pronoun case, sequence of tenses.

Each element is exactly:
{"section":"grammar_adv","type":"error-identification","difficulty":"...",
 "question_text":"<one full English sentence containing the error>",
 "options":["<segment A>","<segment B>","<segment C>","<segment D>"],
 "correct_index":<index of the WRONG segment, 0-3>,
 "explanation":"<casual Indonesian: which segment, why wrong, the correct form>",
 "highlights":{"verb":["..."]}}

The four options must cover the whole sentence in order, with no gaps.
Exactly one segment is wrong. No other keys. No markdown. No code fences.
```

---

## 3. Output contract (backend validation)

The backend MUST validate every AI response before storing. Rules mirror
SRS FR-2.5 / FR-2.6:

| Field | Rule |
|---|---|
| `section` | ∈ {structure, vocabulary, reading, grammar_adv} |
| `type` | ∈ {sentence-completion, vocab-multiple-choice, reading-comprehension, error-identification} |
| `difficulty` | ∈ {easy, medium, hard} |
| `question_text` | non-empty, ≤ 1000 chars; structure type MUST contain `_____` |
| `passage` | REQUIRED and ≤ 4000 chars when type = reading-comprehension |
| `options` | exactly 4, non-empty, ≤ 200 chars, all distinct |
| `correct_index` | integer 0–3 |
| `explanation` | non-empty, ≤ 2000 chars, **Indonesian** |
| `highlights` | keys ⊆ allowed POS keys; values are arrays of strings |

### Highlight normalization (backend, robust against AI offset drift)

The AI returns **phrases**, not offsets. The backend converts phrases to
`highlight_regions` (`{start, end, pos, label}`):

1. For each `pos` → phrase, search `question_text` case-insensitively at
   word boundaries.
2. First match wins; compute `start`/`end` in **rune offsets** (Go `[]rune`)
   and UTF-8-safe indices (JS `string` indices) — never byte offsets.
3. If a phrase is not found, skip it and log a warning; do NOT fail the whole
   item.
4. If the AI returns `highlight_regions` directly (offsets), validate them per
   FR-2.6 and reject on overlap/out-of-bounds.
5. Normalized regions are stored in `questions.highlight_regions`; the
   phrases are discarded.

---

## 4. Example outputs (for reference / prompt validation)

### Structure · sentence-completion

```json
{
  "section": "structure",
  "type": "sentence-completion",
  "difficulty": "medium",
  "question_text": "The committee _____ decided to postpone the meeting until next week.",
  "options": ["have", "has", "are", "were"],
  "correct_index": 1,
  "explanation": "B — has. Subjek utamanya 'the committee' (tunggal, dianggap satu kesatuan), jadi verb-nya 'has'. 'have', 'are', dan 'were' semuanya bentuk jamak dan nggak cocok di sini.",
  "highlights": {
    "noun": ["committee"],
    "verb": ["decided"]
  }
}
```

### Vocabulary · vocab-multiple-choice

```json
{
  "section": "vocabulary",
  "type": "vocab-multiple-choice",
  "difficulty": "hard",
  "question_text": "The word 'mitigate' is closest in meaning to:",
  "options": ["to worsen", "to lessen", "to ignore", "to emphasize"],
  "correct_index": 1,
  "explanation": "B — to lessen. 'Mitigate' artinya ngurangin dampak buruk, jadi padanan paling deket 'to lessen'. 'To worsen' kebalikannya, sedangkan 'ignore' dan 'emphasize' nggak ada hubungannya sama ngurangin.",
  "highlights": {
    "verb": ["mitigate"]
  }
}
```

### Reading · reading-comprehension

```json
{
  "section": "reading",
  "type": "reading-comprehension",
  "difficulty": "medium",
  "passage": "Many migratory birds navigate using the Earth's magnetic field. Researchers discovered that birds rely on small iron-rich particles in their beaks, which may act like a built-in compass. However, this sense works best on clear days. On cloudy days, the birds appear to switch to landmarks such as rivers and coastlines, suggesting that the magnetic sense is not their only navigational tool.",
  "question_text": "Why do birds rely on landmarks on cloudy days?",
  "options": [
    "Landmarks provide more accurate routes than rivers.",
    "Their magnetic compass works poorly without clear skies.",
    "Iron particles in their beaks disappear in the rain.",
    "Clear days make it harder to see the coastline."
  ],
  "correct_index": 1,
  "explanation": "B. Kalimat 'this sense works best on clear days' lalu 'on cloudy days ... switch to landmarks' — indra magnetik menurun saat langit mendung, jadi burung beralih ke penanda darat. A dan C tidak disebut, D kebalikan dari kenyataan.",
  "highlights": {}
}
```

### Advanced Grammar · error-identification

```json
{
  "section": "grammar_adv",
  "type": "error-identification",
  "difficulty": "medium",
  "question_text": "Neither the manager nor the employees was satisfied with the new policy.",
  "options": [
    "Neither the manager",
    "nor the employees was",
    "satisfied with",
    "the new policy"
  ],
  "correct_index": 1,
  "explanation": "B. Aturan 'neither ... nor': verb mengikuti subjek terdekat, yaitu 'the employees' (jamak) → harus 'were satisfied', bukan 'was'. Subjek jamak + was = inkonsisten agreement.",
  "highlights": {
    "verb": ["was satisfied"],
    "conjunction": ["Neither", "nor"]
  }
}
```

---

## 5. Usage notes (drafting outside the app)

- Paste the **system prompt** (§1) + one **instruction** (§2) into Claude/Gemini/ChatGPT. Ask the chat for "only the JSON" if the tool adds prose.
- For a full mock, generate one batch per section (2c/2d/2b for reading / grammar_adv / structure+vocabulary), then assemble the exam in Admin · Exams with the per-section **parts** editor.
- Before generating, re-read `examplequestion.md` and match the question + explanation style in §0.
- If the tool wraps the JSON in a markdown code fence, copy just the JSON block (or the app strips a single fence on import).
- Batch: ask for a JSON array; each item is validated independently on import, so one bad item won't block the rest.
- The app's import rules in §3 mirror the app's server-side validation (SRS FR-2.5 / FR-2.6). If you change a rule here, update the SRS too — keep them in sync.