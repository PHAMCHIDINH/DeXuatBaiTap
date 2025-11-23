import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { Button } from '../../../components/ui/Button';
import { Input } from '../../../components/ui/Input';

const schema = z.object({
  name: z.string().min(1, 'Tên không được trống'),
  gender: z.coerce.number(),
  dob: z.string().min(4, 'Ngày sinh không hợp lệ'),
});

type FormValues = z.infer<typeof schema>;

interface Props {
  defaultValues?: Partial<FormValues>;
  onSubmit: (values: FormValues) => Promise<void> | void;
  submitting?: boolean;
}

export function PatientForm({ defaultValues, onSubmit, submitting }: Props) {
  const { register, handleSubmit, formState } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      name: '',
      gender: 1,
      dob: '',
      ...defaultValues,
    },
  });

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div className="space-y-1">
        <label className="text-sm font-medium text-slate-700">Tên</label>
        <Input placeholder="Nguyễn Văn A" {...register('name')} />
        {formState.errors.name && (
          <p className="text-xs text-red-600">{formState.errors.name.message}</p>
        )}
      </div>
      <div className="space-y-1">
        <label className="text-sm font-medium text-slate-700">Giới tính</label>
        <select
          className="w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm shadow-sm focus:border-blue-500 focus:outline-none"
          {...register('gender')}
        >
          <option value={1}>Male</option>
          <option value={2}>Female</option>
          <option value={0}>Other</option>
        </select>
      </div>
      <div className="space-y-1">
        <label className="text-sm font-medium text-slate-700">Ngày sinh (YYYY-MM-DD)</label>
        <Input type="date" {...register('dob')} />
        {formState.errors.dob && (
          <p className="text-xs text-red-600">{formState.errors.dob.message}</p>
        )}
      </div>
      <Button type="submit" disabled={submitting || formState.isSubmitting}>
        {submitting || formState.isSubmitting ? 'Đang lưu...' : 'Lưu'}
      </Button>
    </form>
  );
}
