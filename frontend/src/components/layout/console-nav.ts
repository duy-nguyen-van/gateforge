export const consoleNavItems = [
  { to: '/console', label: 'Dashboard', icon: 'dashboard', end: true },
  { to: '/console/users', label: 'Users', icon: 'group' },
  { to: '/console/clients', label: 'Clients', icon: 'devices' },
  { to: '/console/tenants', label: 'Tenants', icon: 'domain' },
  { to: '/console/identity-providers', label: 'Identity Providers', icon: 'fingerprint' },
  { to: '/console/audit-logs', label: 'Audit Logs', icon: 'history' },
  { to: '/console/login-history', label: 'Login History', icon: 'login' },
] as const

export const accountNavItems = [
  { to: '/settings/profile', label: 'Profile', icon: 'person', end: true },
  { to: '/settings/security', label: 'Security', icon: 'shield' },
] as const
