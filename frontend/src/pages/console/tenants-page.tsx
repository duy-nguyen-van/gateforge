import { ChevronRightIcon, PlusIcon, SearchIcon } from 'lucide-react'
import { useEffect, useRef, useState } from 'react'
import { useNavigate } from 'react-router'

import { ConsolePageHeader } from '@/components/layout/console-page-header'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { CreateTenantDialog } from '@/features/admin/create-tenant-dialog'
import { ConsolePagination } from '@/features/admin/console-pagination'
import { ConsoleEmptyState, ConsoleErrorState, ConsoleTableSkeleton } from '@/features/admin/console-state'
import { useConsolePagination } from '@/features/admin/use-console-pagination'
import { useSlashFocus } from '@/features/admin/use-slash-focus'
import { useAdminTenants } from '@/features/admin/use-admin-queries'

export function TenantsPage() {
  const navigate = useNavigate()
  const [createOpen, setCreateOpen] = useState(false)
  const [search, setSearch] = useState('')
  const searchRef = useRef<HTMLInputElement>(null)
  useSlashFocus(searchRef)
  const { page, setPage, pageSize, resetPage, queryParams } = useConsolePagination()

  useEffect(() => {
    resetPage()
  }, [search, resetPage])

  const tenantsQuery = useAdminTenants({ search, ...queryParams })
  const tenants = tenantsQuery.data?.data ?? []
  const meta = tenantsQuery.data?.meta

  return (
    <div>
      <ConsolePageHeader
        title="Tenants"
        description="Organization boundaries for users and identity providers."
        actions={
          <Button type="button" onClick={() => setCreateOpen(true)}>
            <PlusIcon className="h-4 w-4" aria-hidden />
            New tenant
          </Button>
        }
      />

      <div className="relative mb-6 w-full max-w-md">
        <SearchIcon className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-on-surface-variant" aria-hidden />
        <Input
          ref={searchRef}
          type="search"
          aria-label="Search tenants"
          placeholder="Search tenants"
          value={search}
          onChange={(event) => setSearch(event.target.value)}
          className="pl-10"
        />
      </div>

      <div className="overflow-hidden rounded-xl bg-surface-container-lowest ghost-border">
        {tenantsQuery.isLoading ? (
          <ConsoleTableSkeleton columns={6} />
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
                  <td className="px-6 py-4 font-semibold">{t.name || 'Not set'}</td>
                  <td className="px-6 py-4 text-on-surface-variant">{t.domain || 'Not set'}</td>
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
                      <ChevronRightIcon className="h-5 w-5" aria-hidden />
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
