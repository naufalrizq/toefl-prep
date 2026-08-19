import type { ReactNode } from 'react';

type Tone = 'neutral' | 'success' | 'danger' | 'warning' | 'primary';

const tones: Record<Tone, string> = {
  neutral: 'bg-card-muted text-ink-muted',
  success: 'bg-success-soft text-success',
  danger: 'bg-danger-soft text-danger',
  warning: 'bg-card-muted text-warning',
  primary: 'bg-primary-soft text-primary',
};

export default function Badge({
  tone = 'neutral',
  children,
  className = '',
}: {
  tone?: Tone;
  children: ReactNode;
  className?: string;
}) {
  return (
    <span
      className={`inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-xs font-semibold ${tones[tone]} ${className}`}
    >
      {children}
    </span>
  );
}