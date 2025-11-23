import { NavLink, useNavigate } from 'react-router-dom';
import { cn } from '../../utils/cn';
import { useAuth } from '../../modules/auth/hooks/useAuth';

const links = [
  { to: '/dashboard', label: 'Dashboard' },
  { to: '/patients', label: 'Patients' },
  { to: '/profile', label: 'Profile' },
  { to: '/exercises/templates', label: 'Exercises' },
];

export function Sidebar() {
  const navigate = useNavigate();
  const { logout } = useAuth();

  return (
    <aside className="hidden w-60 flex-col border-r border-slate-200 bg-white px-4 py-6 md:flex">
      <button
        className="mb-8 text-left text-xl font-bold text-blue-600"
        onClick={() => navigate('/dashboard')}
      >
        TimMach
      </button>
      <nav className="flex flex-1 flex-col gap-2 text-sm">
        {links.map((link) => (
          <NavLink
            key={link.to}
            to={link.to}
            className={({ isActive }) =>
              cn(
                'rounded-lg px-3 py-2 font-medium text-slate-700 hover:bg-slate-100',
                isActive && 'bg-blue-50 text-blue-700 border border-blue-100',
              )
            }
          >
            {link.label}
          </NavLink>
        ))}
        <button
          onClick={logout}
          className="mt-auto rounded-lg px-3 py-2 text-left font-medium text-slate-600 hover:bg-slate-100"
        >
          Logout
        </button>
      </nav>
    </aside>
  );
}
