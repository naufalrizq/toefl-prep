import { useMemo, useState } from 'react';
import { CheckCircle2, Pencil, Plus, Trash2, XCircle } from 'lucide-react';
import Card from '../../components/Card';
import Button from '../../components/Button';
import Badge from '../../components/Badge';
import Spinner from '../../components/Spinner';
import Input from '../../components/Input';
import Pagination from '../../components/Pagination';
import {
  useAdminExams,
  useCreateExam,
  useDeleteExam,
  usePublishExam,
  useQuestions,
  useUpdateExam,
} from './queries';
import { SECTION_LABEL, TYPE_LABEL, ALL_SECTIONS, ALL_TYPES, DEFAULT_TYPE } from '../../lib/format';
import type {
  ExamMode,
  ExamTemplate,
  ExamTemplateInput,
  PartConfig,
  QuestionType,
  Section,
  SectionConfig,
} from '../../types';

interface ExamDraft {
  id?: number;
  title: string;
  sections: Partial<Record<Section, SectionConfig>>;
  mode: ExamMode;
  secondsPerQuestion: string;
  totalMinutes: string;
  shuffle: boolean;
}

function fromExam(e: ExamTemplate): ExamDraft {
  return {
    id: e.id,
    title: e.title,
    sections: Object.fromEntries(
      ALL_SECTIONS.map((s) => [s, e.section_filters[s]]).filter(
        ([, cfg]) => cfg !== undefined,
      ) as [Section, SectionConfig][],
    ),
    mode: e.mode,
    secondsPerQuestion: e.seconds_per_question != null ? String(e.seconds_per_question) : '60',
    totalMinutes: e.total_minutes != null ? String(e.total_minutes) : '15',
    shuffle: e.shuffle ?? true,
  };
}

function empty(): ExamDraft {
  const sections: Partial<Record<Section, SectionConfig>> = {};
  for (const s of ALL_SECTIONS) {
    sections[s] = { parts: [{ title: 'Part 1', type: DEFAULT_TYPE[s], count: 5 }] };
  }
  return { title: '', sections, mode: 'per_question', secondsPerQuestion: '60', totalMinutes: '15', shuffle: true };
}

function toInput(d: ExamDraft): ExamTemplateInput {
  const section_filters: ExamTemplateInput['section_filters'] = {};
  for (const s of ALL_SECTIONS) {
    const cfg = d.sections[s];
    const parts = (cfg?.parts ?? []).filter((p) => p.count > 0);
    if (parts.length > 0) section_filters[s] = { parts };
  }
  return {
    title: d.title,
    section_filters,
    mode: d.mode,
    seconds_per_question: d.mode === 'per_question' || d.mode === 'both' ? Number(d.secondsPerQuestion) : undefined,
    total_minutes: d.mode === 'overall' || d.mode === 'both' ? Number(d.totalMinutes) : undefined,
    shuffle: d.shuffle,
  };
}

function totalFor(d: ExamDraft, s: Section): number {
  return (d.sections[s]?.parts ?? []).reduce((sum, p) => sum + p.count, 0);
}

function updatePart(
  d: ExamDraft,
  s: Section,
  index: number,
  patch: Partial<PartConfig>,
): ExamDraft {
  const cfg = d.sections[s] ?? { parts: [] };
  const parts = cfg.parts.map((p, i) => (i === index ? { ...p, ...patch } : p));
  return { ...d, sections: { ...d.sections, [s]: { parts } } };
}

