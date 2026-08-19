import { useState } from 'react';
import { Flag, Send } from 'lucide-react';
import HighlightedText from '../../components/HighlightedText';
import Timer from '../../components/Timer';
import Button from '../../components/Button';
import QuestionPalette from './QuestionPalette';
import type { QuizItem } from '../../types';
import { SECTION_LABEL, TYPE_LABEL } from '../../lib/format';

const OPTION_KEYS = ['A', 'B', 'C', 'D'];

export default function QuizOverall({
  items,
  answers,
  deadlineMs,
  onAnswer,
  onFlag,
  onSubmit,
  submitting,
}: {
  items: QuizItem[];
  answers: Record<number, number | null>;
  deadlineMs: number;
  onAnswer: (itemId: number, chosenIndex: number | null) => void;
  onFlag: (item: QuizItem, flagged: boolean) => void;
  onSubmit: () => void;
  submitting: boolean;
}) {
  const [current, setCurrent] = useState(0);
  const [confirmOpen, setConfirmOpen] = useState(false);
  const item = items[current];

  const answeredCount = Object.values(answers).filter((v) => v !== null && v !== undefined).length;
  const flaggedCount = items.filter((i) => i.flagged).length;
  const unansweredCount = items.length - answeredCount;

  function handleSubmitDialog() {
    setConfirmOpen(false);
    onSubmit();
  }

  return (
    <div className="flex flex-col gap-6 lg:flex-row">
      <div className="flex-1">
        <div className="mx-auto flex max-w-reading flex-col gap-6">
          <div className="flex items-center justify-between border-b border-border pb-4">
            <div className="flex items-center gap-3">
              <span className="text-sm text-ink-muted">
                Question {current + 1} of {items.length}
              </span>
              <button
                type="button"
                onClick={() => onFlag(item, !item.flagged)}
                className={`inline-flex items-center gap-1.5 rounded-sm px-2.5 py-1.5 text-sm transition-colors duration-200 ${
                  item.flagged ? 'bg-card-muted text-warning' : 'text-ink-muted hover:text-ink'
                }`}
                aria-pressed={item.flagged}
              >
                <Flag className="h-4 w-4" aria-hidden="true" />
                Flag
              </button>
            </div>
            <Button variant="secondary" size="sm" onClick={() => setConfirmOpen(true)} disabled={submitting}>
              Submit
            </Button>
          </div>

          <div>
            <p className="mb-3 text-[13px] font-medium uppercase tracking-wide text-ink-faint">
              {SECTION_LABEL[item.question_snapshot.section]}
              {item.question_snapshot.part ? ` · ${item.question_snapshot.part}` : ''} · {TYPE_LABEL[item.question_snapshot.type]}
            </p>
            {item.question_snapshot.passage ? (
              <div className="mb-4 max-h-64 overflow-y-auto rounded-sm border border-border bg-card-muted p-4 text-[15px] leading-relaxed text-ink">
                {item.question_snapshot.passage}
              </div>
            ) : null}
            <HighlightedText text={item.question_snapshot.question_text} regions={item.question_snapshot.highlight_regions} />
          </div>

          <div className="flex flex-col gap-3">
            {item.question_snapshot.options.map((opt, i) => {
              const selected = answers[item.id] === i;
              return (
                <button
                  key={i}
                  type="button"
                  onClick={() => onAnswer(item.id, i)}
                  className={`flex min-h-[48px] items-center gap-3 rounded-sm border px-4 py-3 text-left transition-all duration-200 ${
                    selected
                      ? 'border-primary bg-primary-soft'
                      : 'border-border-strong bg-card hover:-translate-y-px hover:border-primary/50 hover:shadow-sm'
                  }`}
                  aria-pressed={selected}
                >
                  <span className={`flex h-7 w-7 shrink-0 items-center justify-center rounded-sm text-sm font-semibold ${selected ? 'bg-primary text-on-primary' : 'bg-card-muted text-ink-muted'}`}>
                    {OPTION_KEYS[i]}
                  </span>
                  <span className="question-text text-[17px] leading-snug text-ink">{opt}</span>
                </button>
              );
            })}
          </div>
        </div>
      </div>

      <aside className="print-hidden shrink-0 lg:w-64">
        <div className="sticky top-20 flex flex-col gap-4">
          <div className="flex items-center justify-end">
            <Timer deadlineMs={deadlineMs} warnAfterSec={60} dangerAfterSec={10} onExpire={onSubmit} label="Time left" />
          </div>

          <QuestionPalette
            items={items}
            answers={answers}
            current={current}
            onJump={setCurrent}
          />
        </div>
      </aside>

      {confirmOpen && (
        <div className="fixed inset-0 z-30 flex items-center justify-center bg-black/40 p-4 backdrop-blur-sm" role="dialog" aria-modal="true">
          <div className="w-full max-w-sm rounded-md bg-card p-6 shadow-lg">
            <h2 className="text-lg font-semibold text-ink">Submit now?</h2>
            <p className="mt-2 text-sm text-ink-muted">
              {answeredCount} answered, {unansweredCount} unanswered, {flaggedCount} flagged.
            </p>
            <div className="mt-5 flex justify-end gap-2">
              <Button variant="ghost" onClick={() => setConfirmOpen(false)}>
                Keep working
              </Button>
              <Button onClick={handleSubmitDialog} disabled={submitting} loading={submitting}>
                <Send className="h-4 w-4" aria-hidden="true" />
                Submit
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}