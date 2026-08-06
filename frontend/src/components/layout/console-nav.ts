import type { LucideIcon } from 'lucide-react'
import {
  Building2Icon,
  FingerprintIcon,
  HistoryIcon,
  LayoutDashboardIcon,
  LogInIcon,
  MonitorSmartphoneIcon,
  ShieldIcon,
  UserIcon,
  UsersIcon,
} from 'lucide-react'

export type ConsoleNavItem = {
  to: string
  label: string
  icon: LucideIcon
  end?: boolean
}

export const consoleNavItems: ConsoleNavItem[] = [
  { to: '/console', label: 'Dashboard', icon: LayoutDashboardIcon, end: true },
  { to: '/console/users', label: 'Users', icon: UsersIcon },
  { to: '/console/clients', label: 'Clients', icon: MonitorSmartphoneIcon },
  { to: '/console/tenants', label: 'Tenants', icon: Building2Icon },
  { to: '/console/identity-providers', label: 'Identity Providers', icon: FingerprintIcon },
  { to: '/console/audit-logs', label: 'Audit Logs', icon: HistoryIcon },
  { to: '/console/login-history', label: 'Login History', icon: LogInIcon },
]

export const accountNavItems: ConsoleNavItem[] = [
  { to: '/settings/profile', label: 'Profile', icon: UserIcon, end: true },
  { to: '/settings/security', label: 'Security', icon: ShieldIcon },
]