export default function AdminExamsPage() {
  const exams = useAdminExams();
  const create = useCreateExam();
  const update = useUpdateExam();
  const remove = useDeleteExam();
  const publish = usePublishExam();
  const bank = useQuestions({ limit: 100, page: 1 });
  const [editor, setEditor] = useState<ExamDraft | null>(null);
  const [page, setPage] = useState(1);
  const PAGE_SIZE = 10;

  const allExams = exams.data ?? [];
  const totalPages = Math.max(1, Math.ceil(allExams.length / PAGE_SIZE));
  const current = Math.min(page, totalPages);
  const visible = allExams.slice((current - 1) * PAGE_SIZE, current * PAGE_SIZE);

  const bankCount = useMemo(() => {
    const c: Record<string, number> = {};
    for (const q of bank.data?.items ?? []) {
      if (q.active) c[`${q.section}/${q.type}`] = (c[`${q.section}/${q.type}`] ?? 0) + 1;
    }
    return c;
  }, [bank.data]);

  async function handleSave() {
    if (!editor) return;
    const input = toInput(editor);
    if (editor.id) await update.mutateAsync({ id: editor.id, body: input });
    else await create.mutateAsync(input);
    setEditor(null);
  }

  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <h1 className="text-2xl font-bold tracking-tight text-ink">Admin · Exams</h1>
        <Button size="sm" onClick={() => setEditor(empty())}>
          <Plus className="h-4 w-4" aria-hidden="true" />
          New exam
        </Button>
      </div>

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        <div className="flex flex-col gap-3">
          {exams.isLoading ? (
            <Card className="flex justify-center p-10">
              <Spinner className="h-6 w-6" />
            </Card>
          ) : exams.isError ? (
            <Card className="p-8 text-center">
              <p className="text-danger">Unable to load exams.</p>
            </Card>
          ) : (exams.data?.length ?? 0) === 0 ? (
            <Card className="p-8 text-center">
              <p className="text-sm text-ink-muted">No exam templates yet.</p>
            </Card>
          ) : (
            <>
              {visible.map((e) => (
              <Card key={e.id} className="p-5">
                <div className="flex items-start justify-between gap-3">
                  <div className="min-w-0">
                    <div className="flex flex-wrap items-center gap-2">
                      <h2 className="text-base font-semibold text-ink">{e.title}</h2>
                      <Badge tone={e.active ? (e.published ? 'success' : 'neutral') : 'danger'}>
                        {!e.active ? 'Inactive' : e.published ? 'Published' : 'Draft'}
                      </Badge>
                    </div>
                    <p className="mt-1 text-[13px] text-ink-muted">
                      {ALL_SECTIONS.filter((s) => e.section_filters[s]).map((s) => {
                        const parts = e.section_filters[s]?.parts ?? [];
                        const total = parts.reduce((sum, p) => sum + p.count, 0);
                        return `${SECTION_LABEL[s]} (${parts.length} part${parts.length === 1 ? '' : 's'}, ${total} q)`;
                      }).join(' · ') || 'No sections'}
                      {' · '}
                      {e.mode === 'both' ? 'Both modes' : e.mode === 'overall' ? 'Timed (overall)' : 'Per question'}
                      {e.seconds_per_question ? ` · ${e.seconds_per_question}s/q` : ''}
                      {e.total_minutes ? ` · ${e.total_minutes} min` : ''}
                    </p>
                  </div>
                  <div className="flex shrink-0 items-center gap-1">
                    <button
                      type="button"
                      onClick={() => setEditor(fromExam(e))}
                      className="rounded-sm p-2 text-ink-muted transition-colors duration-200 hover:bg-primary-soft hover:text-primary"
                      aria-label={`Edit exam ${e.id}`}
                    >
                      <Pencil className="h-4 w-4" aria-hidden="true" />
                    </button>
                    <button
                      type="button"
                      onClick={() => {
                        if (window.confirm(`Delete exam "${e.title}"?`)) void remove.mutate(e.id);
                      }}
                      className="rounded-sm p-2 text-ink-muted transition-colors duration-200 hover:bg-danger-soft hover:text-danger"
                      aria-label={`Delete exam ${e.id}`}
                    >
                      <Trash2 className="h-4 w-4" aria-hidden="true" />
                    </button>
                  </div>
                </div>
                {e.active && (
                  <button
                    type="button"
                    onClick={() => void publish.mutate({ id: e.id, published: !e.published })}
                    disabled={publish.isPending}
                    className={`mt-3 inline-flex items-center gap-1.5 rounded-sm px-3 py-1.5 text-[13px] font-medium transition-colors duration-200 ${
                      e.published
                        ? 'bg-danger-soft text-danger hover:bg-danger hover:text-white'
                        : 'bg-success-soft text-success hover:bg-success hover:text-white'
                    }`}
                  >
                    {e.published ? (
                      <>
                        <XCircle className="h-4 w-4" aria-hidden="true" /> Unpublish
                      </>
                    ) : (
                      <>
                        <CheckCircle2 className="h-4 w-4" aria-hidden="true" /> Publish
                      </>
                    )}
                  </button>
                )}
              </Card>
              ))}
              <Pagination page={current} totalPages={totalPages} onPageChange={setPage} />
            </>
          )}
        </div>

        {editor && (
          <Card className="p-5">
            <h2 className="mb-4 text-sm font-semibold uppercase tracking-wide text-ink-muted">
              {editor.id ? `Edit exam #${editor.id}` : 'New exam'}
            </h2>
            <form
              className="flex flex-col gap-4"
              onSubmit={(e) => {
                e.preventDefault();
                void handleSave();
              }}
            >
              <Input
                label="Title"
                value={editor.title}
                onChange={(e) => setEditor({ ...editor, title: e.target.value })}
                placeholder="Mock TOEFL A"
                required
              />
              <div className="flex flex-col gap-4">
                {ALL_SECTIONS.map((s) => {
                  const parts = editor.sections[s]?.parts ?? [];
                  return (
                    <div key={s} className="rounded-sm border border-border bg-card-muted p-3">
                      <div className="mb-2 flex items-center justify-between">
                        <span className="text-sm font-semibold text-ink">{SECTION_LABEL[s]}</span>
                        <span className="text-[13px] text-ink-muted">
                          {totalFor(editor, s)} questions · {parts.length} part{parts.length === 1 ? '' : 's'}
                        </span>
                      </div>
                      <div className="flex flex-col gap-2">
                        {parts.map((p, i) => (
                          <div key={i} className="grid grid-cols-[1fr_1.4fr_80px_auto] items-end gap-2">
                            <Input
                              label="Part title"
                              value={p.title}
                              onChange={(e) => setEditor(updatePart(editor, s, i, { title: e.target.value }))}
                            />
                            <label className="flex flex-col gap-1.5 text-sm font-medium text-ink">
                              Type
                              <select
                                value={p.type}
                                onChange={(e) =>
                                  setEditor(updatePart(editor, s, i, { type: e.target.value as QuestionType }))
                                }
                                className="rounded-sm border border-border-strong bg-card px-3 py-3 text-base"
                              >
                                {ALL_TYPES.map((t) => (
                                  <option key={t} value={t}>
                                    {TYPE_LABEL[t]}
                                  </option>
                                ))}
                              </select>
                            </label>
                            <Input
                              label="Count"
                              type="number"
                              min={0}
                              value={p.count}
                              onChange={(e) => setEditor(updatePart(editor, s, i, { count: Number(e.target.value) }))}
                            />
                            <button
                              type="button"
                              onClick={() =>
                                setEditor({
                                  ...editor,
                                  sections: { ...editor.sections, [s]: { parts: parts.filter((_, j) => j !== i) } },
                                })
                              }
                              className="rounded-sm p-2 text-ink-muted transition-colors duration-200 hover:bg-danger-soft hover:text-danger"
                              aria-label={`Remove ${s} part ${i + 1}`}
                            >
                              <Trash2 className="h-4 w-4" aria-hidden="true" />
                            </button>
                          </div>
                        ))}
                      </div>
                      <button
                        type="button"
                        onClick={() =>
                          setEditor({
                            ...editor,
                            sections: {
                              ...editor.sections,
                              [s]: {
                                parts: [...parts, { title: `Part ${parts.length + 1}`, type: DEFAULT_TYPE[s], count: 5 }],
                              },
                            },
                          })
                        }
                        className="mt-2 inline-flex items-center gap-1 text-[13px] font-medium text-primary hover:underline"
                      >
                        <Plus className="h-3.5 w-3.5" aria-hidden="true" /> Add part
                      </button>
                    </div>
                  );
                })}
              </div>
              <label className="flex flex-col gap-1.5 text-sm font-medium text-ink">
                Mode
                <select
                  value={editor.mode}
                  onChange={(e) => setEditor({ ...editor, mode: e.target.value as ExamMode })}
                  className="rounded-sm border border-border-strong bg-card px-3 py-3 text-base"
                >
                  <option value="per_question">Per question</option>
                  <option value="overall">Overall (timed)</option>
                  <option value="both">Both (student chooses)</option>
                </select>
              </label>
              <div className="grid grid-cols-2 gap-4">
                {(editor.mode === 'per_question' || editor.mode === 'both') && (
                  <Input
                    label="Seconds per question"
                    type="number"
                    min={10}
                    value={editor.secondsPerQuestion}
                    onChange={(e) => setEditor({ ...editor, secondsPerQuestion: e.target.value })}
                  />
                )}
                {(editor.mode === 'overall' || editor.mode === 'both') && (
                  <Input
                    label="Total minutes"
                    type="number"
                    min={1}
                    value={editor.totalMinutes}
                    onChange={(e) => setEditor({ ...editor, totalMinutes: e.target.value })}
                  />
                )}
              </div>
              <label className="flex cursor-pointer items-center gap-2 text-sm font-medium text-ink">
                <input
                  type="checkbox"
                  checked={editor.shuffle}
                  onChange={(e) => setEditor({ ...editor, shuffle: e.target.checked })}
                  className="h-4 w-4 accent-primary"
                />
                Shuffle questions
              </label>

              <div className="flex flex-col gap-1.5 rounded-sm bg-card-muted p-3 text-[13px]">
                <span className="font-medium text-ink">Question bank check</span>
                {bank.isLoading ? (
                  <span className="text-ink-muted">Checking question bank…</span>
                ) : (
                  ALL_SECTIONS.filter((s) => (editor.sections[s]?.parts ?? []).length > 0).map((s) => {
                    const grouped = new Map<QuestionType, number>();
                    for (const p of editor.sections[s]?.parts ?? []) {
                      grouped.set(p.type, (grouped.get(p.type) ?? 0) + p.count);
                    }
                    const rows = [...grouped.entries()].map(([t, needed]) => {
                      const available = bankCount[`${s}/${t}`] ?? 0;
                      const ok = available >= needed;
                      return (
                        <span key={t} className={`flex items-center gap-1.5 ${ok ? 'text-success' : 'text-danger'}`}>
                          {ok ? <CheckCircle2 className="h-4 w-4" aria-hidden="true" /> : <XCircle className="h-4 w-4" aria-hidden="true" />}
                          {TYPE_LABEL[t]}: {available} available, {needed} needed
                        </span>
                      );
                    });
                    return (
                      <div key={s}>
                        <span className="font-medium text-ink">{SECTION_LABEL[s]}</span>
                        <div className="flex flex-col gap-0.5">{rows}</div>
                      </div>
                    );
                  })
                )}
              </div>

              <div className="flex items-center justify-end gap-2 border-t border-border pt-4">
                <Button type="button" variant="ghost" onClick={() => setEditor(null)}>
                  Cancel
                </Button>
                <Button type="submit" loading={create.isPending || update.isPending}>
                  {editor.id ? 'Save changes' : 'Create exam'}
                </Button>
              </div>
            </form>
          </Card>
        )}
      </div>
    </div>
  );
}