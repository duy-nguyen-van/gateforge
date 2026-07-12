import { useState } from 'react'
import { useNavigate } from 'react-router-dom'

import { MaterialIcon } from '@/components/icons/material-icon'
import { CreateClientDialog } from '@/features/admin/create-client-dialog'
import { ConsolePagination } from '@/features/admin/console-pagination'
import {
  ConsoleEmptyState,
  ConsoleErrorState,
  ConsoleLoadingState,
} from '@/features/admin/console-state'
import { useConsolePagination } from '@/features/admin/use-console-pagination'
import { useAdminClients } from '@/features/admin/use-admin-queries'

export function ClientsPage() {
  const navigate = useNavigate()
  const [createOpen, setCreateOpen] = useState(false)
  const { page, setPage, pageSize, queryParams } = useConsolePagination()
  const clientsQuery = useAdminClients(queryParams)
  const clients = clientsQuery.data?.data ?? []
  const meta = clientsQuery.data?.meta

  return (
    <div>
      <header className="mb-10 flex flex-col justify-between gap-6 md:flex-row md:items-end">
        <div>
          <h1 className="font-headline text-4xl font-extrabold tracking-tight text-on-surface">Client Management</h1>
          <p className="mt-1 text-on-surface-variant">OAuth 2.0 / OIDC application registrations.</p>
        </div>
        <button
          type="button"
          onClick={() => setCreateOpen(true)}
          className="flex items-center gap-2 rounded-xl bg-primary px-5 py-2.5 text-sm font-bold text-on-primary shadow-lg shadow-primary/20 transition-opacity hover:opacity-90"
        >
          <MaterialIcon name="add" className="text-sm" />
          Register Client
        </button>
      </header>

      {clientsQuery.isLoading ? (
        <ConsoleLoadingState />
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
                  <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-primary-container">
                    <MaterialIcon name="devices" className="text-primary text-2xl" />
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
                <MaterialIcon name="chevron_right" className="text-on-surface-variant" />
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
