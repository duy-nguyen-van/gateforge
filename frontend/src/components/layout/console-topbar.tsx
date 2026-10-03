import {
  Building2Icon,
  ChevronDownIcon,
  Loader2Icon,
  MenuIcon,
  MonitorIcon,
  MoonIcon,
  SettingsIcon,
  SunIcon,
} from 'lucide-react'
import { useState } from 'react'
import { Link } from 'react-router'

import { switchTenant } from '@/api/client'
import { setTokens } from '@/auth/token-store'
import { DefaultAvatar } from '@/components/avatars/default-avatar'
import { GateForgeBrand } from '@/components/brand/gateforge-brand'
import { useTheme, type ThemePreference } from '@/components/theme/theme'
import { useAuth } from '@/hooks/use-auth'
import { cn } from '@/lib/utils'

export function ConsoleTopbar({
  mobileOpen,
  onOpenMenu,
}: {
  mobileOpen: boolean
  onOpenMenu: () => void
}) {
  const { user, refreshProfile } = useAuth()
  const { theme, setTheme } = useTheme()
  const [switchingTenant, setSwitchingTenant] = useState(false)
  const homeTo = user?.is_platform_admin ? '/console' : '/settings/profile'

  const tenants = user?.tenants ?? []
  const activeTenant = user?.active_tenant_id
  const profileLabel = user?.email ? `Profile for ${user.email}` : 'Profile'
  const nextTheme: Record<ThemePreference, ThemePreference> = {
    light: 'dark',
    dark: 'system',
    system: 'light',
  }
  const ThemeIcon = theme === 'light' ? SunIcon : theme === 'dark' ? MoonIcon : MonitorIcon
  const themeLabel =
    theme === 'light'
      ? 'Theme: Light. Switch to dark.'
      : theme === 'dark'
        ? 'Theme: Dark. Switch to system.'
        : 'Theme: System. Switch to light.'

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
    <nav className="fixed left-0 right-0 top-0 z-50 flex h-16 w-full items-center justify-between border-b border-outline-variant bg-surface-container-low px-4 lg:px-6">
      <div className="flex h-full items-center gap-3">
        <button
          type="button"
          className="inline-flex rounded-lg p-2 text-on-surface-variant hover:bg-surface-container focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring lg:hidden"
          aria-label="Open menu"
          aria-expanded={mobileOpen}
          aria-controls="console-sidebar"
          onClick={onOpenMenu}
        >
          <MenuIcon className="h-5 w-5" aria-hidden />
        </button>
        <GateForgeBrand size="md" layout="horizontal" showTagline={false} linkTo={homeTo} />
      </div>

      <div className="flex items-center gap-2 sm:gap-3">
        {tenants.length > 1 ? (
          <div className="relative hidden sm:block">
            <Building2Icon
              className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-on-surface-variant"
              aria-hidden
            />
            <select
              className={cn(
                'h-10 min-w-[11rem] max-w-[14rem] appearance-none truncate rounded-lg border border-input bg-background pl-9 pr-9',
                'text-sm font-medium text-on-surface',
                'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring',
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
            {switchingTenant ? (
              <Loader2Icon
                className="pointer-events-none absolute right-2.5 top-1/2 h-4 w-4 -translate-y-1/2 animate-spin text-on-surface-variant"
                aria-hidden
              />
            ) : (
              <ChevronDownIcon
                className="pointer-events-none absolute right-2.5 top-1/2 h-4 w-4 -translate-y-1/2 text-on-surface-variant"
                aria-hidden
              />
            )}
          </div>
        ) : null}
        <button
          type="button"
          aria-label={themeLabel}
          onClick={() => setTheme(nextTheme[theme])}
          className="inline-flex rounded-lg p-2 text-on-surface-variant hover:bg-surface-container focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        >
          <ThemeIcon className="h-5 w-5" aria-hidden />
        </button>
        <Link
          to="/settings/security"
          aria-label="Security settings"
          className="inline-flex rounded-lg p-2 text-on-surface-variant hover:bg-surface-container focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        >
          <SettingsIcon className="h-5 w-5" aria-hidden />
        </Link>
        <Link
          to="/settings/profile"
          aria-label={profileLabel}
          className="rounded-full focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        >
          <DefaultAvatar
            seed={user?.email ?? user?.id ?? 'admin'}
            name={user?.first_name ? `${user.first_name} ${user.last_name ?? ''}`.trim() : undefined}
            size="sm"
            title={user?.email ?? 'Profile'}
            className="ring-2 ring-primary-container"
          />
        </Link>
      </div>
    </nav>
  )
}
