import { lazy, Suspense, type ReactNode } from 'react';
import { Navigate, Outlet, Route, Routes } from 'react-router-dom';
import AppShell from './components/AppShell';
import Spinner from './components/Spinner';
import { useAuth } from './hooks/useAuth';

const LoginPage = lazy(() => import('./features/auth/LoginPage'));
const DashboardPage = lazy(() => import('./features/dashboard/DashboardPage'));
const ExamListPage = lazy(() => import('./features/exams/ExamListPage'));
const HistoryPage = lazy(() => import('./features/history/HistoryPage'));
const ReviewPage = lazy(() => import('./features/review/ReviewPage'));
const AdminQuestionsPage = lazy(() => import('./features/admin/AdminQuestionsPage'));
const AdminExamsPage = lazy(() => import('./features/admin/AdminExamsPage'));
const QuizPage = lazy(() => import('./features/quiz/QuizPage'));

function Lazy({ children }: { children: ReactNode }) {
  return <Suspense fallback={<div className="flex justify-center py-24"><Spinner className="h-8 w-8" /></div>}>{children}</Suspense>;
}

function RequireAuth() {
  const { user, isLoading } = useAuth();
  if (isLoading) {
    return <div className="flex justify-center py-24"><Spinner className="h-8 w-8" /></div>;
  }
  if (!user) return <Navigate to="/login" replace />;
  return <Outlet />;
}

function RequireAdmin() {
  const { user, isAdmin, isLoading } = useAuth();
  if (isLoading) {
    return <div className="flex justify-center py-24"><Spinner className="h-8 w-8" /></div>;
  }
  if (!user) return <Navigate to="/login" replace />;
  if (!isAdmin) return <Navigate to="/" replace />;
  return <Outlet />;
}

export default function AppRoutes() {
  return (
    <Routes>
      <Route path="/login" element={<Lazy><LoginPage /></Lazy>} />

      <Route element={<RequireAuth />}>
        <Route element={<AppShell />}>
          <Route index element={<Lazy><DashboardPage /></Lazy>} />
          <Route path="exams" element={<Lazy><ExamListPage /></Lazy>} />
          <Route path="history" element={<Lazy><HistoryPage /></Lazy>} />
          <Route path="quiz/:id" element={<Lazy><QuizPage /></Lazy>} />
          <Route path="attempts/:id/review" element={<Lazy><ReviewPage /></Lazy>} />
        </Route>
      </Route>

      <Route element={<RequireAdmin />}>
        <Route element={<AppShell />}>
          <Route path="admin/questions" element={<Lazy><AdminQuestionsPage /></Lazy>} />
          <Route path="admin/exams" element={<Lazy><AdminExamsPage /></Lazy>} />
        </Route>
      </Route>

      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}