import { NavLink, useNavigate } from 'react-router-dom';
import { cn } from '../../utils/cn';
import { useAuth } from '../../modules/auth/useAuth';
import { Button } from '../ui/Button';
import { LayoutDashboard, Users, User, Activity, LogOut } from 'lucide-react';

const links = [
  { to: '/dashboard', label: 'Dashboard', icon: LayoutDashboard },
  { to: '/patients', label: 'Patients', icon: Users },
  { to: '/profile', label: 'Profile', icon: User },
  { to: '/exercises/templates', label: 'Exercises', icon: Activity },
];

export function Sidebar() {
  const navigate = useNavigate();
  const { logout } = useAuth();

  return (
    <aside className="hidden w-64 flex-col border-r bg-card px-4 py-6 md:flex">
      <div className="mb-8 flex items-center px-2">
        <button
          className="flex items-center gap-2 text-xl font-bold text-primary"
          onClick={() => navigate('/dashboard')}
        >
          <Activity className="h-6 w-6" />
          <span>TimMach</span>
        </button>
      </div>
      <nav className="flex flex-1 flex-col gap-1">
        {links.map((link) => (
          <NavLink
            key={link.to}
            to={link.to}
            className={({ isActive }) =>
              cn(
                'flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground',
                isActive
                  ? 'bg-primary/10 text-primary hover:bg-primary/20'
                  : 'text-muted-foreground',
              )
            }
          >
            <link.icon className="h-4 w-4" />
            {link.label}
          </NavLink>
        ))}
      </nav>
      <div className="mt-auto border-t pt-4">
        <Button
          variant="ghost"
          className="w-full justify-start gap-3 text-muted-foreground hover:text-destructive"
          onClick={logout}
        >
          <LogOut className="h-4 w-4" />
          Logout
        </Button>
      </div>
    </aside>
  );
}
