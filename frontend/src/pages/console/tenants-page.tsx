import { useState } from 'react'
import { useNavigate } from 'react-router-dom'

import { MaterialIcon } from '@/components/icons/material-icon'
import { CreateTenantDialog } from '@/features/admin/create-tenant-dialog'
import { ConsolePagination } from '@/features/admin/console-pagination'
import { ConsoleEmptyState, ConsoleErrorState, ConsoleLoadingState } from '@/features/admin/console-state'
import { useConsolePagination } from '@/features/admin/use-console-pagination'
import { useAdminTenants } from '@/features/admin/use-admin-queries'

export function TenantsPage() {
  const navigate = useNavigate()
  const [createOpen, setCreateOpen] = useState(false)
  const { page, setPage, pageSize, queryParams } = useConsolePagination()
  const tenantsQuery = useAdminTenants(queryParams)
  const tenants = tenantsQuery.data?.data ?? []
  const meta = tenantsQuery.data?.meta

  return (
    <div>
      <header className="mb-10 flex flex-col justify-between gap-6 md:flex-row md:items-end">
        <div>
          <h1 className="font-headline text-4xl font-extrabold tracking-tight text-on-surface">Tenant Management</h1>
          <p className="mt-1 text-on-surface-variant">Multi-tenant isolation and federation boundaries.</p>
        </div>
        <button
          type="button"
          onClick={() => setCreateOpen(true)}
          className="flex items-center gap-2 rounded-xl bg-primary px-5 py-2.5 text-sm font-bold text-on-primary shadow-lg shadow-primary/20 transition-opacity hover:opacity-90"
        >
          <MaterialIcon name="add" className="text-sm" />
          New Tenant
        </button>
      </header>

      <div className="overflow-hidden rounded-xl bg-surface-container-lowest ghost-border">
        {tenantsQuery.isLoading ? (
          <ConsoleLoadingState />
        ) : tenantsQuery.isError ? (
          <div className="p-6">
            <ConsoleErrorState message="Could not load tenants." />
          </div>
        ) : tenants.length === 0 ? (
          <ConsoleEmptyState title="No tenants" description="Create a tenant to get started." />
        ) : (
          <table className="w-full text-left text-sm">
            <thead>
              <tr className="bg-surface-container-low">
                {['Tenant ID', 'Organization', 'Domain', 'Users', 'Created', 'Actions'].map((h) => (
                  <th key={h} className="px-6 py-4 text-[10px] font-bold uppercase tracking-widest text-on-surface-variant">
                    {h}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {tenants.map((t) => (
                <tr
                  key={t.id}
                  className="cursor-pointer border-b border-surface-container transition-colors hover:bg-surface-container-low"
                  onClick={() => navigate(`/console/tenants/${t.id}`)}
                >
                  <td className="px-6 py-4 font-mono text-xs font-bold text-primary">{t.id.slice(0, 8)}…</td>
                  <td className="px-6 py-4 font-semibold">{t.name || '—'}</td>
                  <td className="px-6 py-4 text-on-surface-variant">{t.domain || '—'}</td>
                  <td className="px-6 py-4 font-mono text-xs">{t.user_count.toLocaleString()}</td>
                  <td className="px-6 py-4 text-on-surface-variant">{new Date(t.created_at).toLocaleDateString()}</td>
                  <td className="px-6 py-4">
                    <button
                      type="button"
                      title="View tenant"
                      onClick={(e) => {
                        e.stopPropagation()
                        navigate(`/console/tenants/${t.id}`)
                      }}
                      className="text-on-surface-variant hover:text-primary"
                    >
                      <MaterialIcon name="chevron_right" />
                    </button>
                  </td>
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

      <CreateTenantDialog
        open={createOpen}
        onOpenChange={setCreateOpen}
        onCreated={(tenantId) => navigate(`/console/tenants/${tenantId}`)}
      />
    </div>
  )
}
