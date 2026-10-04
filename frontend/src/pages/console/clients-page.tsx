import { ChevronRightIcon, MonitorSmartphoneIcon, PlusIcon, SearchIcon } from 'lucide-react'
import { useEffect, useRef, useState } from 'react'
import { useNavigate } from 'react-router'

import { ConsolePageHeader } from '@/components/layout/console-page-header'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { CreateClientDialog } from '@/features/admin/create-client-dialog'
import { ConsolePagination } from '@/features/admin/console-pagination'
import {
  ConsoleEmptyState,
  ConsoleErrorState,
  ConsoleTableSkeleton,
} from '@/features/admin/console-state'
import { useConsolePagination } from '@/features/admin/use-console-pagination'
import { useSlashFocus } from '@/features/admin/use-slash-focus'
import { useAdminClients } from '@/features/admin/use-admin-queries'

export function ClientsPage() {
  const navigate = useNavigate()
  const [createOpen, setCreateOpen] = useState(false)
  const [search, setSearch] = useState('')
  const searchRef = useRef<HTMLInputElement>(null)
  useSlashFocus(searchRef)
  const { page, setPage, pageSize, resetPage, queryParams } = useConsolePagination()

  useEffect(() => {
    resetPage()
  }, [search, resetPage])

  const clientsQuery = useAdminClients({ search, ...queryParams })
  const clients = clientsQuery.data?.data ?? []
  const meta = clientsQuery.data?.meta

  return (
    <div>
      <ConsolePageHeader
        title="Clients"
        description="OAuth 2.0 and OIDC application registrations."
        actions={
          <Button type="button" onClick={() => setCreateOpen(true)}>
            <PlusIcon className="h-4 w-4" aria-hidden />
            Register client
          </Button>
        }
      />

      <div className="relative mb-6 w-full max-w-md">
        <SearchIcon className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-on-surface-variant" aria-hidden />
        <Input
          ref={searchRef}
          type="search"
          aria-label="Search clients"
          placeholder="Search clients"
          value={search}
          onChange={(event) => setSearch(event.target.value)}
          className="pl-10"
        />
      </div>

      {clientsQuery.isLoading ? (
        <ConsoleTableSkeleton columns={3} />
      ) : clientsQuery.isError ? (
        <ConsoleErrorState message="Could not load OAuth clients." />
      ) : clients.length === 0 ? (
        <ConsoleEmptyState
          title="No clients registered"
          description="Register an OAuth client to get started."
        />
      ) : (
        <div className="grid gap-4">
          {clients.map((c) => (
            <button
              key={c.id}
              type="button"
              onClick={() => navigate(`/console/clients/${c.id}`)}
              className="rounded-xl bg-surface-container-lowest p-6 text-left ghost-border transition-colors hover:bg-surface-container-low"
            >
              <div className="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
                <div className="flex items-start gap-4">
                  <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-primary-container text-primary">
                    <MonitorSmartphoneIcon className="h-6 w-6" aria-hidden />
                  </div>
                  <div>
                    <h3 className="font-headline text-lg font-bold">{c.name || c.client_id}</h3>
                    <p className="font-mono text-xs text-on-surface-variant">{c.client_id}</p>
                    <div className="mt-2 flex flex-wrap gap-2">
                      <span className="rounded-md bg-tertiary-container px-2 py-0.5 text-[10px] font-bold text-on-tertiary-container">
                        {c.is_public ? 'Public' : 'Confidential'}
                      </span>
                      {c.grant_types?.map((grant) => (
                        <span
                          key={grant}
                          className="rounded-md bg-secondary-container px-2 py-0.5 text-[10px] font-bold text-on-secondary-container"
                        >
                          {grant}
                        </span>
                      ))}
                    </div>
                  </div>
                </div>
                <ChevronRightIcon className="h-5 w-5 text-on-surface-variant" aria-hidden />
              </div>
            </button>
          ))}
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

      <CreateClientDialog
        open={createOpen}
        onOpenChange={setCreateOpen}
        onCreated={(id) => navigate(`/console/clients/${id}`)}
      />
    </div>
  )
}
