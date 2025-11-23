import type { RouteObject } from 'react-router-dom';
import ProfilePage from './pages/ProfilePage';

export function userRoutes(): RouteObject[] {
  return [{ path: 'profile', element: <ProfilePage /> }];
}
