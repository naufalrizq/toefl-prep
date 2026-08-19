import { useMemo, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { ArrowLeft, Check, Printer, X } from 'lucide-react';
import { api } from '../../lib/api';
import Card from '../../components/Card';
import Button from '../../components/Button';
import Spinner from '../../components/Spinner';
import Badge from '../../components/Badge';
import HighlightedText from '../../components/HighlightedText';
import PosLegend from '../../components/PosLegend';
import { formatDateTime, formatDuration, SECTION_LABEL, TYPE_LABEL, POS_LABEL } from '../../lib/format';
import type { Pos, Review } from '../../types';

const OPTION_KEYS = ['A', 'B', 'C', 'D'];

export default function ReviewPage() {
  const { id } = useParams<{ id: string }>();
  const attemptId = Number(id);
  const reviewQuery = useQuery({
    queryKey: ['attempts', attemptId, 'review'],
    queryFn: () => api.get<Review>(`/attempts/${attemptId}/review`),
    enabled: Number.isFinite(attemptId),
  });
  const [posFilter, setPosFilter] = useState<Set<Pos>>(new Set());
  const [wrongFirst, setWrongFirst] = useState(true);

  const review = reviewQuery.data;

  const presentPos = useMemo(() => {
    const s = new Set<Pos>();
    for (const it of review?.items ?? []) {
      for (const r of it.highlight_regions) s.add(r.pos);
    }
    return s;
  }, [review]);

  const items = useMemo(() => {
    if (!review) return [];
    let list = [...review.items];
    if (posFilter.size > 0) {
      list = list.filter((it) =>
        it.highlight_regions.some((r) => posFilter.has(r.pos)),
      );
    }
    list.sort((a, b) => {
      if (a.is_correct !== b.is_correct) return wrongFirst ? Number(a.is_correct) - Number(b.is_correct) : Number(b.is_correct) - Number(a.is_correct);
      return a.id - b.id;
    });
    return list;
  }, [review, posFilter, wrongFirst]);

  if (reviewQuery.isLoading) {
    return (
      <div className="flex items-center justify-center py-24">
        <Spinner className="h-8 w-8" />
      </div>
    );
  }

  if (reviewQuery.isError || !review) {
    return (
      <Card className="mx-auto max-w-reading p-8 text-center">
        <p className="text-danger">This review could not be loaded.</p>
        <Button variant="secondary" className="mt-4" onClick={() => window.history.back()}>
          Back
        </Button>
      </Card>
    );
  }

  const { attempt, report } = review;
  const sections = report.sections;

  return (
    <div className="mx-auto flex max-w-reading flex-col gap-6">
      <div className="print-hidden flex items-center justify-between">
        <Link to="/" className="inline-flex items-center gap-1.5 text-sm text-ink-muted transition-colors duration-200 hover:text-ink">
          <ArrowLeft className="h-4 w-4" aria-hidden="true" />
          Dashboard
        </Link>
        <Button variant="secondary" size="sm" onClick={() => window.print()}>
          <Printer className="h-4 w-4" aria-hidden="true" />
          Export PDF
        </Button>
      </div>

      <div className="print-header mb-4 border-b border-border pb-4">
        <p className="text-sm font-semibold text-primary">TOEFL Prep</p>
        <h1 className="text-xl font-bold text-ink">{attempt.exam_title ?? `Exam #${attempt.exam_template_id}`}</h1>
        <p className="text-sm text-ink-muted">
          {formatDateTime(attempt.started_at)}
          {attempt.finished_at && ` - ${formatDateTime(attempt.finished_at)}`}
        </p>
      </div>

      <Card className="p-6">
        <div className="flex flex-col items-center gap-4 text-center">
          <p className="text-sm font-medium uppercase tracking-wide text-ink-muted">Score</p>
          <p className="tabular text-6xl font-extrabold leading-none text-ink">{report.score_pct}%</p>
          <div className="flex flex-wrap items-center justify-center gap-2">
            <Badge tone="success">Correct {report.correct}</Badge>
            <Badge tone="danger">Wrong {report.wrong}</Badge>
            <Badge tone="neutral">Unanswered {report.unanswered}</Badge>
          </div>
          <div className="mt-1 flex flex-wrap items-center justify-center gap-4 text-[13px] text-ink-muted">
            <span>
              {report.total} questions
            </span>
            <span>
              {attempt.finished_at
                ? formatDuration(new Date(attempt.finished_at).getTime() - new Date(attempt.started_at).getTime())
                : '-'}
            </span>
          </div>
        </div>

        <div className="mt-6 flex flex-col gap-3">
          {sections.map((s) => (
            <div key={s.section} className="flex items-center gap-3">
              <span className="w-28 shrink-0 text-sm text-ink-muted">{SECTION_LABEL[s.section]}</span>
              <div className="h-3 flex-1 overflow-hidden rounded-full bg-card-muted" role="progressbar" aria-valuenow={s.score_pct} aria-valuemin={0} aria-valuemax={100}>
                <div className="h-full rounded-full bg-primary transition-all duration-300" style={{ width: `${s.score_pct}%` }} />
              </div>
              <span className="tabular w-12 text-right text-sm font-semibold text-ink">{s.score_pct}%</span>
            </div>
          ))}
        </div>
      </Card>

      <div className="flex flex-col gap-3">
        <p className="text-sm font-medium uppercase tracking-wide text-ink-muted">Legend</p>
        <PosLegend present={presentPos} detailed onToggle={(pos) => setPosFilter((prev) => {
          const next = new Set(prev);
          if (next.has(pos)) next.delete(pos);
          else next.add(pos);
          return next;
        })} />
        <label className="flex cursor-pointer items-center gap-2 text-sm text-ink-muted">
          <input type="checkbox" checked={wrongFirst} onChange={(e) => setWrongFirst(e.target.checked)} className="h-4 w-4 accent-primary" />
          Wrong answers first
        </label>
      </div>

      <div className="flex flex-col gap-4">
        {items.map((it) => {
          const chosenIndex = it.chosen_index;
          return (
            <div key={it.id} className={`print-break-avoid rounded-md border bg-card p-5 shadow-sm ${it.is_correct ? 'border-success/40' : 'border-danger/40'}`}>
              <div className="mb-3 flex flex-wrap items-center justify-between gap-2">
                <div className="flex items-center gap-2">
                  {it.is_correct ? (
                    <Badge tone="success">
                      <Check className="h-3 w-3" aria-hidden="true" /> Correct
                    </Badge>
                  ) : (
                    <Badge tone="danger">
                      <X className="h-3 w-3" aria-hidden="true" /> {it.is_unanswered ? 'Unanswered' : 'Wrong'}
                    </Badge>
                  )}
                  <span className="text-[13px] text-ink-muted">
                    {SECTION_LABEL[it.section]}
                    {it.part ? ` · ${it.part}` : ''} · {TYPE_LABEL[it.type]} · #{it.id}
                  </span>
                </div>
                <span className="tabular text-[13px] text-ink-muted">{formatDuration(it.time_taken_ms)}</span>
              </div>

              {it.passage ? (
                <div className="mb-4 max-h-64 overflow-y-auto rounded-sm border border-border bg-card-muted p-4 text-[15px] leading-relaxed text-ink">
                  {it.passage}
                </div>
              ) : null}

              <HighlightedText text={it.question_text} regions={it.highlight_regions} showLabels />

              <div className="mt-4 flex flex-col gap-2">
                {it.options.map((opt, i) => {
                  const isCorrect = i === it.correct_index;
                  const isChosen = i === chosenIndex;
                  let cls = 'border-border bg-card-muted/50 text-ink-muted';
                  if (isCorrect) cls = 'border-success bg-success-soft text-success';
                  else if (isChosen) cls = 'border-danger bg-danger-soft text-danger';
                  return (
                    <div key={i} className={`flex items-center gap-2 rounded-sm border px-3 py-2 text-sm ${cls}`}>
                      <span className="w-5 shrink-0 font-semibold">{OPTION_KEYS[i]}</span>
                      <span className="flex-1">{opt}</span>
                      {isCorrect && <span className="text-[12px] font-medium">Correct</span>}
                      {isChosen && !isCorrect && <span className="text-[12px] font-medium">Your answer</span>}
                    </div>
                  );
                })}
              </div>

              {it.explanation && (
                <div className="mt-4 rounded-sm bg-card-muted p-4">
                  <p className="mb-1 text-[12px] font-semibold uppercase tracking-wide text-ink-muted">Explanation</p>
                  <p className="text-[15px] leading-relaxed text-ink">{it.explanation}</p>
                </div>
              )}
            </div>
          );
        })}
        {items.length === 0 && (
          <Card className="p-8 text-center">
            <p className="text-sm text-ink-muted">
              No questions match the selected part of speech ({[...posFilter].map((p) => POS_LABEL[p]).join(', ')}).
            </p>
          </Card>
        )}
      </div>
    </div>
  );
}