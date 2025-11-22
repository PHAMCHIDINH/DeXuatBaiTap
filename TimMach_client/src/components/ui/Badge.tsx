import { ReactNode } from 'react';
import { cn } from '../../utils/cn';

interface Props {
  children: ReactNode;
  color?: 'green' | 'amber' | 'red' | 'gray';
}

export function Badge({ children, color = 'gray' }: Props) {
  const palette: Record<typeof color, string> = {
    green: 'bg-green-100 text-green-700',
    amber: 'bg-amber-100 text-amber-700',
    red: 'bg-red-100 text-red-700',
    gray: 'bg-slate-100 text-slate-700',
  };
  return (
    <span className={cn('inline-flex items-center rounded-full px-3 py-1 text-xs font-medium', palette[color])}>
      {children}
    </span>
  );
}
