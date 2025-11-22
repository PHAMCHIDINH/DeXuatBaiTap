import { Outlet } from 'react-router-dom';
import { MainLayout } from './components/layout/MainLayout';

export function App() {
  return (
    <MainLayout>
      <Outlet />
    </MainLayout>
  );
}

export default App;
