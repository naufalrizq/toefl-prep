import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import { act, fireEvent, render, screen, within } from '@testing-library/react';
import QuizOverall from './QuizOverall';
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

function renderOverall(n = 2, overrides: Partial<Parameters<typeof QuizOverall>[0]> = {}) {
  const props = {
    items: makeItems(n),
    answers: {},
    deadlineMs: Date.now() + 60_000,
    onAnswer: vi.fn(),
    onFlag: vi.fn(),
    onSubmit: vi.fn(),
    submitting: false,
    ...overrides,
  };
  render(<QuizOverall {...props} />);
  return props;
}

describe('QuizOverall', () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });
  afterEach(() => {
    vi.useRealTimers();
  });

  it('renders the first question and the palette', () => {
    renderOverall(3);

    expect(screen.getByText('Question 1 of 3')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /opt A0/ })).toBeInTheDocument();
    expect(screen.getByLabelText('Question palette')).toBeInTheDocument();
    expect(screen.getByRole('timer')).toBeInTheDocument();
  });

  it('selecting an option answers without advancing', () => {
    const { onAnswer } = renderOverall(2);

    fireEvent.click(screen.getByRole('button', { name: /opt B0/ }));
    expect(onAnswer).toHaveBeenCalledWith(1, 1);
    expect(screen.getByText('Question 1 of 2')).toBeInTheDocument();
  });

  it('jumps to a question via the palette and flags it', () => {
    const { onFlag, items } = renderOverall(3);

    act(() => {
      fireEvent.click(screen.getByRole('button', { name: /Question 2,/ }));
    });
    expect(screen.getByText('Question 2 of 3')).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: /flag/i }));
    expect(onFlag).toHaveBeenCalledWith(items[1], true);
  });

  it('shows the confirm dialog with counts; Keep working cancels, Submit calls onSubmit', () => {
    const { onSubmit } = renderOverall(2, {
      answers: { 1: 0 },
      items: makeItems(2).map((it, i) => (i === 0 ? { ...it, flagged: true } : it)),
    });

    fireEvent.click(screen.getByRole('button', { name: /^submit$/i }));
    const dialog = screen.getByRole('dialog');
    expect(dialog).toHaveTextContent(/1 answered/i);
    expect(dialog).toHaveTextContent(/1 unanswered/i);
    expect(dialog).toHaveTextContent(/1 flagged/i);

    fireEvent.click(within(dialog).getByRole('button', { name: /keep working/i }));
    expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
    expect(onSubmit).not.toHaveBeenCalled();

    fireEvent.click(screen.getByRole('button', { name: /^submit$/i }));
    fireEvent.click(within(screen.getByRole('dialog')).getByRole('button', { name: /^submit$/i }));
    expect(onSubmit).toHaveBeenCalledTimes(1);
  });

  it('auto-submits when the overall timer expires', () => {
    const { onSubmit } = renderOverall(2, { deadlineMs: Date.now() + 1000 });

    expect(onSubmit).not.toHaveBeenCalled();
    act(() => {
      vi.advanceTimersByTime(1500);
    });
    expect(onSubmit).toHaveBeenCalledTimes(1);
  });
});
