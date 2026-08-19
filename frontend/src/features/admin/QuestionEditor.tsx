import { useMemo, useRef, useState } from 'react';
import { Plus, Trash2 } from 'lucide-react';
import Button from '../../components/Button';
import HighlightedText from '../../components/HighlightedText';
import Badge from '../../components/Badge';
import { POS_LABEL, SECTION_LABEL, TYPE_LABEL, ALL_SECTIONS, ALL_TYPES, DEFAULT_TYPE } from '../../lib/format';
import type { Difficulty, HighlightRegion, Pos, Question, QuestionType, Section } from '../../types';

export interface QuestionDraft {
  id?: number;
  section: Section;
  type: QuestionType;
  question_text: string;
  passage?: string;
  options: string[];
  correct_index: number;
  explanation: string;
  difficulty: Difficulty;
  active: boolean;
  highlight_regions: HighlightRegion[];
}

export function emptyDraft(): QuestionDraft {
  return {
    section: 'structure',
    type: 'sentence-completion',
    question_text: '',
    passage: '',
    options: ['', '', '', ''],
    correct_index: 0,
    explanation: '',
    difficulty: 'medium',
    active: true,
    highlight_regions: [],
  };
}

export function fromQuestion(q: Question): QuestionDraft {
  return {
    id: q.id,
    section: q.section,
    type: q.type,
    question_text: q.question_text,
    passage: q.passage ?? '',
    options: [...q.options],
    correct_index: q.correct_index,
    explanation: q.explanation,
    difficulty: q.difficulty,
    active: q.active,
    highlight_regions: q.highlight_regions.map((r) => ({ ...r })),
  };
}

const POSES: Pos[] = [
  'verb',
  'noun',
  'pronoun',
  'adjective',
  'adverb',
  'preposition',
  'conjunction',
  'determiner',
  'other',
];

