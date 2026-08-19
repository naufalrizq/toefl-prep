import { useCallback, useMemo, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useMutation } from '@tanstack/react-query';
import { api } from '../../lib/api';
import Spinner from '../../components/Spinner';
import Card from '../../components/Card';
import Button from '../../components/Button';
import QuizPerQuestion from './QuizPerQuestion';
import QuizOverall from './QuizOverall';
import { useAttempt, useQuizItems, useRecordAnswer, useToggleFlag } from './queries';
import type { QuizItem, Report } from '../../types';

export default function QuizPage() {
  const { id } = useParams<{ id: string }>();
  const attemptId = Number(id);
  const navigate = useNavigate();
  const attempt = useAttempt(attemptId);
  const itemsQuery = useQuizItems(attemptId);
  const record = useRecordAnswer(attemptId);
  const toggleFlag = useToggleFlag(attemptId);
  const [answers, setAnswers] = useState<Record<number, number | null>>({});
  const [flags, setFlags] = useState<Record<number, boolean>>({});

  const items = useMemo(
    () => (itemsQuery.data ?? []).map((it) => ({ ...it, flagged: flags[it.id] ?? it.flagged })),
    [itemsQuery.data, flags],
  );

  const submit = useMutation({
    mutationFn: () => api.post<Report>(`/attempts/${attemptId}/submit`),
    onSuccess: (report) => {
      navigate(`/attempts/${attemptId}/review`, { state: { report }, replace: true });
    },
  });

  const handleAnswer = useCallback(
    (itemId: number, chosenIndex: number | null, timeTakenMs?: number) => {
      setAnswers((prev) => ({ ...prev, [itemId]: chosenIndex }));
      void record.mutate({ itemId, chosenIndex, timeTakenMs });
    },
    [record],
  );

  const handleFlag = useCallback(
    (item: QuizItem, flagged: boolean) => {
      setFlags((prev) => ({ ...prev, [item.id]: flagged }));
      void toggleFlag.mutate({ itemId: item.id, flagged });
    },
    [toggleFlag],
  );

  if (attempt.isLoading || itemsQuery.isLoading) {
    return (
      <div className="flex items-center justify-center py-24">
        <Spinner className="h-8 w-8" />
      </div>
    );
  }

  if (itemsQuery.isError || attempt.isError || (!attempt.data && !attempt.isLoading)) {
    return (
      <Card className="mx-auto max-w-reading p-8 text-center">
        <p className="text-danger">This attempt could not be loaded.</p>
        <Button variant="secondary" className="mt-4" onClick={() => navigate('/')}>
          Back to dashboard
        </Button>
      </Card>
    );
  }

  if (!attempt.data) {
    return (
      <Card className="mx-auto max-w-reading p-8 text-center">
        <p className="text-ink-muted">Attempt not found.</p>
      </Card>
    );
  }

  const attemptInfo = attempt.data;
  const deadlineMs = attemptInfo.deadline ? new Date(attemptInfo.deadline).getTime() : Date.now() + 60 * 60 * 1000;

  if (items.length === 0) {
    return (
      <Card className="mx-auto max-w-reading p-8 text-center">
        <p className="text-ink-muted">No questions in this attempt.</p>
      </Card>
    );
  }

  if (attemptInfo.mode === 'overall') {
    return (
      <div className="py-2">
        <QuizOverall
          items={items}
          answers={answers}
          deadlineMs={deadlineMs}
          onAnswer={(itemId, i) => handleAnswer(itemId, i)}
          onFlag={handleFlag}
          onSubmit={() => {
            if (!submit.isPending) void submit.mutate();
          }}
          submitting={submit.isPending}
        />
      </div>
    );
  }

  return (
    <div className="py-2">
      <QuizPerQuestion
        items={items}
        answers={answers}
        secondsPerQuestion={60}
        onAnswer={handleAnswer}
        onFlag={handleFlag}
        onSubmit={() => {
          if (!submit.isPending) void submit.mutate();
        }}
        submitting={submit.isPending}
      />
    </div>
  );
}