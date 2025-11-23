import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { CreateTemplateRequest } from '../../../types/api';
import { Button } from '../../../components/ui/Button';
import { Input } from '../../../components/ui/Input';

const schema = z.object({
  name: z.string().min(1),
  intensity: z.string().min(1),
  description: z.string().min(1),
  duration_min: z.coerce.number().int().positive(),
  freq_per_week: z.coerce.number().int().positive(),
  target_risk_level: z.string().min(1),
  tags: z.string().optional(),
});

type FormValues = z.infer<typeof schema>;

interface Props {
  onSubmit: (payload: CreateTemplateRequest) => Promise<void> | void;
  submitting?: boolean;
}

export function TemplatesForm({ onSubmit, submitting }: Props) {
  const { register, handleSubmit, formState, reset } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      name: '',
      intensity: 'low',
      description: '',
      duration_min: 30,
      freq_per_week: 3,
      target_risk_level: 'low',
      tags: '',
    },
  });

  const handle = async (values: FormValues) => {
    const payload: CreateTemplateRequest = {
      name: values.name,
      intensity: values.intensity,
      description: values.description,
      duration_min: values.duration_min,
      freq_per_week: values.freq_per_week,
      target_risk_level: values.target_risk_level,
      tags: values.tags ? values.tags.split(',').map((s) => s.trim()).filter(Boolean) : [],
    };
    await onSubmit(payload);
    reset();
  };

  return (
    <form className="space-y-3" onSubmit={handleSubmit(handle)}>
      <div className="grid gap-3 md:grid-cols-2">
        <div className="space-y-1">
          <label className="text-sm font-medium text-slate-700">Tên</label>
          <Input {...register('name')} />
        </div>
        <div className="space-y-1">
          <label className="text-sm font-medium text-slate-700">Cường độ</label>
          <Input {...register('intensity')} placeholder="low/medium/high" />
        </div>
        <div className="space-y-1">
          <label className="text-sm font-medium text-slate-700">Thời lượng (phút)</label>
          <Input type="number" {...register('duration_min')} />
        </div>
        <div className="space-y-1">
          <label className="text-sm font-medium text-slate-700">Số buổi/tuần</label>
          <Input type="number" {...register('freq_per_week')} />
        </div>
        <div className="space-y-1">
          <label className="text-sm font-medium text-slate-700">Mức nguy cơ mục tiêu</label>
          <select
            className="w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm shadow-sm focus:border-blue-500 focus:outline-none"
            {...register('target_risk_level')}
          >
            <option value="low">low</option>
            <option value="medium">medium</option>
            <option value="high">high</option>
            <option value="none">none</option>
          </select>
        </div>
        <div className="space-y-1 md:col-span-2">
          <label className="text-sm font-medium text-slate-700">Mô tả</label>
          <textarea
            className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm shadow-sm focus:border-blue-500 focus:outline-none"
            rows={3}
            {...register('description')}
          />
        </div>
        <div className="space-y-1 md:col-span-2">
          <label className="text-sm font-medium text-slate-700">Tags (cách nhau bởi dấu phẩy)</label>
          <Input placeholder="cardio,stretch" {...register('tags')} />
        </div>
      </div>
      <Button type="submit" disabled={formState.isSubmitting || submitting}>
        {formState.isSubmitting || submitting ? 'Đang lưu...' : 'Thêm template'}
      </Button>
    </form>
  );
}
