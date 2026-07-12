import { NavLink } from 'react-router-dom'

import { MaterialIcon } from '@/components/icons/material-icon'
import { accountNavItems, consoleNavItems } from '@/components/layout/console-nav'
import { useAuth } from '@/hooks/use-auth'
import { cn } from '@/lib/utils'

function SidebarNavLink({
  to,
  label,
  icon,
  end,
}: {
  to: string
  label: string
  icon: string
  end?: boolean
}) {
  return (
    <NavLink
      to={to}
      end={end}
      className={({ isActive }) =>
        cn(
          'flex items-center gap-3 rounded-md px-3 py-2 text-sm transition-all duration-200 ease-in-out',
          isActive
            ? 'bg-white font-bold text-blue-700 shadow-sm dark:bg-slate-900 dark:text-blue-300'
            : 'text-slate-600 hover:bg-slate-200 hover:text-blue-600 dark:text-slate-400 dark:hover:bg-slate-800 dark:hover:text-blue-300',
        )
      }
    >
      <MaterialIcon name={icon} className="text-xl" />
      <span>{label}</span>
    </NavLink>
  )
}

export function ConsoleSidebar() {
  const { user, logout } = useAuth()
  const isAdmin = user?.is_platform_admin

  return (
    <aside className="fixed left-0 top-0 z-40 flex h-screen w-64 flex-col border-r border-slate-200 bg-slate-100 p-4 pt-20 dark:border-slate-800 dark:bg-slate-950">
      <div className="mb-8 px-2">
        <h2 className="font-headline text-lg font-black text-slate-900 dark:text-slate-50">
          {isAdmin ? 'System Architect' : 'My Account'}
        </h2>
        <p className="text-xs font-medium uppercase tracking-widest text-slate-500">
          {isAdmin ? 'Identity & Access' : 'GateForge Identity'}
        </p>
      </div>

      <nav className="flex-1 space-y-6 overflow-y-auto">
        <div className="space-y-1">
          <p className="px-3 text-[10px] font-bold uppercase tracking-widest text-slate-400">Account</p>
          {accountNavItems.map(({ to, label, icon, ...rest }) => (
            <SidebarNavLink key={to} to={to} label={label} icon={icon} end={'end' in rest ? rest.end : false} />
          ))}
        </div>

        {isAdmin ? (
          <div className="space-y-1">
            <p className="px-3 text-[10px] font-bold uppercase tracking-widest text-slate-400">Administration</p>
            {consoleNavItems.map(({ to, label, icon, ...rest }) => (
              <SidebarNavLink key={to} to={to} label={label} icon={icon} end={'end' in rest ? rest.end : false} />
            ))}
          </div>
        ) : null}
      </nav>

      <div className="mt-auto space-y-1 border-t border-slate-200 pt-4 dark:border-slate-800">
        <a
          href="#"
          className="flex items-center gap-3 rounded-md px-3 py-2 text-sm text-slate-600 transition-all duration-200 hover:bg-slate-200 hover:text-blue-600 dark:text-slate-400 dark:hover:bg-slate-800"
        >
          <MaterialIcon name="contact_support" className="text-xl" />
          <span>Support</span>
        </a>
        <button
          type="button"
          onClick={() => void logout()}
          className="flex w-full items-center gap-3 rounded-md px-3 py-2 text-sm text-slate-600 transition-all duration-200 hover:text-error dark:text-slate-400"
        >
          <MaterialIcon name="logout" className="text-xl" />
          <span>Sign Out</span>
        </button>
      </div>
    </aside>
  )
}
