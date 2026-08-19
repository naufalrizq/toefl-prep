import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { api } from '../../lib/api';
import type {
  Difficulty,
  ExamTemplate,
  ExamTemplateInput,
  ImportResult,
  Draft,
  Question,
  QuestionPage,
  QuestionType,
  Section,
} from '../../types';

export interface QuestionFilters {
  section?: Section;
  type?: QuestionType;
  difficulty?: Difficulty;
  search?: string;
  page?: number;
  limit?: number;
}

export function useQuestions(filters: QuestionFilters) {
  const params = new URLSearchParams();
  if (filters.section) params.set('section', filters.section);
  if (filters.type) params.set('type', filters.type);
  if (filters.difficulty) params.set('difficulty', filters.difficulty);
  if (filters.search) params.set('search', filters.search);
  params.set('page', String(filters.page ?? 1));
  params.set('limit', String(filters.limit ?? 50));
  return useQuery({
    queryKey: ['questions', filters],
    queryFn: () => api.get<QuestionPage>(`/questions?${params.toString()}`),
  });
}

export function useCreateQuestion() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: Partial<Question>) => api.post<Question>('/questions', body),
    onSuccess: () => void qc.invalidateQueries({ queryKey: ['questions'] }),
  });
}

export function useUpdateQuestion() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, body }: { id: number; body: Partial<Question> }) =>
      api.put<Question>(`/questions/${id}`, body),
    onSuccess: () => void qc.invalidateQueries({ queryKey: ['questions'] }),
  });
}

export function useDeleteQuestion() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => api.del<void>(`/questions/${id}`),
    onSuccess: () => void qc.invalidateQueries({ queryKey: ['questions'] }),
  });
}

export function useImportDrafts() {
  return useMutation({
    mutationFn: (drafts: Draft[]) => api.post<ImportResult[]>('/questions/import', drafts),
  });
}

export function useSeed() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: () => api.post<{ seeded: number }>('/seed'),
    onSuccess: () => void qc.invalidateQueries({ queryKey: ['questions'] }),
  });
}

export function useAdminExams() {
  return useQuery({
    queryKey: ['exams', 'admin'],
    queryFn: () => api.get<ExamTemplate[]>('/exams'),
  });
}

export function useCreateExam() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: ExamTemplateInput) => api.post<ExamTemplate>('/exams', body),
    onSuccess: () => void qc.invalidateQueries({ queryKey: ['exams', 'admin'] }),
  });
}

export function useUpdateExam() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, body }: { id: number; body: ExamTemplateInput }) =>
      api.put<ExamTemplate>(`/exams/${id}`, body),
    onSuccess: () => void qc.invalidateQueries({ queryKey: ['exams', 'admin'] }),
  });
}

export function useDeleteExam() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => api.del<void>(`/exams/${id}`),
    onSuccess: () => void qc.invalidateQueries({ queryKey: ['exams', 'admin'] }),
  });
}

export function usePublishExam() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, published }: { id: number; published: boolean }) =>
      api.post<void>(`/exams/${id}/publish`, { published }),
    onSuccess: () => void qc.invalidateQueries({ queryKey: ['exams', 'admin'] }),
  });
}