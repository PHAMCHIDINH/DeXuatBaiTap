import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../useAuth';
import { Button } from '../../../components/ui/Button';
import { Input } from '../../../components/ui/Input';

const schema = z.object({
  username: z.string().min(1, 'Không được bỏ trống'),
  password: z.string().min(1, 'Không được bỏ trống'),
});

type FormValues = z.infer<typeof schema>;

export function LoginForm() {
  const { login } = useAuth();
  const navigate = useNavigate();
  const [error, setError] = useState<string | null>(null);
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { username: '', password: '' },
  });

  const onSubmit = async (values: FormValues) => {
    setError(null);
    try {
      await login(values);
      navigate('/dashboard');
    } catch (err) {
      console.error(err);
      const message =
        err instanceof Error ? err.message : 'Đăng nhập thất bại. Vui lòng kiểm tra thông tin.';
      setError(message);
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div className="space-y-1">
        <label className="text-sm font-medium text-slate-700">Email / Username</label>
        <Input type="text" placeholder="you@example.com" {...register('username')} />
        {errors.username && <p className="text-xs text-red-600">{errors.username.message}</p>}
      </div>
      <div className="space-y-1">
        <label className="text-sm font-medium text-slate-700">Password</label>
        <Input type="password" placeholder="••••••" {...register('password')} />
        {errors.password && <p className="text-xs text-red-600">{errors.password.message}</p>}
      </div>
      {error && <p className="text-sm text-red-600">{error}</p>}
      <Button type="submit" className="w-full" disabled={isSubmitting}>
        {isSubmitting ? 'Đang đăng nhập...' : 'Đăng nhập'}
      </Button>
    </form>
  );
}
