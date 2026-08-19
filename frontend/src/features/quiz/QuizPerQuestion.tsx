import { useEffect, useState } from 'react';
import { ArrowLeft, CheckCircle2, Flag, Send } from 'lucide-react';
import HighlightedText from '../../components/HighlightedText';
import Timer from '../../components/Timer';
import Button from '../../components/Button';
import QuestionPalette from './QuestionPalette';
import type { QuizItem } from '../../types';
import { SECTION_LABEL, TYPE_LABEL } from '../../lib/format';

const OPTION_KEYS = ['A', 'B', 'C', 'D'];

export default function QuizPerQuestion({
  items,
  answers,
  secondsPerQuestion,
  onAnswer,
  onFlag,
  onSubmit,
  submitting,
}: {
  items: QuizItem[];
  answers: Record<number, number | null>;
  secondsPerQuestion: number;
  onAnswer: (itemId: number, chosenIndex: number | null, timeTakenMs: number) => void;
  onFlag: (item: QuizItem, flagged: boolean) => void;
  onSubmit: () => void;
  submitting: boolean;
}) {
  const [index, setIndex] = useState(0);
  const [questionStartedAt, setQuestionStartedAt] = useState<number>(Date.now());
  const [deadline, setDeadline] = useState<number>(Date.now() + secondsPerQuestion * 1000);

  const item = items[index];
  const chosen = answers[item.id] ?? null;
  const isLast = index === items.length - 1;
  const isFirst = index === 0;

  useEffect(() => {
    setQuestionStartedAt(Date.now());
    setDeadline(Date.now() + secondsPerQuestion * 1000);
  }, [index, secondsPerQuestion]);

  function answer(choice: number | null) {
    onAnswer(item.id, choice, Date.now() - questionStartedAt);
    if (!isLast) setIndex((i) => i + 1);
  }

  function handleExpire() {
    const alreadyAnswered = answers[item.id] !== null && answers[item.id] !== undefined;
    if (alreadyAnswered) {
      if (!isLast) setIndex((i) => i + 1);
      return;
    }
    answer(null);
  }

  function go(delta: number) {
    setIndex((i) => Math.max(0, Math.min(items.length - 1, i + delta)));
  }

  useEffect(() => {
    function onKey(e: KeyboardEvent) {
      if (e.repeat) return;
      const k = e.key.toLowerCase();
      if (['1', '2', '3', '4'].includes(k)) answer(Number(k) - 1);
      else if (['a', 'b', 'c', 'd'].includes(k)) answer(k.charCodeAt(0) - 97);
    }
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [index, item?.id, questionStartedAt, chosen, isLast]);

  return (
    <div className="flex flex-col gap-6 lg:flex-row">
      <div className="min-w-0 flex-1">
        <div className="mx-auto flex max-w-reading flex-col gap-6">
          <div className="flex items-center justify-between border-b border-border pb-4">
            <div className="flex items-center gap-3">
              <span className="text-sm text-ink-muted">
                Q {index + 1} / {items.length}
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
            <Timer deadlineMs={deadline} warnAfterSec={10} dangerAfterSec={3} onExpire={handleExpire} label="Question" />
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

          <div className="flex flex-col gap-3" role="listbox" aria-label="Options">
            {item.question_snapshot.options.map((opt, i) => {
              const selected = chosen === i;
              return (
                <button
                  key={i}
                  type="button"
                  role="option"
                  aria-selected={selected}
                  onClick={() => answer(i)}
                  className={`flex min-h-[48px] items-center gap-3 rounded-sm border px-4 py-3 text-left transition-all duration-200 ${
                    selected
                      ? 'border-primary bg-primary-soft'
                      : 'border-border-strong bg-card hover:-translate-y-px hover:border-primary/50 hover:shadow-sm'
                  }`}
                >
                  <span className={`flex h-7 w-7 shrink-0 items-center justify-center rounded-sm text-sm font-semibold ${selected ? 'bg-primary text-on-primary' : 'bg-card-muted text-ink-muted'}`}>
                    {OPTION_KEYS[i]}
                  </span>
                  <span className="question-text text-[17px] leading-snug text-ink">{opt}</span>
                </button>
              );
            })}
          </div>

          <div className="flex items-center justify-between border-t border-border pt-4">
            <Button variant="secondary" size="sm" onClick={() => go(-1)} disabled={isFirst}>
              <ArrowLeft className="h-4 w-4" aria-hidden="true" />
              Previous
            </Button>
            {isLast && (
              <Button size="sm" onClick={onSubmit} disabled={submitting} loading={submitting}>
                <Send className="h-4 w-4" aria-hidden="true" />
                Finish & submit
              </Button>
            )}
          </div>
        </div>
      </div>

      <aside className="print-hidden shrink-0 lg:w-64">
        <div className="sticky top-20">
          <QuestionPalette
            items={items}
            answers={answers}
            current={index}
            onJump={(i) => setIndex(i)}
          />
          {item.flagged && (
            <p className="mt-3 flex items-center gap-1.5 text-[13px] text-warning">
              <CheckCircle2 className="h-4 w-4" aria-hidden="true" />
              Flagged. Review this one before submitting.
            </p>
          )}
        </div>
      </aside>
    </div>
  );
}