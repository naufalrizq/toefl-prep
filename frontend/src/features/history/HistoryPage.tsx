import { Link } from 'react-router-dom';
import { ArrowRight } from 'lucide-react';
import Card from '../../components/Card';
import Badge from '../../components/Badge';
import Spinner from '../../components/Spinner';
import { useAttempts } from '../dashboard/queries';
import { formatDuration, formatTimeRange, SECTION_LABEL } from '../../lib/format';
import type { Section } from '../../types';

export default function HistoryPage() {
  const { data, isLoading, isError } = useAttempts();

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-24">
        <Spinner className="h-8 w-8" />
      </div>
    );
  }

  if (isError) {
    return (
      <Card className="p-8 text-center">
        <p className="text-danger">Unable to load your attempts.</p>
      </Card>
    );
  }

  const attempts = data ?? [];

  return (
    <div className="mx-auto max-w-reading">
      <h1 className="mb-6 text-2xl font-bold tracking-tight text-ink">History</h1>
      {attempts.length === 0 ? (
        <Card className="p-8 text-center">
          <p className="text-sm text-ink-muted">No attempts yet.</p>
          <Link to="/exams" className="mt-3 inline-block text-sm font-medium text-primary hover:underline">
            Start your first exam →
          </Link>
        </Card>
      ) : (
        <Card className="divide-y divide-border">
          {attempts.map((a) => {
            const inProgress = a.status === 'in_progress';
            const sections = (a.summary?.sections ?? []) as Array<{ section: Section; score_pct: number }>;
            const duration =
              a.finished_at
                ? formatDuration(new Date(a.finished_at).getTime() - new Date(a.started_at).getTime())
                : undefined;
            return (
              <Link
                key={a.id}
                to={inProgress ? `/quiz/${a.id}` : `/attempts/${a.id}/review`}
                className="block px-6 py-4 transition-colors duration-200 hover:bg-card-muted/60"
              >
                <div className="flex items-center justify-between gap-3">
                  <div className="min-w-0">
                    <p className="truncate font-medium text-ink">{a.exam_title ?? `Exam #${a.exam_template_id}`}</p>
                    <p className="text-[13px] text-ink-muted">
                      {a.mode === 'overall' ? 'Timed' : 'Per-question'}
                      {inProgress ? ' · started ' : ' · '}
                      {formatTimeRange(a.started_at, a.finished_at ?? undefined)}
                      {duration && ` · ${duration}`}
                    </p>
                    {sections.length > 0 && (
                      <div className="mt-1 flex flex-wrap gap-1.5">
                        {sections.map((s) => (
                          <span key={s.section} className="rounded-full bg-card-muted px-2 py-0.5 text-[12px] font-medium text-ink-muted">
                            {SECTION_LABEL[s.section]} {s.score_pct}%
                          </span>
                        ))}
                      </div>
                    )}
                  </div>
                  <div className="flex shrink-0 items-center gap-3">
                    {inProgress ? (
                      <Badge tone="warning">In progress</Badge>
                    ) : (
                      <span className="tabular text-lg font-bold text-ink">{a.score_pct}%</span>
                    )}
                    <ArrowRight className="h-4 w-4 text-ink-faint" aria-hidden="true" />
                  </div>
                </div>
              </Link>
            );
          })}
        </Card>
      )}
    </div>
  );
}