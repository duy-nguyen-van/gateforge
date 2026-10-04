import type { LucideIcon } from 'lucide-react'
import { useEffect, useRef, useState } from 'react'
import { LogOutIcon } from 'lucide-react'
import { NavLink } from 'react-router'

import { accountNavItems, consoleNavItems } from '@/components/layout/console-nav'
import { useAuth } from '@/hooks/use-auth'
import { cn } from '@/lib/utils'

function SidebarNavLink({
  to,
  label,
  icon: Icon,
  end,
  onNavigate,
}: {
  to: string
  label: string
  icon: LucideIcon
  end?: boolean
  onNavigate?: () => void
}) {
  return (
    <NavLink
      to={to}
      end={end}
      onClick={onNavigate}
      className={({ isActive }) =>
        cn(
          'flex items-center gap-3 rounded-lg px-3 py-2 text-sm transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring',
          isActive
            ? 'bg-surface-container-lowest font-semibold text-primary'
            : 'text-on-surface-variant hover:bg-surface-container hover:text-on-surface',
        )
      }
    >
      <Icon className="h-5 w-5 shrink-0" aria-hidden />
      <span>{label}</span>
    </NavLink>
  )
}

export function ConsoleSidebar({
  mobileOpen,
  onClose,
}: {
  mobileOpen: boolean
  onClose: () => void
}) {
  const { user, logout } = useAuth()
  const isAdmin = user?.is_platform_admin
  const asideRef = useRef<HTMLElement>(null)
  const [isDesktop, setIsDesktop] = useState(
    () => typeof window !== 'undefined' && window.matchMedia('(min-width: 1024px)').matches,
  )

  useEffect(() => {
    const media = window.matchMedia('(min-width: 1024px)')
    const onChange = () => setIsDesktop(media.matches)
    media.addEventListener('change', onChange)
    return () => media.removeEventListener('change', onChange)
  }, [])

  useEffect(() => {
    if (!mobileOpen) {
      return
    }
    const aside = asideRef.current
    const previous = document.activeElement instanceof HTMLElement ? document.activeElement : null
    aside?.querySelector<HTMLElement>('a, button')?.focus()

    function onKeyDown(event: KeyboardEvent) {
      if (event.key === 'Escape') {
        onClose()
      }
    }

    window.addEventListener('keydown', onKeyDown)
    return () => {
      window.removeEventListener('keydown', onKeyDown)
      previous?.focus()
    }
  }, [mobileOpen, onClose])

  return (
    <aside
      ref={asideRef}
      id="console-sidebar"
      aria-label="Console"
      role={mobileOpen ? 'dialog' : undefined}
      aria-modal={mobileOpen ? true : undefined}
      inert={!isDesktop && !mobileOpen ? true : undefined}
      className={cn(
        'fixed left-0 top-0 z-40 flex min-h-[100dvh] w-64 flex-col border-r border-outline-variant bg-surface-container-low p-4 pt-20',
        'transition-transform lg:translate-x-0',
        mobileOpen ? 'translate-x-0' : '-translate-x-full',
      )}
    >
      <div className="mb-8 px-2">
        <h2 className="font-headline text-lg font-semibold text-on-surface">
          {isAdmin ? 'Admin console' : 'Account'}
        </h2>
        <p className="text-xs font-medium text-on-surface-variant">GateForge Identity</p>
      </div>

      <nav className="flex-1 space-y-6 overflow-y-auto">
        <div className="space-y-1">
          <p className="px-3 text-xs font-medium text-on-surface-variant">Account</p>
          {accountNavItems.map(({ to, label, icon, ...rest }) => (
            <SidebarNavLink
              key={to}
              to={to}
              label={label}
              icon={icon}
              end={'end' in rest ? rest.end : false}
              onNavigate={onClose}
            />
          ))}
        </div>

        {isAdmin ? (
          <div className="space-y-1">
            <p className="px-3 text-xs font-medium text-on-surface-variant">Administration</p>
            {consoleNavItems.map(({ to, label, icon, ...rest }) => (
              <SidebarNavLink
                key={to}
                to={to}
                label={label}
                icon={icon}
                end={'end' in rest ? rest.end : false}
                onNavigate={onClose}
              />
            ))}
          </div>
        ) : null}
      </nav>

      <div className="mt-auto border-t border-outline-variant pt-4">
        <button
          type="button"
          onClick={() => void logout()}
          className="flex w-full items-center gap-3 rounded-lg px-3 py-2 text-sm text-on-surface-variant transition-colors hover:bg-surface-container hover:text-error focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        >
          <LogOutIcon className="h-5 w-5 shrink-0" aria-hidden />
          <span>Sign out</span>
        </button>
      </div>
    </aside>
  )
}
