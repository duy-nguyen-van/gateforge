import { useState } from 'react'
import { Link, useLocation } from 'react-router-dom'

import { switchTenant } from '@/api/client'
import { setTokens } from '@/auth/token-store'
import { DefaultAvatar } from '@/components/avatars/default-avatar'
import { GateForgeBrand } from '@/components/brand/gateforge-brand'
import { MaterialIcon } from '@/components/icons/material-icon'
import { useAuth } from '@/hooks/use-auth'
import { cn } from '@/lib/utils'

const topNavItems = [
  { label: 'Overview', to: '/console', match: '/console' },
  { label: 'System Logs', to: '/console/audit-logs', match: '/console/audit-logs' },
] as const

export function ConsoleTopbar() {
  const location = useLocation()
  const { user, refreshProfile } = useAuth()
  const [switchingTenant, setSwitchingTenant] = useState(false)
  const homeTo = user?.is_platform_admin ? '/console' : '/settings/profile'

  const tenants = user?.tenants ?? []
  const activeTenant = user?.active_tenant_id

  const onSwitchTenant = async (tenantId: string) => {
    if (!tenantId || tenantId === activeTenant || switchingTenant) return
    setSwitchingTenant(true)
    try {
      const envelope = await switchTenant({ tenant_id: tenantId })
      setTokens(envelope.data, true)
      await refreshProfile()
    } finally {
      setSwitchingTenant(false)
    }
  }

  return (
    <nav className="fixed left-0 right-0 top-0 z-50 flex h-16 w-full items-center justify-between border-b border-slate-200 bg-slate-50 px-6 dark:border-slate-800 dark:bg-slate-900">
      <div className="flex h-full items-center gap-8">
        <GateForgeBrand size="md" layout="horizontal" showTagline={false} linkTo={homeTo} />
        {user?.is_platform_admin ? (
          <div className="hidden h-full items-stretch gap-6 font-manrope text-sm font-medium tracking-tight md:flex">
            {topNavItems.map(({ label, to, match }) => {
              const isActive = match ? location.pathname === match : false
              return (
                <Link
                  key={label}
                  to={to}
                  className={cn(
                    'inline-flex h-full items-center border-b-2 px-1 transition-colors',
                    isActive
                      ? 'border-blue-700 font-semibold text-blue-700 dark:border-blue-400 dark:text-blue-400'
                      : 'border-transparent text-slate-500 hover:text-slate-800 dark:text-slate-400 dark:hover:text-slate-200',
                  )}
                >
                  {label}
                </Link>
              )
            })}
          </div>
        ) : null}
      </div>

      <div className="flex items-center gap-4">
        {tenants.length > 1 ? (
          <div className="relative">
            <MaterialIcon
              name="domain"
              className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-base text-on-surface-variant"
            />
            <select
              className={cn(
                'h-10 min-w-[11rem] max-w-[14rem] appearance-none truncate rounded-xl bg-surface-container-low pl-9 pr-9',
                'text-sm font-medium text-on-surface ghost-border',
                'transition-colors hover:bg-surface-container',
                'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/40',
                'disabled:cursor-wait disabled:opacity-60',
              )}
              value={activeTenant ?? ''}
              disabled={switchingTenant}
              onChange={(e) => void onSwitchTenant(e.target.value)}
              aria-label="Active organization"
            >
              {tenants.map((t) => (
                <option key={t.id} value={t.id}>
                  {t.name || t.domain || t.id.slice(0, 8)}
                </option>
              ))}
            </select>
            <MaterialIcon
              name={switchingTenant ? 'progress_activity' : 'expand_more'}
              className={cn(
                'pointer-events-none absolute right-2.5 top-1/2 -translate-y-1/2 text-base text-on-surface-variant',
                switchingTenant && 'animate-spin',
              )}
            />
          </div>
        ) : null}
        <MaterialIcon name="help" className="cursor-pointer rounded-full p-2 text-slate-500 hover:bg-slate-100" />
        <Link to="/settings/security">
          <MaterialIcon name="settings" className="cursor-pointer rounded-full p-2 text-slate-500 hover:bg-slate-100" />
        </Link>
        <Link to="/settings/profile">
          <DefaultAvatar
            seed={user?.email ?? user?.id ?? 'admin'}
            name={
              user?.first_name ? `${user.first_name} ${user.last_name ?? ''}`.trim() : undefined
            }
            size="sm"
            title={user?.email ?? 'Administrator profile'}
            className="ml-2 ring-2 ring-primary-container"
          />
        </Link>
      </div>
    </nav>
  )
}
