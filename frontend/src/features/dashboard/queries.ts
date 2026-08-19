import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { api } from '../../lib/api';
import type {
  Attempt,
  DashboardStats,
  ExamTemplate,
  StartResult,
} from '../../types';

export function useDashboard() {
  return useQuery({
    queryKey: ['dashboard', 'stats'],
    queryFn: () => api.get<DashboardStats>('/dashboard/stats'),
  });
}

export function useAttempts() {
  return useQuery({
    queryKey: ['attempts'],
    queryFn: () => api.get<Attempt[]>('/attempts'),
  });
}

export function usePublishedExams() {
  return useQuery({
    queryKey: ['exams', 'published'],
    queryFn: () => api.get<ExamTemplate[]>('/exams'),
  });
}

export function useStartAttempt() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ examTemplateId, mode }: { examTemplateId: number; mode: 'per_question' | 'overall' }) =>
      api.post<StartResult>('/attempts', { exam_template_id: examTemplateId, mode }),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: ['attempts'] });
    },
  });
}