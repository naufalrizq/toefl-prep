import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Play, Timer, Clock } from 'lucide-react';
import Card from '../../components/Card';
import Button from '../../components/Button';
import Spinner from '../../components/Spinner';
import { usePublishedExams, useStartAttempt } from '../dashboard/queries';
import { SECTION_LABEL } from '../../lib/format';
import type { Section, SectionConfig } from '../../types';

export default function ExamListPage() {
  const exams = usePublishedExams();
  const start = useStartAttempt();
  const navigate = useNavigate();
  const [modeFor, setModeFor] = useState<Record<number, 'per_question' | 'overall'>>({});

  async function handleStart(examId: number, availableModes: ('per_question' | 'overall')[]) {
    const mode = modeFor[examId] ?? availableModes[0];
    const result = await start.mutateAsync({ examTemplateId: examId, mode });
    navigate(`/quiz/${result.attempt.id}`);
  }

  if (exams.isLoading) {
    return (
      <div className="flex items-center justify-center py-24">
        <Spinner className="h-8 w-8" />
      </div>
    );
  }

  if (exams.isError) {
    return (
      <Card className="p-8 text-center">
        <p className="text-danger">Unable to load exams.</p>
      </Card>
    );
  }

  const list = (exams.data ?? []).filter((e) => e.published && e.active);

  return (
    <div className="mx-auto max-w-reading">
      <h1 className="mb-1 text-2xl font-bold tracking-tight text-ink">Start exam</h1>
      <p className="mb-6 text-sm text-ink-muted">Pick a published exam to begin.</p>

      {list.length === 0 ? (
        <Card className="p-8 text-center">
          <p className="text-sm text-ink-muted">No published exams available right now.</p>
        </Card>
      ) : (
        <div className="flex flex-col gap-4">
          {list.map((exam) => {
            const filters = Object.entries(exam.section_filters) as [Section, SectionConfig][];
            const modes: ('per_question' | 'overall')[] =
              exam.mode === 'both' ? ['per_question', 'overall'] : [exam.mode];
            const selected = modeFor[exam.id] ?? modes[0];
            const totalQuestions = filters.reduce(
              (acc, [, cfg]) => acc + (cfg?.parts ?? []).reduce((s, p) => s + p.count, 0),
              0,
            );
            const filterLabel = filters
              .map(([s, cfg]) => `${SECTION_LABEL[s]} ${(cfg?.parts ?? []).reduce((acc, p) => acc + p.count, 0)}`)
              .join(' · ');
            return (
              <Card key={exam.id} className="p-6">
                <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
                  <div>
                    <h2 className="text-lg font-semibold text-ink">{exam.title}</h2>
                    <p className="mt-1 flex flex-wrap items-center gap-x-3 gap-y-1 text-[13px] text-ink-muted">
                      <span>
                        {filterLabel} · {totalQuestions} questions
                      </span>
                      {exam.seconds_per_question ? (
                        <span className="inline-flex items-center gap-1">
                          <Timer className="h-3.5 w-3.5" aria-hidden="true" />
                          {exam.seconds_per_question}s/question
                        </span>
                      ) : null}
                      {exam.total_minutes ? (
                        <span className="inline-flex items-center gap-1">
                          <Clock className="h-3.5 w-3.5" aria-hidden="true" />
                          {exam.total_minutes} min
                        </span>
                      ) : null}
                    </p>
                    {modes.length === 2 && (
                      <div className="mt-3 inline-flex rounded-sm border border-border-strong p-0.5" role="radiogroup" aria-label="Mode">
                        {modes.map((m) => (
                          <button
                            key={m}
                            type="button"
                            role="radio"
                            aria-checked={selected === m}
                            onClick={() => setModeFor((s) => ({ ...s, [exam.id]: m }))}
                            className={`rounded-sm px-3 py-1.5 text-[13px] font-medium transition-colors duration-200 ${
                              selected === m ? 'bg-primary text-on-primary' : 'text-ink-muted hover:text-ink'
                            }`}
                          >
                            {m === 'overall' ? 'Timed (overall)' : 'Per question'}
                          </button>
                        ))}
                      </div>
                    )}
                  </div>
                  <Button
                    onClick={() => void handleStart(exam.id, modes)}
                    loading={start.isPending}
                    className="shrink-0"
                    disabled={start.isPending}
                  >
                    {start.isPending ? 'Preparing…' : (
                      <>
                        <Play className="h-4 w-4" aria-hidden="true" />
                        Start
                      </>
                    )}
                  </Button>
                </div>
              </Card>
            );
          })}
        </div>
      )}
    </div>
  );
}