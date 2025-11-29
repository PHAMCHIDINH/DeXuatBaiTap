import { ReactNode } from 'react';
import { cn } from '../../utils/cn';
import { cva, type VariantProps } from 'class-variance-authority';

const badgeVariants = cva(
  'inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2',
  {
    variants: {
      variant: {
        default:
          'border-transparent bg-primary text-primary-foreground hover:bg-primary/80',
        secondary:
          'border-transparent bg-secondary text-secondary-foreground hover:bg-secondary/80',
        destructive:
          'border-transparent bg-destructive text-destructive-foreground hover:bg-destructive/80',
        outline: 'text-foreground',
        success: 'border-transparent bg-green-500 text-white hover:bg-green-600',
        warning: 'border-transparent bg-yellow-500 text-white hover:bg-yellow-600',
      },
    },
    defaultVariants: {
      variant: 'default',
    },
  },
);

interface Props extends VariantProps<typeof badgeVariants> {
  children: ReactNode;
  className?: string;
  color?: 'green' | 'amber' | 'red' | 'gray'; // Keep for backward compatibility
}

export function Badge({ children, className, variant, color }: Props) {
  // Map old colors to new variants if variant is not provided
  let finalVariant = variant;
  if (!finalVariant && color) {
    switch (color) {
      case 'green':
        finalVariant = 'success';
        break;
      case 'amber':
        finalVariant = 'warning';
        break;
      case 'red':
        finalVariant = 'destructive';
        break;
      case 'gray':
        finalVariant = 'secondary';
        break;
    }
  }

  return (
    <div className={cn(badgeVariants({ variant: finalVariant }), className)}>
      {children}
    </div>
  );
}
