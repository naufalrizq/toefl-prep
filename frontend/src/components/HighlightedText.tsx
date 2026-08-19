import { useMemo } from 'react';
import type { HighlightRegion, Pos } from '../types';
import { POS_LABEL } from '../lib/format';

const POS_TEXT: Record<Pos, string> = {
  verb: 'text-pos-verb',
  noun: 'text-pos-noun',
  pronoun: 'text-pos-pronoun',
  adjective: 'text-pos-adjective',
  adverb: 'text-pos-adverb',
  preposition: 'text-pos-preposition',
  conjunction: 'text-pos-conjunction',
  determiner: 'text-pos-determiner',
  other: 'text-pos-other',
};

const POS_UNDERLINE: Record<Pos, string> = {
  verb: 'pos-underline-verb',
  noun: 'pos-underline-noun',
  pronoun: 'pos-underline-pronoun',
  adjective: 'pos-underline-adjective',
  adverb: 'pos-underline-adverb',
  preposition: 'pos-underline-preposition',
  conjunction: 'pos-underline-conjunction',
  determiner: 'pos-underline-determiner',
  other: 'pos-underline-other',
};

function regionKey(r: HighlightRegion) {
  return `${r.pos}:${r.start}:${r.end}`;
}

export default function HighlightedText({
  text,
  regions = [],
  showLabels = false,
}: {
  text: string;
  regions?: HighlightRegion[];
  showLabels?: boolean;
}) {
  const parts = useMemo(() => {
    const sorted = [...regions].sort((a, b) => a.start - b.start);
    const out: Array<{ key: string; text: string; region: HighlightRegion | null }> = [];
    let cursor = 0;
    for (const r of sorted) {
      const s = Math.max(cursor, Math.min(r.start, text.length));
      const e = Math.min(r.end, text.length);
      if (s > cursor) out.push({ key: `t-${cursor}`, text: text.slice(cursor, s), region: null });
      if (e > s) out.push({ key: `r-${regionKey(r)}`, text: text.slice(s, e), region: r });
      cursor = Math.max(cursor, e);
    }
    if (cursor < text.length) out.push({ key: `t-${cursor}`, text: text.slice(cursor), region: null });
    return out;
  }, [text, regions]);

  return (
    <p className="question-text text-ink">
      {parts.map((p) =>
        p.region ? (
          <span
            key={p.key}
            className={`group relative font-medium ${POS_TEXT[p.region.pos]} ${POS_UNDERLINE[p.region.pos]}`}
            aria-label={POS_LABEL[p.region.pos]}
            title={POS_LABEL[p.region.pos]}
          >
            {p.text}
            {showLabels && (
              <span className="pointer-events-none absolute -top-5 left-1/2 -translate-x-1/2 whitespace-nowrap rounded-full bg-ink px-1.5 py-0.5 text-[10px] font-semibold text-card opacity-0 transition-opacity duration-200 group-hover:opacity-100 group-focus-within:opacity-100">
                {POS_LABEL[p.region.pos]}
              </span>
            )}
          </span>
        ) : (
          <span key={p.key}>{p.text}</span>
        ),
      )}
    </p>
  );
}