import { Link } from 'react-router-dom'

import { MfaAvatarPreviewStack } from '@/components/avatars/default-avatar'
import { MaterialIcon } from '@/components/icons/material-icon'
import {
  ConsoleEmptyState,
  ConsoleErrorState,
  ConsoleLoadingState,
} from '@/features/admin/console-state'
import { useAdminStats, useTenantIdentityProviders } from '@/features/admin/use-admin-queries'

const defaultTenantId = import.meta.env.VITE_DEFAULT_TENANT_ID ?? '00000000-0000-0000-0000-000000000001'

export function DashboardPage() {
  const statsQuery = useAdminStats()
  const providersQuery = useTenantIdentityProviders(defaultTenantId)
  const stats = statsQuery.data?.data
  const idpRows = providersQuery.data?.data ?? []

  return (
    <div>
      <header className="mb-10 flex items-end justify-between">
        <div>
          <h1 className="font-headline text-4xl font-extrabold tracking-tight text-on-surface">Systems Overview</h1>
          <p className="mt-1 font-body text-on-surface-variant">Enterprise identity management and global access state.</p>
        </div>
        <div className="flex gap-3">
          <div className="flex items-center gap-2 rounded-xl bg-surface-container px-4 py-2">
            <span className="h-2 w-2 animate-pulse rounded-circle bg-green-500" />
            <span className="font-label text-sm font-semibold text-on-surface">Global Status: Optimal</span>
          </div>
        </div>
      </header>

      <div className="grid grid-cols-12 gap-6">
        <div className="col-span-12 grid grid-cols-2 gap-6 lg:col-span-8">
          {statsQuery.isLoading ? (
            <div className="col-span-2">
              <ConsoleLoadingState label="Loading dashboard stats…" />
            </div>
          ) : statsQuery.isError ? (
            <div className="col-span-2">
              <ConsoleErrorState message="Could not load dashboard statistics." />
            </div>
          ) : stats ? (
            <>
              <div className="flex min-h-[160px] flex-col justify-between rounded-full bg-surface-container-lowest p-6 shadow-sm ghost-border">
                <div className="flex items-start justify-between">
                  <MaterialIcon name="person_check" className="rounded-lg bg-primary-container p-2 text-primary-dim" />
                  <span className="rounded-full bg-green-50 px-2 py-0.5 text-xs font-bold text-green-600">
                    {stats.total_users.toLocaleString()} users
                  </span>
                </div>
                <div>
                  <p className="mb-1 font-label text-sm font-bold uppercase tracking-wider text-on-surface-variant">
                    Active Sessions
                  </p>
                  <h3 className="font-headline text-5xl font-extrabold text-on-primary-fixed">
                    {stats.active_sessions.toLocaleString()}
                  </h3>
                </div>
              </div>

              <div className="flex min-h-[160px] flex-col justify-between rounded-full bg-surface-container-lowest p-6 shadow-sm ghost-border">
                <div className="flex items-start justify-between">
                  <MaterialIcon name="security" className="rounded-lg bg-primary-container p-2 text-primary-dim" />
                  <MfaAvatarPreviewStack />
                </div>
                <div>
                  <p className="mb-1 font-label text-sm font-bold uppercase tracking-wider text-on-surface-variant">
                    MFA Adoption Rate
                  </p>
                  <div className="flex items-baseline gap-2">
                    <h3 className="font-headline text-5xl font-extrabold text-on-primary-fixed">
                      {stats.mfa_enabled_percent}%
                    </h3>
                    <span className="font-body text-sm text-on-surface-variant">
                      {stats.mfa_enabled_count} enrolled
                    </span>
                  </div>
                </div>
              </div>
            </>
          ) : null}

          <div className="col-span-2 rounded-full bg-surface-container-lowest p-8 shadow-sm ghost-border">
            <div className="mb-4">
              <h4 className="font-headline text-xl font-bold">Authentication Trends</h4>
              <p className="text-sm text-on-surface-variant">Session frequency vs. risk score</p>
            </div>
            <ConsoleEmptyState
              title="Trend data not available yet"
              description="Authentication trend metrics will appear here when the backend exposes this data."
            />
          </div>
        </div>

        <div className="col-span-12 space-y-6 lg:col-span-4">
          <div className="relative h-full overflow-hidden rounded-full bg-surface-container-high p-6 shadow-sm">
            <div className="mb-6 flex items-center justify-between">
              <h4 className="font-headline text-xl font-extrabold tracking-tight">Security Alerts</h4>
            </div>
            <ConsoleEmptyState
              title="No alerts"
              description="Security alert data is not available from the API yet."
            />
            <div className="mt-6">
              <Link
                to="/console/audit-logs"
                className="block w-full rounded-xl bg-surface-container-highest py-4 text-center text-xs font-black uppercase tracking-widest text-on-surface-variant transition-all hover:bg-primary-container hover:text-on-primary-container"
              >
                View audit logs
              </Link>
            </div>
          </div>
        </div>

        <div className="col-span-12">
          <div className="overflow-hidden rounded-full bg-surface-container-lowest shadow-sm ghost-border">
            <div className="flex items-center justify-between border-b border-surface-container px-8 py-6">
              <h4 className="font-headline text-lg font-bold">Identity Provider Health</h4>
              <Link to="/console/identity-providers" className="text-xs font-bold text-primary">
                Manage providers
              </Link>
            </div>
            <div className="overflow-x-auto">
              {providersQuery.isLoading ? (
                <div className="px-8 py-8">
                  <ConsoleLoadingState label="Loading providers…" />
                </div>
              ) : providersQuery.isError ? (
                <div className="px-8 py-8">
                  <ConsoleErrorState message="Could not load identity providers." />
                </div>
              ) : idpRows.length === 0 ? (
                <div className="px-8 py-8">
                  <ConsoleEmptyState
                    title="No identity providers configured"
                    description="Enable Google federation from the identity providers page."
                  />
                </div>
              ) : (
                <table className="w-full border-collapse text-left">
                  <thead>
                    <tr className="bg-surface-container-low">
                      {['Provider', 'Protocol', 'Status', 'Actions'].map((h) => (
                        <th
                          key={h}
                          className={`px-8 py-4 text-[10px] font-bold uppercase tracking-widest text-on-surface-variant ${h === 'Actions' ? 'text-right' : ''}`}
                        >
                          {h}
                        </th>
                      ))}
                    </tr>
                  </thead>
                  <tbody className="text-sm">
                    {idpRows.map((p) => (
                      <tr key={p.provider} className="group transition-colors hover:bg-surface-container-low">
                        <td className="px-8 py-5">
                          <div className="flex items-center gap-3">
                            <div className="flex h-8 w-8 items-center justify-center rounded bg-surface-container-highest text-xs font-bold text-primary">
                              {p.name.slice(0, 2).toUpperCase()}
                            </div>
                            <span className="font-semibold">{p.name}</span>
                          </div>
                        </td>
                        <td className="px-8 py-5">
                          <span className="rounded-md bg-tertiary-container px-2 py-1 text-[10px] font-bold text-on-tertiary-container">
                            OIDC
                          </span>
                        </td>
                        <td className="px-8 py-5">
                          <span
                            className={`rounded-full px-2 py-0.5 text-[10px] font-bold uppercase ${
                              p.enabled ? 'bg-green-50 text-green-700' : 'bg-amber-50 text-amber-700'
                            }`}
                          >
                            {p.enabled ? 'Enabled' : 'Disabled'}
                          </span>
                        </td>
                        <td className="px-8 py-5 text-right">
                          <Link to="/console/identity-providers" className="text-xs font-bold text-primary">
                            Manage
                          </Link>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
