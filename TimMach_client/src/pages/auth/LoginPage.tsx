import { Link } from 'react-router-dom';
import { LoginForm } from '../../components/auth/LoginForm';
import { Card } from '../../components/ui/Card';

function LoginPage() {
  return (
    <Card className="w-full">
      <div className="mb-4 space-y-1">
        <p className="text-xs uppercase tracking-wide text-slate-500">Welcome back</p>
        <h2 className="text-2xl font-semibold text-slate-900">Đăng nhập</h2>
      </div>
      <LoginForm />
      <p className="mt-4 text-sm text-slate-600">
        Chưa có tài khoản?{' '}
        <Link to="/register" className="text-blue-600 hover:underline">
          Đăng ký
        </Link>
      </p>
    </Card>
  );
}

export default LoginPage;
