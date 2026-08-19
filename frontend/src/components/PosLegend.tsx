import type { Pos } from '../types';
import { POS_LABEL, POS_ID } from '../lib/format';

const POS_STYLE: Record<Pos, string> = {
  verb: 'bg-pos-verb',
  noun: 'bg-pos-noun',
  pronoun: 'bg-pos-pronoun',
  adjective: 'bg-pos-adjective',
  adverb: 'bg-pos-adverb',
  preposition: 'bg-pos-preposition',
  conjunction: 'bg-pos-conjunction',
  determiner: 'bg-pos-determiner',
  other: 'bg-pos-other',
};

const ORDER: Pos[] = [
  'verb',
  'noun',
  'pronoun',
  'adjective',
  'adverb',
  'preposition',
  'conjunction',
  'determiner',
  'other',
];

export default function PosLegend({
  present,
  onToggle,
  detailed = false,
}: {
  present: Set<Pos>;
  onToggle?: (pos: Pos) => void;
  detailed?: boolean;
}) {
  if (detailed) {
    return (
      <div className="grid grid-cols-1 gap-1 sm:grid-cols-2" aria-label="Part of speech legend">
        {ORDER.map((pos) => {
          const active = present.has(pos);
          const row = (
            <span className="flex items-start gap-2 rounded-sm px-2 py-1.5">
              <span
                className={`pos-chip shrink-0 ${active ? `${POS_STYLE[pos]} text-on-pos` : 'bg-card-muted text-ink-faint'}`}
              >
                {POS_LABEL[pos]}
              </span>
              <span className={`text-[13px] leading-snug ${active ? 'text-ink' : 'text-ink-faint'}`}>
                {POS_ID[pos]}
              </span>
            </span>
          );
          if (onToggle) {
            return (
              <button
                key={pos}
                type="button"
                onClick={() => onToggle(pos)}
                aria-pressed={active}
                className="w-full rounded-sm text-left transition-colors duration-200 hover:bg-card-muted"
              >
                {row}
              </button>
            );
          }
          return <div key={pos}>{row}</div>;
        })}
      </div>
    );
  }

  return (
    <div className="flex flex-wrap items-center gap-1.5" aria-label="Part of speech legend">
      {ORDER.map((pos) => {
        const active = present.has(pos);
        const label = POS_LABEL[pos];
        if (onToggle) {
          return (
            <button
              key={pos}
              type="button"
              onClick={() => onToggle(pos)}
              aria-pressed={active}
              title={POS_ID[pos]}
              className={`pos-chip transition-opacity duration-200 ${
                active ? `${POS_STYLE[pos]} text-on-pos` : 'bg-card-muted text-ink-faint'
              }`}
            >
              {label}
            </button>
          );
        }
        return (
          <span key={pos} title={POS_ID[pos]} className={`pos-chip ${active ? `${POS_STYLE[pos]} text-on-pos` : 'bg-card-muted text-ink-faint'}`}>
            {label}
          </span>
        );
      })}
    </div>
  );
}