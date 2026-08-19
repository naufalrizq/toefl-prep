import type { ReactNode } from 'react';
import Card from './Card';
import { TrendingUp, TrendingDown } from 'lucide-react';

export default function StatCard({
  label,
  value,
  delta,
  deltaPositiveIsGood = true,
}: {
  label: string;
  value: ReactNode;
  delta?: number;
  deltaPositiveIsGood?: boolean;
}) {
  const showDelta = delta !== undefined && !Number.isNaN(delta);
  const deltaIsGood = deltaPositiveIsGood ? (delta ?? 0) >= 0 : (delta ?? 0) <= 0;
  return (
    <Card className="p-6">
      <p className="text-[13px] font-medium uppercase tracking-wide text-ink-muted">{label}</p>
      <p className="mt-1 text-4xl font-extrabold tabular text-ink">{value}</p>
      {showDelta && (
        <p className={`mt-1 flex items-center gap-1 text-[13px] font-medium ${deltaIsGood ? 'text-success' : 'text-danger'}`}>
          {delta! >= 0 ? <TrendingUp className="h-4 w-4" aria-hidden="true" /> : <TrendingDown className="h-4 w-4" aria-hidden="true" />}
          {delta! >= 0 ? '+' : ''}
          {delta} pts
        </p>
      )}
    </Card>
  );
}