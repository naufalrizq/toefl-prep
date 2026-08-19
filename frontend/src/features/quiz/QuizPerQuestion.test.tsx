import { useState } from 'react';
import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import { act, fireEvent, render, screen } from '@testing-library/react';
import QuizPerQuestion from './QuizPerQuestion';
import type { QuizItem } from '../../types';

function makeItems(n: number): QuizItem[] {
  return Array.from({ length: n }, (_, i) => ({
    id: i + 1,
    question_snapshot: {
      question_text: `Question ${i + 1} with a ____ blank`,
      options: [`opt A${i}`, `opt B${i}`, `opt C${i}`, `opt D${i}`],
      explanation: '',
      highlight_regions: [],
      section: 'structure',
      type: 'sentence-completion',
    },
    flagged: false,
  }));
}

function renderQuiz(n: number, overrides: Partial<Parameters<typeof QuizPerQuestion>[0]> = {}) {
  const props = {
    items: makeItems(n),
    answers: {},
    secondsPerQuestion: 60,
    onAnswer: vi.fn(),
    onFlag: vi.fn(),
    onSubmit: vi.fn(),
    submitting: false,
    ...overrides,
  };
  render(<QuizPerQuestion {...props} />);
  return props;
}

function StatefulQuiz({ n = 2 }: { n?: number }) {
  const [answers, setAnswers] = useState<Record<number, number | null>>({});
  return (
    <QuizPerQuestion
      items={makeItems(n)}
      answers={answers}
      secondsPerQuestion={60}
      onAnswer={(id, chosen) => setAnswers((prev) => ({ ...prev, [id]: chosen }))}
      onFlag={() => {}}
      onSubmit={() => {}}
      submitting={false}
    />
  );
}

describe('QuizPerQuestion', () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });
  afterEach(() => {
    vi.useRealTimers();
  });

  it('auto-advances to the next question when an option is selected', () => {
    const { onAnswer, onSubmit } = renderQuiz(2);

    expect(screen.getByText('Q 1 / 2')).toBeInTheDocument();
    fireEvent.click(screen.getByRole('option', { name: /opt A0/ }));
    expect(onAnswer).toHaveBeenCalledWith(1, 0, expect.any(Number));
    expect(screen.getByText('Q 2 / 2')).toBeInTheDocument();
    expect(onSubmit).not.toHaveBeenCalled();
  });

  it('answers via the A–D keyboard shortcuts', () => {
    const { onAnswer } = renderQuiz(1);
    fireEvent.keyDown(window, { key: 'c' });
    expect(onAnswer).toHaveBeenCalledWith(1, 2, expect.any(Number));
  });

  it('marks the question unanswered and advances when the timer expires', () => {
    const { onAnswer } = renderQuiz(2, { secondsPerQuestion: 2 });

    act(() => {
      vi.advanceTimersByTime(2100);
    });

    expect(onAnswer).toHaveBeenCalledWith(1, null, expect.any(Number));
    expect(screen.getByText('Q 2 / 2')).toBeInTheDocument();
  });

  it('does not auto-submit after the last question, only via Finish & submit', () => {
    const { onAnswer, onSubmit } = renderQuiz(1);

    fireEvent.click(screen.getByRole('option', { name: /opt B0/ }));
    expect(onAnswer).toHaveBeenCalledWith(1, 1, expect.any(Number));
    expect(onSubmit).not.toHaveBeenCalled();

    fireEvent.click(screen.getByRole('button', { name: /finish & submit/i }));
    expect(onSubmit).toHaveBeenCalledTimes(1);
  });

  it('goes back to a previous question and restores the chosen answer', () => {
    const onAnswer = vi.fn();
    render(
      <QuizPerQuestion
        items={makeItems(2)}
        answers={{ 1: 0 }}
        secondsPerQuestion={60}
        onAnswer={onAnswer}
        onFlag={() => {}}
        onSubmit={() => {}}
        submitting={false}
      />,
    );

    fireEvent.click(screen.getByRole('button', { name: /previous/i }));
    expect(screen.getByText('Q 1 / 2')).toBeInTheDocument();
    expect(screen.getByRole('option', { name: /opt A0/ })).toHaveAttribute('aria-selected', 'true');
    expect(onAnswer).not.toHaveBeenCalled();
  });

  it('keeps an answered question selected when answering then going back', () => {
    render(<StatefulQuiz n={2} />);

    fireEvent.click(screen.getByRole('option', { name: /opt A0/ }));
    expect(screen.getByText('Q 2 / 2')).toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: /previous/i }));
    expect(screen.getByRole('option', { name: /opt A0/ })).toHaveAttribute('aria-selected', 'true');
  });

  it('jumps to any question from the palette', () => {
    renderQuiz(3);
    fireEvent.click(screen.getByRole('button', { name: 'Question 3, unanswered' }));
    expect(screen.getByText('Q 3 / 3')).toBeInTheDocument();
  });

  it('flags and unflags the current question', () => {
    const { onFlag } = renderQuiz(1);
    fireEvent.click(screen.getByRole('button', { name: /flag/i }));
    expect(onFlag).toHaveBeenCalledWith(expect.objectContaining({ id: 1 }), true);
  });
});