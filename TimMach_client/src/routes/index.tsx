import { Navigate, Outlet, createBrowserRouter } from 'react-router-dom';
import App from '../App';
import { useAuth } from '../modules/auth/hooks/useAuth';
import LoginPage from '../modules/auth/pages/LoginPage';
import RegisterPage from '../modules/auth/pages/RegisterPage';
import DashboardPage from '../modules/dashboard/pages/DashboardPage';
import PatientDetailPage from '../modules/patients/pages/PatientDetailPage';
import PatientFormPage from '../modules/patients/pages/PatientFormPage';
import PatientHistoryPage from '../modules/patients/pages/PatientHistoryPage';
import PatientPredictPage from '../modules/patients/pages/PatientPredictPage';
import PatientsListPage from '../modules/patients/pages/PatientsListPage';
import TemplatesPage from '../modules/exercises/pages/TemplatesPage';
import ProfilePage from '../modules/users/pages/ProfilePage';

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
      { path: 'dashboard', element: <DashboardPage /> },
      { path: 'patients', element: <PatientsListPage /> },
      { path: 'patients/new', element: <PatientFormPage mode="create" /> },
      { path: 'patients/:id', element: <PatientDetailPage /> },
      { path: 'patients/:id/edit', element: <PatientFormPage mode="edit" /> },
      { path: 'patients/:id/predict', element: <PatientPredictPage /> },
      { path: 'patients/:id/history', element: <PatientHistoryPage /> },
      { path: 'profile', element: <ProfilePage /> },
      { path: 'exercises/templates', element: <TemplatesPage /> },
    ],
  },
  {
    element: <AuthLayout />,
    children: [
      { path: '/login', element: <LoginPage /> },
      { path: '/register', element: <RegisterPage /> },
    ],
  },
  { path: '*', element: <Navigate to="/dashboard" replace /> },
]);
