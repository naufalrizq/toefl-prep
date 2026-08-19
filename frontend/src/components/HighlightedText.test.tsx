import { describe, expect, it } from 'vitest';
import { render, screen } from '@testing-library/react';
import HighlightedText from './HighlightedText';
import type { HighlightRegion } from '../types';

const regions: HighlightRegion[] = [
  { start: 4, end: 13, pos: 'noun' }, // "committee"
  { start: 14, end: 25, pos: 'verb' }, // "has decided"
];

describe('HighlightedText', () => {
  it('splits text into plain and highlighted spans', () => {
    const text = 'The committee has decided to postpone.';
    const { container } = render(<HighlightedText text={text} regions={regions} />);
    expect(container.querySelector('p')).toHaveTextContent(
      'The committee has decided to postpone.',
    );
    expect(screen.getByLabelText('Noun')).toHaveTextContent('committee');
    expect(screen.getByLabelText('Verb')).toHaveTextContent('has decided');
  });

  it('labels highlighted spans with their part of speech', () => {
    const text = 'The committee has decided to postpone.';
    render(<HighlightedText text={text} regions={regions} />);
    expect(screen.getByLabelText('Noun')).toHaveTextContent('committee');
    expect(screen.getByLabelText('Verb')).toHaveTextContent('has decided');
  });

  it('does not render overlapping text twice and clamps out-of-bounds regions', () => {
    const text = 'abc';
    const over: HighlightRegion[] = [
      { start: 0, end: 2, pos: 'noun' },
      { start: 1, end: 3, pos: 'verb' },
      { start: 99, end: 120, pos: 'other' },
    ];
    const { container } = render(<HighlightedText text={text} regions={over} />);
    expect(container.textContent).toBe('abc');
    expect(screen.getByLabelText('Noun')).toHaveTextContent('ab');
    expect(screen.getByLabelText('Verb')).toHaveTextContent('c');
  });

  it('ignores regions when the list is empty', () => {
    render(<HighlightedText text="plain sentence" regions={[]} />);
    expect(screen.getByText('plain sentence')).toBeInTheDocument();
  });
});