import { useEffect, useState } from 'react'

import { ConsolePageHeader } from '@/components/layout/console-page-header'
import { Input } from '@/components/ui/input'
import {
  auditResultBadgeClass,
  formatAuditAction,
  formatAuditTimestamp,
} from '@/features/admin/admin-utils'
import { ConsolePagination } from '@/features/admin/console-pagination'
import {
  ConsoleEmptyState,
  ConsoleErrorState,
  ConsoleTableSkeleton,
} from '@/features/admin/console-state'
import { OrganizationFilter } from '@/features/admin/organization-filter'
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
      <ConsolePageHeader
        title="Audit logs"
        description="Event trail for security review and forensics."
      />

      <div className="mb-6 grid grid-cols-1 gap-3 md:grid-cols-3">
        <Input
          type="text"
          aria-label="Filter by action"
          placeholder="Filter by action (for example auth.login)"
          value={actionFilter}
          onChange={(e) => setActionFilter(e.target.value)}
        />
        <select
          aria-label="Filter by result"
          value={resultFilter}
          onChange={(e) => setResultFilter(e.target.value)}
          className="flex h-10 w-full rounded-lg border border-input bg-background px-3 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        >
          {RESULT_OPTIONS.map((opt) => (
            <option key={opt || 'all'} value={opt}>
              {opt ? opt.charAt(0).toUpperCase() + opt.slice(1) : 'All results'}
            </option>
          ))}
        </select>
        <OrganizationFilter value={tenantFilter} onChange={setTenantFilter} />
      </div>

      <div className="overflow-hidden rounded-xl bg-surface-container-lowest ghost-border">
        {auditQuery.isLoading ? (
          <ConsoleTableSkeleton columns={6} />
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
                        <span className="text-on-surface-variant">Not set</span>
                      )}
                    </td>
                    <td className="whitespace-nowrap px-6 py-4 font-mono text-xs text-on-surface-variant">
                      {log.ip_address ?? 'Not set'}
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
