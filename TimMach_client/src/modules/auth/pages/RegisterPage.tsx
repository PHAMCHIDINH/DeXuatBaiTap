import { Link } from 'react-router-dom';
import { RegisterForm } from '../components/RegisterForm';
import { Card } from '../../../components/ui/Card';

function RegisterPage() {
  return (
    <Card className="w-full">
      <div className="mb-4 space-y-1">
        <p className="text-xs uppercase tracking-wide text-slate-500">Tạo tài khoản</p>
        <h2 className="text-2xl font-semibold text-slate-900">Đăng ký</h2>
      </div>
      <RegisterForm />
      <p className="mt-4 text-sm text-slate-600">
        Đã có tài khoản?{' '}
        <Link to="/login" className="text-blue-600 hover:underline">
          Đăng nhập
        </Link>
      </p>
    </Card>
  );
}

export default RegisterPage;
