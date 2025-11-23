import type { RouteObject } from 'react-router-dom';
import { Navigate, Outlet, createBrowserRouter } from 'react-router-dom';
import App from '../App';
import { useAuth } from '../modules/auth/useAuth';
import { authRoutes } from '../modules/auth/routes';
import { dashboardRoutes } from '../modules/dashboard/routes';
import { exerciseRoutes } from '../modules/exercises/routes';
import { patientRoutes } from '../modules/patients/routes';
import { userRoutes } from '../modules/users/routes';

function ProtectedShell() {
  const { token, loading } = useAuth();
  if (loading) return <div className="p-8 text-sm text-slate-600">Đang tải...</div>;
  if (!token) return <Navigate to="/login" replace />;
  return <App />;
}

function AuthLayout() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-slate-50">
      <div className="mx-auto flex min-h-screen max-w-md flex-col justify-center px-4 py-10">
        <Outlet />
      </div>
    </div>
  );
}

export const router = createBrowserRouter([
  {
    path: '/',
    element: <ProtectedShell />,
    children: [
      { index: true, element: <Navigate to="/dashboard" replace /> },
      ...dashboardRoutes(),
      ...patientRoutes(),
      ...userRoutes(),
      ...exerciseRoutes(),
    ] as RouteObject[],
  },
  {
    element: <AuthLayout />,
    children: authRoutes(),
  },
  { path: '*', element: <Navigate to="/dashboard" replace /> },
]);
