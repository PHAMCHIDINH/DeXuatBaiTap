import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { Button } from '../../../components/ui/Button';
import { Input } from '../../../components/ui/Input';

const schema = z.object({
  age_years: z.coerce.number().positive(),
  gender: z.coerce.number(),
  height: z.coerce.number().positive(),
  weight: z.coerce.number().positive(),
  ap_hi: z.coerce.number(),
  ap_lo: z.coerce.number(),
  cholesterol: z.coerce.number(),
  gluc: z.coerce.number(),
  smoke: z.coerce.number(),
  alco: z.coerce.number(),
  active: z.coerce.number(),
});

type FormValues = z.infer<typeof schema>;

interface Props {
  defaultValues?: Partial<FormValues>;
  onSubmit: (values: FormValues) => Promise<void> | void;
}

export function PredictForm({ defaultValues, onSubmit }: Props) {
  // defaultValues lấy từ lần dự đoán gần nhất (nếu có) để tiết kiệm thao tác nhập lại.
  const { register, handleSubmit, formState } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      age_years: 30,
      gender: 1,
      height: 170,
      weight: 65,
      ap_hi: 120,
      ap_lo: 80,
      cholesterol: 1,
      gluc: 1,
      smoke: 0,
      alco: 0,
      active: 1,
      ...defaultValues,
    },
  });

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="grid grid-cols-1 gap-4 md:grid-cols-2">
      {(
        [
          ['age_years', 'Tuổi (năm)'],
          ['gender', 'Giới tính (0/1/2)'],
          ['height', 'Chiều cao (cm)'],
          ['weight', 'Cân nặng (kg)'],
          ['ap_hi', 'Huyết áp tâm thu'],
          ['ap_lo', 'Huyết áp tâm trương'],
          ['cholesterol', 'Cholesterol (1-3)'],
          ['gluc', 'Glucose (1-3)'],
          ['smoke', 'Smoke (0/1)'],
          ['alco', 'Alco (0/1)'],
          ['active', 'Active (0/1)'],
        ] as const
      ).map(([key, label]) => (
        <div key={key} className="space-y-1">
          <label className="text-sm font-medium text-slate-700">{label}</label>
          <Input type="number" step="any" {...register(key as keyof FormValues)} />
          {formState.errors[key as keyof FormValues] && (
            <p className="text-xs text-red-600">
              {formState.errors[key as keyof FormValues]?.message as string}
            </p>
          )}
        </div>
      ))}
      <div className="md:col-span-2">
        <Button type="submit" disabled={formState.isSubmitting}>
          {formState.isSubmitting ? 'Đang gửi...' : 'Predict now'}
        </Button>
      </div>
    </form>
  );
}
