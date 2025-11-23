import type { RouteObject } from 'react-router-dom';
import TemplatesPage from './pages/TemplatesPage';

export function exerciseRoutes(): RouteObject[] {
  return [{ path: 'exercises/templates', element: <TemplatesPage /> }];
}
