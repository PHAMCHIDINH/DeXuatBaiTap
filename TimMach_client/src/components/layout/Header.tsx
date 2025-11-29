import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../modules/auth/useAuth';
import { Button } from '../ui/Button';
import { cn } from '../../utils/cn';
import { User, LogOut } from 'lucide-react';

export function Header() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  return (
    <header
      className={cn(
        'sticky top-0 z-10 flex h-16 items-center justify-between border-b bg-background/95 px-6 backdrop-blur supports-[backdrop-filter]:bg-background/60',
      )}
    >
      <div>
        <p className="text-xs font-medium uppercase tracking-wider text-muted-foreground">
          TimMach Dashboard
        </p>
        <h1 className="text-lg font-semibold tracking-tight text-foreground">
          Heart Risk Prediction
        </h1>
      </div>
      <div className="flex items-center gap-4">
        {user && (
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <span className="hidden sm:inline-block">Welcome,</span>
            <span className="font-medium text-foreground">{user.email}</span>
          </div>
        )}
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => navigate('/profile')}
            className="gap-2"
          >
            <User className="h-4 w-4" />
            <span className="hidden sm:inline-block">Profile</span>
          </Button>
          <Button
            variant="ghost"
            size="sm"
            onClick={logout}
            className="text-muted-foreground hover:text-destructive"
          >
            <LogOut className="h-4 w-4" />
            <span className="sr-only">Logout</span>
          </Button>
        </div>
      </div>
    </header>
  );
}
