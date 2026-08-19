import { useState } from 'react';
import { Database, FileJson2, Pencil, Plus } from 'lucide-react';
import Card from '../../components/Card';
import Button from '../../components/Button';
import Badge from '../../components/Badge';
import Spinner from '../../components/Spinner';
import Pagination from '../../components/Pagination';
import QuestionEditor, { emptyDraft, fromQuestion, type QuestionDraft } from './QuestionEditor';
import {
  useCreateQuestion,
  useDeleteQuestion,
  useImportDrafts,
  useQuestions,
  useSeed,
  useUpdateQuestion,
} from './queries';
import { SECTION_LABEL, TYPE_LABEL, ALL_SECTIONS } from '../../lib/format';
import type { Difficulty, QuestionType, Section } from '../../types';

export default function AdminQuestionsPage() {
  const [filters, setFilters] = useState<{
    section?: Section;
    type?: QuestionType;
    difficulty?: Difficulty;
    search?: string;
    page: number;
    size: number;
  }>({
    page: 1,
    size: 10,
  });
  const [editor, setEditor] = useState<{ draft: QuestionDraft; isNew: boolean } | null>(null);
  const [importOpen, setImportOpen] = useState(false);
  const [importText, setImportText] = useState('');

  const page = useQuestions({ ...filters, limit: filters.size });
  const create = useCreateQuestion();
  const update = useUpdateQuestion();
  const remove = useDeleteQuestion();
  const importDrafts = useImportDrafts();
  const seed = useSeed();

  const totalPages = Math.max(1, Math.ceil((page.data?.total ?? 0) / filters.size));

  async function handleSave() {
    if (!editor) return;
    const { draft } = editor;
    const body = {
      section: draft.section,
      type: draft.type,
      question_text: draft.question_text,
      passage: draft.passage,
      options: draft.options,
      correct_index: draft.correct_index,
      explanation: draft.explanation,
      difficulty: draft.difficulty,
      active: draft.active,
      highlight_regions: draft.highlight_regions,
    };
    if (draft.id) await update.mutateAsync({ id: draft.id, body });
    else await create.mutateAsync(body);
    setEditor(null);
  }

  async function handleImport() {
    let parsed: unknown;
    try {
      parsed = JSON.parse(importText);
    } catch {
      alert('Invalid JSON. Check the pasted content.');
      return;
    }
    const list = Array.isArray(parsed) ? parsed : [parsed];
    const results = await importDrafts.mutateAsync(list as never);
    const firstValid = results.find((r) => r.valid && r.question);
    if (firstValid?.question) {
      setEditor({ draft: fromQuestion(firstValid.question), isNew: true });
      setImportOpen(false);
      setImportText('');
    }
  }

  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <h1 className="text-2xl font-bold tracking-tight text-ink">Admin · Questions</h1>
        <div className="flex items-center gap-2">
          <Button variant="secondary" size="sm" onClick={() => setImportOpen((v) => !v)}>
            <FileJson2 className="h-4 w-4" aria-hidden="true" />
            Import JSON
          </Button>
          <Button variant="secondary" size="sm" onClick={() => void seed.mutate()} loading={seed.isPending}>
            <Database className="h-4 w-4" aria-hidden="true" />
            Seed
          </Button>
          <Button size="sm" onClick={() => setEditor({ draft: emptyDraft(), isNew: true })}>
            <Plus className="h-4 w-4" aria-hidden="true" />
            New
          </Button>
        </div>
      </div>

      <div className="flex flex-wrap items-center gap-2">
        <select
          value={filters.section ?? ''}
          onChange={(e) => setFilters((f) => ({ ...f, section: (e.target.value as Section) || undefined, page: 1 }))}
          className="rounded-sm border border-border-strong bg-card px-3 py-2 text-sm"
          aria-label="Filter by section"
        >
          <option value="">All sections</option>
          {ALL_SECTIONS.map((s) => (
            <option key={s} value={s}>
              {SECTION_LABEL[s]}
            </option>
          ))}
        </select>
        <select
          value={filters.difficulty ?? ''}
          onChange={(e) => setFilters((f) => ({ ...f, difficulty: (e.target.value as Difficulty) || undefined, page: 1 }))}
          className="rounded-sm border border-border-strong bg-card px-3 py-2 text-sm"
          aria-label="Filter by difficulty"
        >
          <option value="">All difficulties</option>
          <option value="easy">Easy</option>
          <option value="medium">Medium</option>
          <option value="hard">Hard</option>
        </select>
        <input
          type="search"
          value={filters.search ?? ''}
          onChange={(e) => setFilters((f) => ({ ...f, search: e.target.value || undefined, page: 1 }))}
          placeholder="Search text…"
          className="w-56 rounded-sm border border-border-strong bg-card px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
          aria-label="Search questions"
        />
      </div>

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        <div className="flex flex-col gap-3">
          {page.isLoading ? (
            <Card className="flex justify-center p-10">
              <Spinner className="h-6 w-6" />
            </Card>
          ) : page.isError ? (
            <Card className="p-8 text-center">
              <p className="text-danger">Unable to load questions.</p>
            </Card>
          ) : (page.data?.items.length ?? 0) === 0 ? (
            <Card className="p-8 text-center">
              <p className="text-sm text-ink-muted">No questions match the current filters.</p>
              <Button variant="secondary" size="sm" className="mt-3" onClick={() => setEditor({ draft: emptyDraft(), isNew: true })}>
                Add your first question
              </Button>
            </Card>
          ) : (
            <div className="flex flex-col gap-2">
              {(page.data?.items ?? []).map((q) => (
                <Card key={q.id} interactive className="p-4" onClick={() => setEditor({ draft: fromQuestion(q), isNew: false })}>
                  <div className="flex items-center justify-between gap-3">
                    <div className="min-w-0">
                      <p className="truncate text-sm font-medium text-ink">{q.question_text}</p>
                      <div className="mt-1 flex flex-wrap items-center gap-1.5 text-[12px] text-ink-muted">
                        <Badge tone="neutral">{SECTION_LABEL[q.section]}</Badge>
                        <span>{TYPE_LABEL[q.type]}</span>
                        <span className="capitalize">{q.difficulty}</span>
                        {!q.active && <Badge tone="danger">Inactive</Badge>}
                      </div>
                    </div>
                    <div className="flex shrink-0 items-center gap-1">
                      <button
                        type="button"
                        onClick={() => setEditor({ draft: fromQuestion(q), isNew: false })}
                        className="rounded-sm p-2 text-ink-muted transition-colors duration-200 hover:bg-primary-soft hover:text-primary"
                        aria-label={`Edit question ${q.id}`}
                      >
                        <Pencil className="h-4 w-4" aria-hidden="true" />
                      </button>
                      <button
                        type="button"
                        onClick={() => {
                          if (window.confirm(`Delete question ${q.id}?`)) void remove.mutate(q.id);
                        }}
                        className="rounded-sm p-2 text-ink-muted transition-colors duration-200 hover:bg-danger-soft hover:text-danger"
                        aria-label={`Delete question ${q.id}`}
                      >
                        ✕
                      </button>
                    </div>
                  </div>
                </Card>
              ))}
              <div className="mt-2">
                <Pagination
                  page={filters.page}
                  totalPages={totalPages}
                  size={filters.size}
                  onSizeChange={(size) => setFilters((f) => ({ ...f, size, page: 1 }))}
                  onPageChange={(p) => setFilters((f) => ({ ...f, page: p }))}
                />
              </div>
            </div>
          )}
        </div>

        <div className="flex flex-col gap-4">
          {importOpen && (
            <Card className="p-5">
              <h2 className="mb-2 text-sm font-semibold uppercase tracking-wide text-ink-muted">Import from AI</h2>
              <p className="mb-3 text-[13px] text-ink-muted">
                Paste the JSON output of the question generator (run it outside the app). Valid items become drafts for you to review.
              </p>
              <textarea
                rows={8}
                value={importText}
                onChange={(e) => setImportText(e.target.value)}
                placeholder='[{"section":"structure","question_text":"...","options":["...","...","...","..."],"correct_index":0,"explanation":"...","highlights":{"verb":["..."]}}]'
                className="w-full rounded-sm border border-border-strong bg-card px-3 py-2.5 font-mono text-[13px] leading-relaxed focus:outline-none focus:ring-2 focus:ring-primary/30"
                aria-label="Pasted AI JSON"
              />
              <div className="mt-3 flex items-center justify-end gap-2">
                <Button variant="ghost" size="sm" onClick={() => setImportOpen(false)}>
                  Cancel
                </Button>
                <Button size="sm" onClick={() => void handleImport()} loading={importDrafts.isPending}>
                  Validate & import
                </Button>
              </div>
              {importDrafts.data && (
                <div className="mt-3 flex flex-col gap-1">
                  {importDrafts.data.map((r) => (
                    <p key={r.index} className={`text-[13px] ${r.valid ? 'text-success' : 'text-danger'}`}>
                      Draft #{r.index}: {r.valid ? 'valid, opened in editor' : r.error ?? 'invalid'}
                    </p>
                  ))}
                </div>
              )}
              {seed.data && (
                <p className="mt-2 text-[13px] text-success">Seeded {seed.data.seeded} questions.</p>
              )}
            </Card>
          )}

          {editor && (
            <Card className="p-5">
              <h2 className="mb-4 text-sm font-semibold uppercase tracking-wide text-ink-muted">
                {editor.isNew ? 'New question' : `Edit question #${editor.draft.id}`}
              </h2>
              <QuestionEditor
                draft={editor.draft}
                onChange={(d) => setEditor({ ...editor, draft: d })}
                onSubmit={() => void handleSave()}
                submitting={create.isPending || update.isPending}
                submitLabel={editor.isNew ? 'Create question' : 'Save changes'}
                onCancel={() => setEditor(null)}
              />
            </Card>
          )}
        </div>
      </div>
    </div>
  );
}