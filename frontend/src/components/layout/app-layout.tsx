import { LogOutIcon, ShieldIcon, UserIcon } from 'lucide-react'
import { Link, NavLink, Outlet } from 'react-router-dom'

import { useAuth } from '@/hooks/use-auth'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'

const navItems = [
  { to: '/settings/profile', label: 'Profile', icon: UserIcon },
  { to: '/settings/security', label: 'Security', icon: ShieldIcon },
]

export function AppLayout() {
  const { user, logout } = useAuth()

  return (
    <div className="min-h-screen bg-background">
      <header className="border-b bg-card/80 backdrop-blur">
        <div className="mx-auto flex h-16 max-w-5xl items-center justify-between px-4">
          <Link to="/settings/profile" className="flex items-center gap-2 font-semibold">
            <ShieldIcon className="h-5 w-5 text-primary" />
            IAM Console
          </Link>
          <div className="flex items-center gap-4">
            <span className="hidden text-sm text-muted-foreground sm:inline">{user?.email}</span>
            <Button variant="outline" size="sm" onClick={() => void logout()}>
              <LogOutIcon className="h-4 w-4" />
              Sign out
            </Button>
          </div>
        </div>
      </header>

      <div className="mx-auto grid max-w-5xl gap-6 px-4 py-8 md:grid-cols-[220px_1fr]">
        <nav className="space-y-1">
          {navItems.map(({ to, label, icon: Icon }) => (
            <NavLink
              key={to}
              to={to}
              className={({ isActive }) =>
                cn(
                  'flex items-center gap-2 rounded-lg px-3 py-2 text-sm font-medium transition-colors',
                  isActive
                    ? 'bg-primary text-primary-foreground'
                    : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground',
                )
              }
            >
              <Icon className="h-4 w-4" />
              {label}
            </NavLink>
          ))}
        </nav>
        <main>
          <Outlet />
        </main>
      </div>
    </div>
  )
}
