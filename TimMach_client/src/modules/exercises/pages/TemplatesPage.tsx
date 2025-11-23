import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Card } from '../../../components/ui/Card';
import { createTemplate, listTemplates } from '../api';
import { TemplatesForm } from '../components/TemplatesForm';
import { TemplatesTable } from '../components/TemplatesTable';
import { TemplateResponse, CreateTemplateRequest } from '../types';

function TemplatesPage() {
  const qc = useQueryClient();
  const { data: templates = [], isLoading } = useQuery<TemplateResponse[]>({
    queryKey: ['exercise-templates'],
    queryFn: listTemplates,
  });
  const createMutation = useMutation({
    mutationFn: (payload: CreateTemplateRequest) => createTemplate(payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['exercise-templates'] });
    },
  });

  return (
    <div className="space-y-4">
      <div>
        <p className="text-xs uppercase tracking-wide text-slate-500">Bài tập</p>
        <h2 className="text-2xl font-semibold text-slate-900">Exercise templates</h2>
      </div>

      <Card title="Thêm template">
        <TemplatesForm
          onSubmit={async (payload) => {
            await createMutation.mutateAsync(payload);
          }}
          submitting={createMutation.isPending}
        />
      </Card>

      <div>
        {isLoading ? (
          <p className="text-sm text-slate-600">Đang tải...</p>
        ) : (
          <TemplatesTable templates={templates} />
        )}
      </div>
    </div>
  );
}

export default TemplatesPage;
