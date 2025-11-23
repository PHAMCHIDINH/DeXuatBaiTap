import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../modules/auth/hooks/useAuth';
import { Button } from '../ui/Button';
import { cn } from '../../utils/cn';

export function Header() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  return (
    <header className={cn('sticky top-0 z-10 flex items-center justify-between gap-4 border-b border-slate-200 bg-white/90 px-6 py-3 backdrop-blur')}> 
      <div>
        <p className="text-xs uppercase tracking-wide text-slate-500">TimMach Dashboard</p>
        <h1 className="text-lg font-semibold text-slate-900">Heart Risk Prediction</h1>
      </div>
      <div className="flex items-center gap-3 text-sm">
        {user && <span className="text-slate-700">{user.email}</span>}
        <Button variant="secondary" size="sm" onClick={() => navigate('/profile')}>
          Profile
        </Button>
        <Button variant="ghost" size="sm" onClick={logout}>
          Logout
        </Button>
      </div>
    </header>
  );
}
