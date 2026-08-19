import { Link } from 'react-router-dom';
import { Area, AreaChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts';
import { ArrowRight, BookOpenCheck } from 'lucide-react';
import Card from '../../components/Card';
import StatCard from '../../components/StatCard';
import Spinner from '../../components/Spinner';
import Badge from '../../components/Badge';
import PosLegend from '../../components/PosLegend';
import { useDashboard, useAttempts } from './queries';
import { formatDate, formatScore, formatTimeRange, SECTION_LABEL, POS_LABEL } from '../../lib/format';
import type { Pos, Section } from '../../types';

function EmptyState() {
  return (
    <Card className="flex flex-col items-center gap-4 p-10 text-center">
      <BookOpenCheck className="h-12 w-12 text-ink-faint" aria-hidden="true" />
      <div>
        <h2 className="text-lg font-semibold text-ink">No attempts yet</h2>
        <p className="mt-1 text-sm text-ink-muted">
          Take your first exam and your progress will show up here.
        </p>
      </div>
      <Link
        to="/exams"
        className="inline-flex items-center gap-2 rounded-sm bg-primary px-5 py-2.5 text-sm font-medium text-on-primary transition-all duration-200 hover:-translate-y-px hover:bg-primary/90"
      >
        Start your first exam
        <ArrowRight className="h-4 w-4" aria-hidden="true" />
      </Link>
    </Card>
  );
}

export default function DashboardPage() {
  const dashboard = useDashboard();
  const attempts = useAttempts();

  if (dashboard.isLoading || attempts.isLoading) {
    return (
      <div className="flex items-center justify-center py-24">
        <Spinner className="h-8 w-8" />
      </div>
    );
  }

  if (dashboard.isError || attempts.isError) {
    return (
      <Card className="p-8 text-center">
        <p className="text-danger">Unable to load your progress.</p>
      </Card>
    );
  }

  const stats = dashboard.data!;
  const recent = attempts.data ?? [];

  if (stats.total_attempts === 0) {
    return (
      <div className="mx-auto max-w-reading">
        <EmptyState />
      </div>
    );
  }

  const chartData = stats.series.map((s) => ({
    date: formatDate(s.started_at),
    score: s.score_pct,
    id: s.id,
  }));

  const sections: Section[] = ['structure', 'vocabulary'];
  const presentPos = new Set(
    (Object.keys(stats.worst_pos ?? {}) as Pos[]).filter((p) => (stats.worst_pos?.[p] ?? 0) > 0),
  );
  const weakestPos = Object.entries(stats.worst_pos ?? {}).sort((a, b) => b[1] - a[1])[0];

  return (
    <div className="flex flex-col gap-6">
      <div className="grid grid-cols-2 gap-4 lg:grid-cols-3">
        <StatCard label="Exams taken" value={stats.total_attempts} />
        <StatCard label="Average score" value={formatScore(stats.average_score)} />
        <StatCard label="Best score" value={formatScore(stats.best_score)} />
        {/* Trend temporarily hidden from the dashboard. Re-enable when ready:
        <StatCard label="Trend" value={stats.trend >= 0 ? `+${stats.trend}` : stats.trend} delta={stats.trend} />
        */}
      </div>

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        <Card className="p-6">
          <h2 className="mb-4 text-sm font-semibold uppercase tracking-wide text-ink-muted">Score over time</h2>
          <div className="h-56">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={chartData} margin={{ top: 4, right: 8, left: -18, bottom: 0 }}>
                <defs>
                  <linearGradient id="scoreFill" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" style={{ stopColor: 'var(--primary)' }} stopOpacity={0.25} />
                    <stop offset="95%" style={{ stopColor: 'var(--primary)' }} stopOpacity={0} />
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" vertical={false} />
                <XAxis dataKey="date" tick={{ fontSize: 12 }} tickLine={false} axisLine={false} />
                <YAxis domain={[0, 100]} tick={{ fontSize: 12 }} tickLine={false} axisLine={false} tickFormatter={(v) => `${v}%`} />
                <Tooltip
                  formatter={(value: number) => [`${value}%`, 'Score']}
                  contentStyle={{ borderRadius: 8, border: '1px solid var(--border)', background: 'var(--card)', fontSize: 13, color: 'var(--ink)' }}
                />
                <Area
                  type="monotone"
                  dataKey="score"
                  style={{ stroke: 'var(--primary)' }}
                  strokeWidth={2}
                  fill="url(#scoreFill)"
                  connectNulls
                />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </Card>

        <div className="flex flex-col gap-6">
          <Card className="p-6">
            <h2 className="mb-4 text-sm font-semibold uppercase tracking-wide text-ink-muted">Section accuracy</h2>
            <div className="flex flex-col gap-4">
              {sections.map((s) => {
                const sec = stats.sections?.[s];
                if (!sec || sec.total === 0) return null;
                return (
                  <div key={s} className="flex items-center gap-3">
                    <span className="w-28 shrink-0 text-sm text-ink-muted">{SECTION_LABEL[s]}</span>
                    <div className="h-3 flex-1 overflow-hidden rounded-full bg-card-muted" role="progressbar" aria-valuenow={sec.accuracy} aria-valuemin={0} aria-valuemax={100}>
                      <div className="h-full rounded-full bg-primary transition-all duration-300" style={{ width: `${sec.accuracy}%` }} />
                    </div>
                    <span className="tabular w-12 text-right text-sm font-semibold text-ink">{sec.accuracy}%</span>
                  </div>
                );
              })}
            </div>
          </Card>

          <Card className="p-6">
            <h2 className="mb-3 text-sm font-semibold uppercase tracking-wide text-ink-muted">Weakest part of speech</h2>
            {weakestPos && weakestPos[1] > 0 ? (
              <>
                <p className="text-lg font-semibold text-ink">
                  {POS_LABEL[weakestPos[0] as Pos]}
                  <span className="ml-2 text-sm font-normal text-ink-muted">
                    {weakestPos[1]} mistake{weakestPos[1] > 1 ? 's' : ''}
                  </span>
                </p>
                <div className="mt-3">
                  <PosLegend present={presentPos} />
                </div>
              </>
            ) : (
              <p className="text-sm text-ink-muted">No wrong answers yet. Keep it up!</p>
            )}
          </Card>
        </div>
      </div>

      <Card className="p-6">
        <div className="mb-4 flex items-center justify-between">
          <h2 className="text-sm font-semibold uppercase tracking-wide text-ink-muted">Recent attempts</h2>
          <Link to="/history" className="text-sm font-medium text-primary hover:underline">
            View all
          </Link>
        </div>
        {recent.length === 0 ? (
          <p className="text-sm text-ink-muted">Nothing here yet.</p>
        ) : (
          <ul className="flex flex-col divide-y divide-border">
            {recent.slice(0, 6).map((a) => (
              <li key={a.id}>
                <Link
                  to={a.status === 'in_progress' ? `/quiz/${a.id}` : `/attempts/${a.id}/review`}
                  className="flex items-center justify-between gap-3 py-3 transition-colors duration-200 hover:bg-card-muted/60"
                >
                  <div className="flex items-center gap-3">
                    <span className="text-sm font-medium text-ink">{a.exam_title ?? `Exam #${a.exam_template_id}`}</span>
                    {a.status === 'in_progress' && <Badge tone="warning">In progress</Badge>}
                  </div>
                  <div className="flex items-center gap-3">
                    <span className="text-sm text-ink-muted">{formatTimeRange(a.started_at, a.finished_at)}</span>
                    {a.score_pct !== undefined && a.score_pct !== null && (
                      <span className="tabular text-sm font-semibold text-ink">{a.score_pct}%</span>
                    )}
                    <ArrowRight className="h-4 w-4 text-ink-faint" aria-hidden="true" />
                  </div>
                </Link>
              </li>
            ))}
          </ul>
        )}
      </Card>
    </div>
  );
}