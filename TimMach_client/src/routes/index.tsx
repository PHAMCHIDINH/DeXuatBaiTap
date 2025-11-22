import { Navigate, Outlet, createBrowserRouter } from 'react-router-dom';
import App from '../App';
import { useAuth } from '../hooks/useAuth';
import LoginPage from '../pages/auth/LoginPage';
import RegisterPage from '../pages/auth/RegisterPage';
import DashboardPage from '../pages/dashboard/DashboardPage';
import PatientsListPage from '../pages/patients/PatientsListPage';
import PatientDetailPage from '../pages/patients/PatientDetailPage';
import PatientFormPage from '../pages/patients/PatientFormPage';
import PatientPredictPage from '../pages/patients/PatientPredictPage';
import PatientHistoryPage from '../pages/patients/PatientHistoryPage';
import ProfilePage from '../pages/profile/ProfilePage';
import TemplatesPage from '../pages/exercises/TemplatesPage';

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
