import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { GraduationCap } from 'lucide-react';
import Button from '../../components/Button';
import Input from '../../components/Input';
import Card from '../../components/Card';
import { useAuth } from '../../hooks/useAuth';
import { ApiError } from '../../lib/api';

export default function LoginPage() {
  const { login } = useAuth();
  const navigate = useNavigate();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [errors, setErrors] = useState<{ username?: string; password?: string; form?: string }>({});
  const [submitting, setSubmitting] = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setErrors({});
    const next: typeof errors = {};
    if (!username.trim()) next.username = 'Username is required';
    if (!password) next.password = 'Password is required';
    if (Object.keys(next).length) {
      setErrors(next);
      return;
    }
    setSubmitting(true);
    try {
      await login(username.trim(), password);
      navigate('/', { replace: true });
    } catch (err) {
      if (err instanceof ApiError) {
        if (err.status === 401) {
          setErrors({ form: 'Incorrect username or password.' });
        } else if (err.status === 429) {
          setErrors({ form: 'Too many attempts. Please wait a minute and try again.' });
        } else {
          setErrors({ form: err.message });
        }
      } else {
        setErrors({ form: 'Unable to reach the server. Try again.' });
      }
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-paper px-4">
      <Card className="w-full max-w-sm p-8">
        <div className="mb-8 flex flex-col items-center gap-2">
          <GraduationCap className="h-10 w-10 text-primary" aria-hidden="true" />
          <h1 className="text-xl font-bold tracking-tight text-ink">TOEFL Prep</h1>
          <p className="text-sm text-ink-muted">Your personal practice room.</p>
        </div>
        <form onSubmit={handleSubmit} className="flex flex-col gap-4" noValidate>
          <Input
            label="Username"
            type="text"
            autoComplete="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            error={errors.username}
          />
          <Input
            label="Password"
            type="password"
            autoComplete="current-password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            error={errors.password}
          />
          {errors.form && (
            <p className="rounded-sm border border-danger-soft bg-danger-soft/50 px-3 py-2 text-sm text-danger" role="alert">
              {errors.form}
            </p>
          )}
          <Button type="submit" loading={submitting} className="mt-1">
            Sign in
          </Button>
          <p className="text-center text-[13px] text-ink-faint">Seeded accounts only</p>
        </form>
      </Card>
    </div>
  );
}