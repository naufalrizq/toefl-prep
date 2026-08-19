import type { QuizItem } from '../../types';

export default function QuestionPalette({
  items,
  answers,
  current,
  onJump,
}: {
  items: QuizItem[];
  answers: Record<number, number | null>;
  current: number;
  onJump: (index: number) => void;
}) {
  return (
    <div className="rounded-md bg-card p-4 shadow-sm">
      <p className="mb-3 text-[13px] font-semibold uppercase tracking-wide text-ink-muted">Questions</p>
      <div className="grid grid-cols-10 gap-1 overflow-x-auto pb-1 lg:grid-cols-7" aria-label="Question palette">
        {items.map((it, i) => {
          const answered = answers[it.id] !== null && answers[it.id] !== undefined;
          const isCurrent = i === current;
          return (
            <button
              key={it.id}
              type="button"
              onClick={() => onJump(i)}
              aria-label={`Question ${i + 1}, ${answered ? 'answered' : 'unanswered'}${it.flagged ? ', flagged' : ''}`}
              className={`relative flex h-7 items-center justify-center rounded-sm text-[11px] font-semibold transition-colors duration-200 ${
                isCurrent
                  ? 'ring-2 ring-primary ring-offset-1 ring-offset-card'
                  : answered
                    ? 'bg-primary text-on-primary hover:bg-primary/90'
                    : 'bg-card-muted text-ink-muted hover:bg-border-strong'
              }`}
            >
              {i + 1}
              {it.flagged && (
                <span className="absolute right-0.5 top-0.5 h-1.5 w-1.5 rounded-full bg-warning" aria-hidden="true" />
              )}
            </button>
          );
        })}
      </div>
      <div className="mt-3 flex flex-col gap-1 text-[12px] text-ink-muted">
        <span className="flex items-center gap-1.5">
          <span className="inline-block h-2.5 w-2.5 rounded-sm bg-primary" /> answered
        </span>
        <span className="flex items-center gap-1.5">
          <span className="inline-block h-2.5 w-2.5 rounded-sm border border-border-strong bg-card-muted" /> unanswered
        </span>
        <span className="flex items-center gap-1.5">
          <span className="inline-block h-2.5 w-2.5 rounded-full bg-warning" /> flagged
        </span>
      </div>
    </div>
  );
}