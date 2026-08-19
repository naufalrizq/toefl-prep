import { useEffect, useMemo, useState } from 'react';
import { formatClock } from '../lib/format';

interface TimerProps {
  deadlineMs: number | null;
  label?: string;
  warnAfterSec?: number;
  dangerAfterSec?: number;
  onExpire?: () => void;
}

export default function Timer({
  deadlineMs,
  label = 'Time left',
  warnAfterSec = 60,
  dangerAfterSec = 10,
  onExpire,
}: TimerProps) {
  const [now, setNow] = useState(() => Date.now());

  useEffect(() => {
    if (deadlineMs == null) return;
    const id = window.setInterval(() => setNow(Date.now()), 500);
    return () => window.clearInterval(id);
  }, [deadlineMs]);

  const secondsLeft = useMemo(() => {
    if (deadlineMs == null) return null;
    return Math.max(0, Math.ceil((deadlineMs - now) / 1000));
  }, [deadlineMs, now]);

  useEffect(() => {
    if (secondsLeft === 0 && deadlineMs != null) onExpire?.();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [secondsLeft === 0]);

  if (secondsLeft == null) return null;

  const urgent = secondsLeft <= dangerAfterSec;
  const warning = secondsLeft <= warnAfterSec && !urgent;

  const tone = urgent
    ? 'text-danger'
    : warning
      ? 'text-warning animate-pulse'
      : 'text-ink';

  return (
    <div className="flex flex-col items-end" role="timer" aria-label={`${label}: ${formatClock(secondsLeft)}`}>
      <span className="text-xs uppercase tracking-wide text-ink-muted">{label}</span>
      <span
        className={`tabular text-xl font-semibold ${tone} transition-colors duration-200`}
        style={{ minWidth: '3.4em', textAlign: 'right' }}
        aria-live="off"
      >
        {formatClock(secondsLeft)}
      </span>
    </div>
  );
}