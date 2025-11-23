import { TemplatesForm } from '../components/TemplatesForm';
import { TemplatesTable } from '../components/TemplatesTable';
import { useCreateTemplate, useTemplates } from '../hooks/useExercises';
import { Card } from '../../../components/ui/Card';

function TemplatesPage() {
  const { data, isLoading } = useTemplates();
  const createMutation = useCreateTemplate();
  const templates = data?.templates ?? [];

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
