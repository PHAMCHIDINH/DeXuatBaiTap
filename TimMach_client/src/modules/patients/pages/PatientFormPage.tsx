import { useNavigate, useParams } from 'react-router-dom';
import { Button } from '../../../components/ui/Button';
import { Card } from '../../../components/ui/Card';
import { PatientForm } from '../components/PatientForm';
import { useCreatePatient, usePatient, useUpdatePatient } from '../hooks/usePatients';

interface Props {
  mode: 'create' | 'edit';
}

function PatientFormPage({ mode }: Props) {
  const { id } = useParams<{ id: string }>();
  const patientId = id ?? '';
  const navigate = useNavigate();
  const createMutation = useCreatePatient();
  const updateMutation = useUpdatePatient(patientId);
  const { data: patient, isLoading } = usePatient(mode === 'edit' ? id : undefined);

  const isEdit = mode === 'edit';

  const handleSubmit = async (values: { name: string; gender: number; dob: string }) => {
    if (isEdit && patientId) {
      await updateMutation.mutateAsync(values);
      navigate(`/patients/${patientId}`);
      return;
    }
    const created = await createMutation.mutateAsync(values);
    navigate(`/patients/${created.id}`);
  };

  if (isEdit && !id) {
    return <p className="text-sm text-red-600">Thiếu patient id.</p>;
  }

  if (isEdit && isLoading) {
    return <p className="text-sm text-slate-600">Đang tải...</p>;
  }

  return (
    <Card className="max-w-2xl">
      <div className="mb-4 flex items-center justify-between">
        <div>
          <p className="text-xs uppercase tracking-wide text-slate-500">{isEdit ? 'Chỉnh sửa' : 'Tạo mới'}</p>
          <h2 className="text-2xl font-semibold text-slate-900">Patient</h2>
        </div>
        <Button variant="secondary" onClick={() => navigate(-1)}>
          Quay lại
        </Button>
      </div>
      <PatientForm
        defaultValues={
          isEdit && patient ? { name: patient.name, gender: patient.gender, dob: patient.dob } : undefined
        }
        onSubmit={handleSubmit}
        submitting={createMutation.isPending || updateMutation.isPending}
      />
    </Card>
  );
}

export default PatientFormPage;
