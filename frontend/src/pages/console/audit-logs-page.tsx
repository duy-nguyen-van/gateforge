import { useEffect, useState } from 'react'

import { MaterialIcon } from '@/components/icons/material-icon'
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
import { useAdminAuditLogs } from '@/features/admin/use-admin-queries'

const RESULT_OPTIONS = ['', 'success', 'failure', 'denied'] as const

export function AuditLogsPage() {
  const [actionFilter, setActionFilter] = useState('')
  const [resultFilter, setResultFilter] = useState('')
  const [tenantFilter, setTenantFilter] = useState('')
  const { page, setPage, pageSize, resetPage, queryParams } = useConsolePagination()

  useEffect(() => {
    resetPage()
  }, [actionFilter, resultFilter, tenantFilter, resetPage])

  const auditQuery = useAdminAuditLogs({
    action: actionFilter || undefined,
    result: resultFilter || undefined,
    tenant_id: tenantFilter || undefined,
    ...queryParams,
  })

  const logs = auditQuery.data?.data ?? []
  const meta = auditQuery.data?.meta

  return (
    <div>
      <header className="mb-10 flex flex-col justify-between gap-6 md:flex-row md:items-end">
        <div>
          <h1 className="font-headline text-4xl font-extrabold tracking-tight text-on-surface">Audit Logs</h1>
          <p className="mt-1 text-on-surface-variant">Immutable event trail for compliance and forensics.</p>
        </div>
        <div className="flex gap-3">
          <button type="button" disabled className="flex items-center gap-2 rounded-xl bg-surface-container-highest/60 px-4 py-2 text-sm font-bold">
            <MaterialIcon name="download" className="text-lg" />
            Export
          </button>
        </div>
      </header>

      <div className="mb-6 grid grid-cols-1 gap-3 md:grid-cols-3">
        <input
          type="text"
          placeholder="Filter by action (e.g. auth.login)"
          value={actionFilter}
          onChange={(e) => setActionFilter(e.target.value)}
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
        {auditQuery.isLoading ? (
          <ConsoleLoadingState label="Loading audit logs…" />
        ) : auditQuery.isError ? (
          <ConsoleErrorState message="Could not load audit logs." />
        ) : logs.length === 0 ? (
          <ConsoleEmptyState
            title="No audit events yet"
            description="Security and admin actions will appear here as they occur."
          />
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full min-w-[960px] text-left text-sm">
              <thead className="border-b border-surface-container bg-surface-container-low/50 text-xs font-bold uppercase tracking-wider text-on-surface-variant">
                <tr>
                  <th className="px-6 py-4">Time</th>
                  <th className="px-6 py-4">Action</th>
                  <th className="px-6 py-4">Result</th>
                  <th className="px-6 py-4">Actor</th>
                  <th className="px-6 py-4">Resource</th>
                  <th className="px-6 py-4">IP</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-surface-container">
                {logs.map((log) => (
                  <tr key={log.id} className="hover:bg-surface-container-low/40">
                    <td className="whitespace-nowrap px-6 py-4 text-on-surface-variant">
                      {formatAuditTimestamp(log.created_at)}
                    </td>
                    <td className="px-6 py-4 font-medium text-on-surface">{formatAuditAction(log.action)}</td>
                    <td className="px-6 py-4">
                      <span className={`rounded-full px-2.5 py-1 text-xs font-bold uppercase ${auditResultBadgeClass(log.result)}`}>
                        {log.result}
                      </span>
                    </td>
                    <td className="px-6 py-4">
                      <div className="text-on-surface">{log.actor_type}</div>
                      {log.actor_id ? (
                        <div className="mt-0.5 max-w-[180px] truncate font-mono text-xs text-on-surface-variant">{log.actor_id}</div>
                      ) : null}
                    </td>
                    <td className="px-6 py-4">
                      {log.resource_type ? (
                        <>
                          <div className="text-on-surface">{log.resource_type}</div>
                          {(log.resource_name || log.resource_id) && (
                            <div className="mt-0.5 max-w-[180px] truncate text-xs text-on-surface-variant">
                              {log.resource_name ?? log.resource_id}
                            </div>
                          )}
                        </>
                      ) : (
                        <span className="text-on-surface-variant">—</span>
                      )}
                    </td>
                    <td className="whitespace-nowrap px-6 py-4 font-mono text-xs text-on-surface-variant">
                      {log.ip_address ?? '—'}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
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
