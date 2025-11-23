import type { RouteObject } from 'react-router-dom';
import PatientDetailPage from './pages/PatientDetailPage';
import PatientFormPage from './pages/PatientFormPage';
import PatientHistoryPage from './pages/PatientHistoryPage';
import PatientPredictPage from './pages/PatientPredictPage';
import PatientsListPage from './pages/PatientsListPage';

export function patientRoutes(): RouteObject[] {
  return [
    { path: 'patients', element: <PatientsListPage /> },
    { path: 'patients/new', element: <PatientFormPage mode="create" /> },
    { path: 'patients/:id', element: <PatientDetailPage /> },
    { path: 'patients/:id/edit', element: <PatientFormPage mode="edit" /> },
    { path: 'patients/:id/predict', element: <PatientPredictPage /> },
    { path: 'patients/:id/history', element: <PatientHistoryPage /> },
  ];
}
