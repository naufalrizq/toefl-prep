import { useMutation, useQuery } from '@tanstack/react-query';
import { api } from '../../lib/api';
import type { QuizItem, Attempt } from '../../types';

export function useAttempt(id: number) {
  return useQuery({
    queryKey: ['attempts', 'list', 'for-quiz'],
    queryFn: () => api.get<Attempt[]>('/attempts'),
    select: (list) => list.find((a) => a.id === id),
    enabled: Number.isFinite(id),
  });
}

export function useQuizItems(id: number) {
  return useQuery({
    queryKey: ['attempts', id, 'questions'],
    queryFn: () => api.get<QuizItem[]>(`/attempts/${id}/questions`),
    enabled: Number.isFinite(id),
  });
}

export function useRecordAnswer(attemptId: number) {
  return useMutation({
    mutationFn: ({ itemId, chosenIndex, timeTakenMs }: { itemId: number; chosenIndex: number | null; timeTakenMs?: number }) =>
      api.put<void>(`/attempts/${attemptId}/answers/${itemId}`, {
        chosen_index: chosenIndex,
        time_taken_ms: timeTakenMs,
      }),
  });
}

export function useToggleFlag(attemptId: number) {
  return useMutation({
    mutationFn: ({ itemId, flagged }: { itemId: number; flagged: boolean }) =>
      api.put<void>(`/attempts/${attemptId}/flag/${itemId}`, { flagged }),
  });
}