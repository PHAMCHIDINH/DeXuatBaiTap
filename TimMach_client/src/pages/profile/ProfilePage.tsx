import { useAuth } from '../../hooks/useAuth';
import { Card } from '../../components/ui/Card';
import { Button } from '../../components/ui/Button';
import { formatDate } from '../../utils/format';

function ProfilePage() {
  const { user, refreshProfile } = useAuth();

  return (
    <Card className="max-w-xl" title="Hồ sơ người dùng">
      {user ? (
        <div className="space-y-2 text-sm">
          <p className="text-slate-600">Email</p>
          <p className="font-medium text-slate-900">{user.email}</p>
          <p className="text-slate-600">Tạo lúc</p>
          <p className="font-medium text-slate-900">{formatDate(user.created_at)}</p>
          <Button variant="secondary" size="sm" onClick={refreshProfile}>
            Làm mới
          </Button>
        </div>
      ) : (
        <p className="text-sm text-slate-600">Không có thông tin user.</p>
      )}
    </Card>
  );
}

export default ProfilePage;
