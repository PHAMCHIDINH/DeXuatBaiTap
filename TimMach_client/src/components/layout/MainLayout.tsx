import { ReactNode } from 'react';
import { Sidebar } from './Sidebar';
import { Header } from './Header';
import { cn } from '../../utils/cn';

interface Props {
  children: ReactNode;
}

export function MainLayout({ children }: Props) {
  return (
    <div className={cn('flex min-h-screen bg-muted/40')}>
      <Sidebar />
      <div className="flex min-h-screen flex-1 flex-col">
        <Header />
        <main className="flex-1 p-6 md:p-8 pt-6">{children}</main>
      </div>
    </div>
  );
}
