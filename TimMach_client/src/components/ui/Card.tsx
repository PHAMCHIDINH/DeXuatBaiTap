import { ReactNode } from 'react';
import { cn } from '../../utils/cn';

interface Props {
  title?: string;
  action?: ReactNode;
  children: ReactNode;
  className?: string;
}

export function Card({ title, action, className, children }: Props) {
  return (
    <div
      className={cn(
        'rounded-lg border bg-card text-card-foreground shadow-sm',
        className,
      )}
    >
      {(title || action) && (
        <div className="flex flex-col space-y-1.5 p-6 pb-4">
          <div className="flex items-center justify-between">
            {title && (
              <h3 className="text-2xl font-semibold leading-none tracking-tight">
                {title}
              </h3>
            )}
            {action}
          </div>
        </div>
      )}
      <div className={cn('p-6 pt-0', !title && !action && 'p-6')}>{children}</div>
    </div>
  );
}
