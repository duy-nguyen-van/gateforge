import { useEffect, useState } from 'react'

import {
  auditResultBadgeClass,
  formatAuditAction,
  formatAuditTimestamp,
} from '@/features/admin/admin-utils'
import { ConsolePagination } from '@/features/admin/console-pagination'
import {
  ConsoleEmptyState,
  ConsoleErrorState,
  ConsoleLoadingState,
} from '@/features/admin/console-state'
import { useConsolePagination } from '@/features/admin/use-console-pagination'
import { useAdminLoginHistory } from '@/features/admin/use-admin-queries'

const RESULT_OPTIONS = ['', 'success', 'failure', 'denied'] as const

export function LoginHistoryPage() {
  const [resultFilter, setResultFilter] = useState('')
  const [tenantFilter, setTenantFilter] = useState('')
  const [actorFilter, setActorFilter] = useState('')
  const { page, setPage, pageSize, resetPage, queryParams } = useConsolePagination()

  useEffect(() => {
    resetPage()
  }, [actorFilter, resultFilter, tenantFilter, resetPage])

  const historyQuery = useAdminLoginHistory({
    result: resultFilter || undefined,
    tenant_id: tenantFilter || undefined,
    actor_id: actorFilter || undefined,
    ...queryParams,
  })

  const logs = historyQuery.data?.data ?? []
  const meta = historyQuery.data?.meta

  return (
    <div>
      <header className="mb-10 flex flex-col justify-between gap-6 md:flex-row md:items-end">
        <div>
          <h1 className="font-headline text-4xl font-extrabold tracking-tight text-on-surface">Login History</h1>
          <p className="mt-1 text-on-surface-variant">
            Password, passkey, federation, and OIDC sign-in events across all tenants.
          </p>
        </div>
      </header>

      <div className="mb-6 grid grid-cols-1 gap-3 md:grid-cols-3">
        <input
          type="text"
          placeholder="Filter by actor user ID"
          value={actorFilter}
          onChange={(e) => setActorFilter(e.target.value)}
          className="rounded-xl border-none bg-surface-container-low px-4 py-2.5 text-sm focus:ring-1 focus:ring-primary"
        />
        <select
          value={resultFilter}
          onChange={(e) => setResultFilter(e.target.value)}
          className="rounded-xl border-none bg-surface-container-low px-4 py-2.5 text-sm focus:ring-1 focus:ring-primary"
        >
          {RESULT_OPTIONS.map((opt) => (
            <option key={opt || 'all'} value={opt}>
              {opt ? opt.charAt(0).toUpperCase() + opt.slice(1) : 'All results'}
            </option>
          ))}
        </select>
        <input
          type="text"
          placeholder="Filter by tenant ID"
          value={tenantFilter}
          onChange={(e) => setTenantFilter(e.target.value)}
          className="rounded-xl border-none bg-surface-container-low px-4 py-2.5 text-sm focus:ring-1 focus:ring-primary"
        />
      </div>

      <div className="overflow-hidden rounded-xl bg-surface-container-lowest ghost-border">
        {historyQuery.isLoading ? (
          <ConsoleLoadingState label="Loading login history…" />
        ) : historyQuery.isError ? (
          <ConsoleErrorState message="Could not load login history." />
        ) : logs.length === 0 ? (
          <ConsoleEmptyState title="No login events" description="Sign-in attempts will appear here once users authenticate." />
        ) : (
          <table className="w-full text-left text-sm">
            <thead>
              <tr className="bg-surface-container-low">
                {['Time', 'Action', 'Result', 'Actor', 'IP', 'Tenant'].map((h) => (
                  <th key={h} className="px-6 py-3 text-[10px] font-bold uppercase tracking-widest text-on-surface-variant">
                    {h}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {logs.map((log) => (
                <tr key={log.id} className="border-b border-surface-container hover:bg-surface-container-low">
                  <td className="px-6 py-3 font-mono text-xs text-on-surface-variant">
                    {formatAuditTimestamp(log.created_at)}
                  </td>
                  <td className="px-6 py-3">{formatAuditAction(log.action)}</td>
                  <td className="px-6 py-3">
                    <span className={`rounded-full px-2 py-0.5 text-xs font-bold ${auditResultBadgeClass(log.result)}`}>
                      {log.result}
                    </span>
                  </td>
                  <td className="px-6 py-3 font-mono text-xs">{log.actor_id?.slice(0, 8) ?? '—'}…</td>
                  <td className="px-6 py-3 font-mono text-xs">{log.ip_address ?? '—'}</td>
                  <td className="px-6 py-3 font-mono text-xs">{log.tenant_id?.slice(0, 8) ?? '—'}…</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
        {meta?.total != null && meta.total > 0 ? (
          <ConsolePagination
            page={meta.page ?? page}
            pageSize={meta.page_size ?? pageSize}
            total={meta.total}
            onPageChange={setPage}
          />
        ) : null}
      </div>
    </div>
  )
}
