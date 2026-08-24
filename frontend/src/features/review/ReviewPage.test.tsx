import { describe, expect, it, vi, beforeAll } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import ReviewPage from './ReviewPage';
import { api } from '../../lib/api';
import type { Review } from '../../types';

vi.mock('../../lib/api', () => ({
  api: { get: vi.fn() },
}));

const review: Review = {
  attempt: {
    id: 7,
    user_id: 1,
    exam_template_id: 2,
    exam_title: 'Mock TOEFL A',
    mode: 'overall',
    started_at: '2026-08-24T10:00:00Z',
    finished_at: '2026-08-24T10:15:00Z',
    status: 'submitted',
    score_pct: 50,
  },
  report: {
    score_pct: 50,
    correct: 1,
    wrong: 1,
    unanswered: 0,
    total: 2,
    sections: [{ section: 'structure', correct: 1, total: 2, score_pct: 50 }],
  },
  items: [
    {
      id: 11,
      section: 'structure',
      type: 'sentence-completion',
      question_text: 'The committee has decided the policy.',
      options: ['has', 'have', 'having', 'to have'],
      correct_index: 0,
      chosen_index: 0,
      is_correct: true,
      is_unanswered: false,
      explanation: 'Subjek committee bersifat kolektif tunggal.',
      highlight_regions: [{ start: 14, end: 17, pos: 'verb', label: 'Verb' }],
      time_taken_ms: 12_000,
    },
    {
      id: 12,
      section: 'structure',
      type: 'sentence-completion',
      question_text: 'She rarely goes to the library.',
      options: ['rarely', 'goes', 'to', 'library'],
      correct_index: 1,
      chosen_index: 3,
      is_correct: false,
      is_unanswered: false,
      explanation: 'Setelah rarely, verbo tetap conjugated.',
      highlight_regions: [],
      time_taken_ms: 9_000,
    },
  ],
};

function renderReview() {
  vi.mocked(api.get).mockResolvedValue(review);
  const qc = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={qc}>
      <MemoryRouter initialEntries={['/attempts/7/review']}>
        <Routes>
          <Route path="/attempts/:id/review" element={<ReviewPage />} />
        </Routes>
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe('ReviewPage', () => {
  beforeAll(() => {
    window.print = vi.fn();
  });

  it('renders the score summary with counts and section bars', async () => {
    renderReview();

    expect(await screen.findByText('Mock TOEFL A')).toBeInTheDocument();
    expect(screen.getByText('50%', { selector: 'p.text-6xl' })).toBeInTheDocument();
    expect(screen.getByText(/correct\s*1/i)).toBeInTheDocument();
    expect(screen.getByText(/wrong\s*1/i)).toBeInTheDocument();
    expect(screen.getByRole('progressbar')).toHaveAttribute('aria-valuenow', '50');
  });

  it('renders highlights on the question text and marks chosen vs correct answers', async () => {
    renderReview();

    const mark = await screen.findByText('has', { selector: 'span[title="Verb"]' });
    expect(mark).toBeInTheDocument();

    // Wrong item shows "Your answer" on D; correct options are labeled.
    expect(screen.getByText('Your answer')).toBeInTheDocument();
    expect(screen.getAllByText('Explanation')).toHaveLength(2);
  });

  it('sorts wrong answers first by default', async () => {
    renderReview();

    await screen.findByText('50%', { selector: 'p.text-6xl' });
    const wrongCard = screen.getByText('Setelah rarely, verbo tetap conjugated.');
    const correctCard = screen.getByText('Subjek committee bersifat kolektif tunggal.');
    expect(
      wrongCard.compareDocumentPosition(correctCard) & Node.DOCUMENT_POSITION_FOLLOWING,
    ).toBeTruthy();
  });

  it('shows the print header and Export PDF triggers window.print()', async () => {
    const printSpy = vi.spyOn(window, 'print').mockImplementation(() => {});
    renderReview();

    expect(await screen.findByText('Export PDF')).toBeInTheDocument();
    expect(document.querySelector('.print-header')).toBeInTheDocument();

    await userEvent.click(screen.getByRole('button', { name: /export pdf/i }));
    expect(printSpy).toHaveBeenCalledTimes(1);
    printSpy.mockRestore();
  });
});