export default function QuestionEditor({
  draft,
  onChange,
  onSubmit,
  submitting,
  submitLabel,
  onCancel,
}: {
  draft: QuestionDraft;
  onChange: (d: QuestionDraft) => void;
  onSubmit: () => void;
  submitting: boolean;
  submitLabel: string;
  onCancel?: () => void;
}) {
  const taRef = useRef<HTMLTextAreaElement>(null);
  const [sel, setSel] = useState<{ start: number; end: number } | null>(null);
  const [errors, setErrors] = useState<Record<string, string>>({});

  const validation = useMemo(() => {
    const v: Record<string, string> = {};
    if (!draft.question_text.trim()) v.question_text = 'Question text is required.';
    else {
      for (const r of draft.highlight_regions) {
        if (r.start < 0 || r.end > draft.question_text.length || r.start >= r.end) {
          v.highlight_regions = `Region out of bounds for "${draft.question_text.slice(r.start, r.end) || draft.question_text}"`;
        }
        if (draft.highlight_regions.filter((o) => !(o.start >= r.end || o.end <= r.start)).length > 1) {
          v.highlight_regions = 'Overlapping highlight regions are not allowed.';
        }
      }
    }
    if (draft.type === 'reading-comprehension' && !draft.passage?.trim()) v.passage = 'Reading questions require a passage.';
    draft.options.forEach((o, i) => {
      if (!o.trim()) v[`option_${i}`] = 'Option is required.';
    });
    if (!draft.explanation.trim()) v.explanation = 'Explanation (Bahasa Indonesia) is required.';
    return v;
  }, [draft]);

  function captureSelection() {
    const ta = taRef.current;
    if (!ta) return;
    const s = ta.selectionStart;
    const e = ta.selectionEnd;
    if (e > s) setSel({ start: s, end: e });
  }

  function addRegion(pos: Pos) {
    if (!sel) return;
    const region: HighlightRegion = { start: sel.start, end: sel.end, pos };
    onChange({ ...draft, highlight_regions: [...draft.highlight_regions, region] });
    setSel(null);
  }

  function removeRegion(i: number) {
    onChange({
      ...draft,
      highlight_regions: draft.highlight_regions.filter((_, idx) => idx !== i),
    });
  }

  function handleSubmit() {
    const v = { ...validation, ...errors };
    if (Object.keys(v).length) {
      setErrors(v);
      return;
    }
    setErrors({});
    onSubmit();
  }

  return (
    <form
      className="flex flex-col gap-4"
      onSubmit={(e) => {
        e.preventDefault();
        handleSubmit();
      }}
      noValidate
    >
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <label className="flex flex-col gap-1.5 text-sm font-medium text-ink">
          Section
          <select
            value={draft.section}
            onChange={(e) => {
              const section = e.target.value as Section;
              onChange({ ...draft, section, type: DEFAULT_TYPE[section] });
            }}
            className="rounded-sm border border-border-strong bg-card px-3 py-3 text-base"
          >
            {ALL_SECTIONS.map((s) => (
              <option key={s} value={s}>
                {SECTION_LABEL[s]}
              </option>
            ))}
          </select>
        </label>
        <label className="flex flex-col gap-1.5 text-sm font-medium text-ink">
          Type
          <select
            value={draft.type}
            onChange={(e) => onChange({ ...draft, type: e.target.value as QuestionType })}
            className="rounded-sm border border-border-strong bg-card px-3 py-3 text-base"
          >
            {ALL_TYPES.map((t) => (
              <option key={t} value={t}>
                {TYPE_LABEL[t]}
              </option>
            ))}
          </select>
        </label>
      </div>

      {draft.type === 'reading-comprehension' && (
        <div className="flex flex-col gap-1.5">
          <label className="text-sm font-medium text-ink" htmlFor="qp">
            Passage
          </label>
          <textarea
            id="qp"
            rows={6}
            value={draft.passage}
            onChange={(e) => onChange({ ...draft, passage: e.target.value })}
            aria-invalid={errors.passage ? true : undefined}
            className={`rounded-sm border px-4 py-3 text-base leading-relaxed text-ink transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-primary/30 ${
              errors.passage ? 'border-danger' : 'border-border-strong'
            }`}
            placeholder="The passage the question refers to…"
          />
          {errors.passage && (
            <p className="text-[13px] text-danger" role="alert">{errors.passage}</p>
          )}
        </div>
      )}

      <div className="flex flex-col gap-1.5">
        <label className="text-sm font-medium text-ink" htmlFor="qt">
          Question text
        </label>
        <p className="text-[13px] text-ink-muted">Select a span of text, then pick a part of speech to highlight it.</p>
        <textarea
          id="qt"
          ref={taRef}
          rows={4}
          value={draft.question_text}
          onChange={(e) => onChange({ ...draft, question_text: e.target.value })}
          onMouseUp={captureSelection}
          onKeyUp={captureSelection}
          className={`rounded-sm border px-4 py-3 text-base leading-relaxed text-ink transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-primary/30 ${
            errors.question_text ? 'border-danger' : 'border-border-strong'
          }`}
          aria-invalid={errors.question_text ? true : undefined}
        />
        {sel && (
          <div className="flex flex-wrap items-center gap-1.5 rounded-sm bg-primary-soft p-2">
            <span className="text-[13px] font-medium text-primary">
              “{draft.question_text.slice(sel.start, sel.end)}”
            </span>
            {POSES.map((p) => (
              <button
                key={p}
                type="button"
                onClick={() => addRegion(p)}
                className="rounded-full px-2.5 py-1 text-[12px] font-semibold text-primary transition-colors duration-200 hover:bg-primary hover:text-on-primary"
              >
                {POS_LABEL[p]}
              </button>
            ))}
          </div>
        )}
        {errors.question_text && (
          <p className="text-[13px] text-danger" role="alert">{errors.question_text}</p>
        )}
        {draft.highlight_regions.length > 0 && (
          <div className="rounded-sm bg-card-muted p-3">
            <HighlightedText text={draft.question_text} regions={draft.highlight_regions} showLabels />
            <ul className="mt-2 flex flex-wrap gap-1.5">
              {draft.highlight_regions.map((r, i) => (
                <li key={i} className="flex items-center gap-1">
                  <Badge tone="primary">{POS_LABEL[r.pos]}: {draft.question_text.slice(r.start, r.end) || '…'}</Badge>
                  <button
                    type="button"
                    onClick={() => removeRegion(i)}
                    aria-label={`Remove highlight ${POS_LABEL[r.pos]}`}
                    className="rounded-sm p-1 text-ink-muted transition-colors duration-200 hover:bg-danger-soft hover:text-danger"
                  >
                    <Trash2 className="h-3.5 w-3.5" aria-hidden="true" />
                  </button>
                </li>
              ))}
            </ul>
          </div>
        )}
        {errors.highlight_regions && (
          <p className="text-[13px] text-danger" role="alert">{errors.highlight_regions}</p>
        )}
      </div>

      <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
        {draft.options.map((opt, i) => (
          <div key={i} className="flex items-center gap-2">
            <span className="flex h-8 w-8 shrink-0 items-center justify-center rounded-sm bg-card-muted text-sm font-semibold text-ink-muted">
              {String.fromCharCode(65 + i)}
            </span>
            <input
              type="text"
              value={opt}
              onChange={(e) =>
                onChange({ ...draft, options: draft.options.map((o, oi) => (oi === i ? e.target.value : o)) })
              }
              aria-label={`Option ${String.fromCharCode(65 + i)}`}
              aria-invalid={errors[`option_${i}`] ? true : undefined}
              className={`w-full rounded-sm border px-3 py-2.5 text-base transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-primary/30 ${
                errors[`option_${i}`] ? 'border-danger' : 'border-border-strong'
              }`}
            />
          </div>
        ))}
      </div>
      {errors.option_0 && <p className="text-[13px] text-danger" role="alert">All four options are required.</p>}

      <div className="flex items-center gap-2">
        <span className="text-sm font-medium text-ink">Correct answer</span>
        {draft.options.map((_, i) => (
          <button
            key={i}
            type="button"
            onClick={() => onChange({ ...draft, correct_index: i })}
            aria-pressed={draft.correct_index === i}
            className={`flex h-8 w-8 items-center justify-center rounded-full text-sm font-semibold transition-all duration-200 ${
              draft.correct_index === i ? 'bg-success text-white' : 'bg-card-muted text-ink-muted hover:bg-border-strong'
            }`}
          >
            {String.fromCharCode(65 + i)}
          </button>
        ))}
      </div>

      <label className="flex flex-col gap-1.5 text-sm font-medium text-ink">
        Explanation <span className="text-[12px] font-normal text-ink-muted">in Bahasa Indonesia</span>
        <textarea
          rows={3}
          value={draft.explanation}
          onChange={(e) => onChange({ ...draft, explanation: e.target.value })}
          aria-invalid={errors.explanation ? true : undefined}
          className={`rounded-sm border px-4 py-3 text-base leading-relaxed transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-primary/30 ${
            errors.explanation ? 'border-danger' : 'border-border-strong'
          }`}
        />
        {errors.explanation && <span className="text-[13px] text-danger" role="alert">{errors.explanation}</span>}
      </label>

      <div className="flex flex-wrap items-center gap-2">
        <span className="text-sm font-medium text-ink">Difficulty</span>
        {(['easy', 'medium', 'hard'] as Difficulty[]).map((d) => (
          <button
            key={d}
            type="button"
            onClick={() => onChange({ ...draft, difficulty: d })}
            aria-pressed={draft.difficulty === d}
            className={`rounded-full px-3 py-1 text-[13px] font-semibold capitalize transition-colors duration-200 ${
              draft.difficulty === d ? 'bg-primary text-on-primary' : 'bg-card-muted text-ink-muted hover:bg-border-strong'
            }`}
          >
            {d}
          </button>
        ))}
      </div>

      <div className="flex items-center justify-end gap-2 border-t border-border pt-4">
        {onCancel && (
          <Button type="button" variant="ghost" onClick={onCancel}>
            Cancel
          </Button>
        )}
        <Button type="submit" loading={submitting}>
          <Plus className="h-4 w-4" aria-hidden="true" />
          {submitLabel}
        </Button>
      </div>
    </form>
  );
}