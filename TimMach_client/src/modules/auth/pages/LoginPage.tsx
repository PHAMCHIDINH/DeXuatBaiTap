import { LoginForm } from '../components/LoginForm';
import { Card } from '../../../components/ui/Card';

function LoginPage() {
  return (
    <Card className="w-full">
      <div className="mb-4 space-y-1">
        <p className="text-xs uppercase tracking-wide text-slate-500">Welcome back</p>
        <h2 className="text-2xl font-semibold text-slate-900">Đăng nhập SSO</h2>
        <p className="text-sm text-slate-600">Sử dụng Keycloak để truy cập dashboard.</p>
      </div>
      <LoginForm />
    </Card>
  );
}

export default LoginPage;
