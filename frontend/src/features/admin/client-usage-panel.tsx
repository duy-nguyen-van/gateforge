import { ConsoleErrorState, ConsoleLoadingState } from '@/features/admin/console-state'
import { useAdminClientUsage } from '@/features/admin/use-admin-queries'

export function ClientUsagePanel({ clientId }: { clientId: string }) {
  const usageQuery = useAdminClientUsage(clientId)
  const usage = usageQuery.data?.data

  if (usageQuery.isLoading) {
    return <ConsoleLoadingState label="Loading usage…" />
  }
  if (usageQuery.isError) {
    return <ConsoleErrorState message="Could not load client usage." />
  }
  if (!usage) {
    return null
  }

  return (
    <div className="grid grid-cols-2 gap-3 sm:grid-cols-3">
      {[
        { label: 'Refresh tokens (total)', value: usage.total_refresh_tokens.toLocaleString() },
        { label: 'Active refresh tokens', value: usage.active_refresh_tokens.toLocaleString() },
        { label: 'Authorize events (30d)', value: usage.authorize_events_30d.toLocaleString() },
        { label: 'Token issues (30d)', value: usage.token_issue_events_30d.toLocaleString() },
        {
          label: 'Last token issued',
          value: usage.last_token_issued_at
            ? new Date(usage.last_token_issued_at).toLocaleString()
            : '—',
        },
      ].map((item) => (
        <div key={item.label} className="rounded-lg bg-surface-container-low px-3 py-2">
          <p className="text-[10px] font-bold uppercase tracking-wider text-on-surface-variant">{item.label}</p>
          <p className="mt-1 text-sm font-semibold text-on-surface">{item.value}</p>
        </div>
      ))}
    </div>
  )
}
