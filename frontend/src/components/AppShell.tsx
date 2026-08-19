import { NavLink, Outlet } from 'react-router-dom';
import { GraduationCap, LogOut, LayoutDashboard, FilePlus2, Library, ClipboardList, Moon, Sun } from 'lucide-react';
import { useAuth } from '../hooks/useAuth';
import { useTheme } from '../hooks/useTheme';

const studentLinks = [
  { to: '/', label: 'Dashboard', icon: LayoutDashboard, end: true },
  { to: '/exams', label: 'Start exam', icon: FilePlus2 },
  { to: '/history', label: 'History', icon: ClipboardList },
];

const adminLinks = [
  { to: '/admin/questions', label: 'Questions', icon: Library, end: false },
  { to: '/admin/exams', label: 'Exams', icon: ClipboardList, end: false },
];

function linkClass({ isActive }: { isActive: boolean }) {
  return `inline-flex items-center gap-2 rounded-sm px-3 py-2 text-sm font-medium transition-colors duration-200 ${
    isActive ? 'bg-primary-soft text-primary' : 'text-ink-muted hover:bg-card-muted hover:text-ink'
  }`;
}

export default function AppShell() {
  const { user, isAdmin, logout } = useAuth();
  const { theme, toggle } = useTheme();
  const links = [...studentLinks, ...(isAdmin ? adminLinks : [])];

  return (
    <div className="min-h-screen bg-paper">
      <header className="print-hidden sticky top-0 z-20 border-b border-border bg-paper/95 backdrop-blur-sm">
        <div className="mx-auto flex h-14 max-w-app items-center justify-between gap-4 px-4 sm:px-8">
          <NavLink to="/" className="flex items-center gap-2 text-ink">
            <GraduationCap className="h-6 w-6 text-primary" aria-hidden="true" />
            <span className="text-base font-bold tracking-tight">TOEFL Prep</span>
          </NavLink>
          <nav className="flex items-center gap-1 overflow-x-auto" aria-label="Primary">
            {links.map((l) => (
              <NavLink key={l.to} to={l.to} end={l.end} className={linkClass}>
                <l.icon className="h-4 w-4" aria-hidden="true" />
                {l.label}
              </NavLink>
            ))}
          </nav>
          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={toggle}
              className="inline-flex items-center rounded-sm p-2 text-ink-muted transition-colors duration-200 hover:bg-card-muted hover:text-ink"
              aria-label={theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
            >
              {theme === 'dark' ? <Sun className="h-4 w-4" aria-hidden="true" /> : <Moon className="h-4 w-4" aria-hidden="true" />}
            </button>
            <span className="hidden text-sm text-ink-muted sm:inline">{user?.email}</span>
            <button
              type="button"
              onClick={() => void logout()}
              className="inline-flex items-center gap-1.5 rounded-sm px-2.5 py-2 text-sm text-ink-muted transition-colors duration-200 hover:bg-card-muted hover:text-danger"
              aria-label="Log out"
            >
              <LogOut className="h-4 w-4" aria-hidden="true" />
              <span className="hidden md:inline">Log out</span>
            </button>
          </div>
        </div>
      </header>
      <main className="mx-auto max-w-app px-4 py-8 sm:px-8">
        <Outlet />
      </main>
    </div>
  );
}