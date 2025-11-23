import type { RouteObject } from 'react-router-dom';
import DashboardPage from './pages/DashboardPage';

export function dashboardRoutes(): RouteObject[] {
  return [
    {
      path: 'dashboard',
      element: <DashboardPage />,
    },
  ];
}
