import type { HTMLAttributes } from 'react';

interface CardProps extends HTMLAttributes<HTMLDivElement> {
  interactive?: boolean;
}

export default function Card({ interactive = false, className = '', children, ...rest }: CardProps) {
  return (
    <div
      className={`rounded-md bg-card shadow-sm transition-all duration-200 ${
        interactive ? 'cursor-pointer hover:-translate-y-px hover:shadow-md' : ''
      } ${className}`}
      {...rest}
    >
      {children}
    </div>
  );
}