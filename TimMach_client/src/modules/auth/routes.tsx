import type { RouteObject } from 'react-router-dom';
import LoginPage from './pages/LoginPage';

export function authRoutes(): RouteObject[] {
  return [
    { path: '/login', element: <LoginPage /> },
  ];
}
