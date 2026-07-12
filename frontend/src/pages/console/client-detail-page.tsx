import { useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'

import { MaterialIcon } from '@/components/icons/material-icon'
import { ClientUsagePanel } from '@/features/admin/client-usage-panel'
import { ConsoleErrorState, ConsoleLoadingState } from '@/features/admin/console-state'
import { DeleteClientDialog } from '@/features/admin/delete-client-dialog'
import { EditClientDialog } from '@/features/admin/edit-client-dialog'
import { useAdminClient } from '@/features/admin/use-admin-queries'

const defaultTenantId = import.meta.env.VITE_DEFAULT_TENANT_ID ?? ''
const devClientId = 'oidc-dev'

export function ClientDetailPage() {
  const { clientId = '' } = useParams()
  const navigate = useNavigate()
  const clientQuery = useAdminClient(clientId)

  const [editOpen, setEditOpen] = useState(false)
  const [deleteOpen, setDeleteOpen] = useState(false)

  const client = clientQuery.data?.data
  const isDevClient = client?.tenant_id === defaultTenantId && client?.client_id === devClientId

  if (clientQuery.isLoading) {
    return <ConsoleLoadingState label="Loading client…" />
  }

  if (clientQuery.isError || !client) {
    return (
      <div className="p-6">
        <ConsoleErrorState message="Could not load OAuth client." />
        <Link to="/console/clients" className="mt-4 inline-flex text-sm font-semibold text-primary">
          Back to clients
        </Link>
      </div>
    )
  }

  return (
    <div>
      <div className="mb-6">
        <Link
          to="/console/clients"
          className="inline-flex items-center gap-1 text-sm font-semibold text-on-surface-variant hover:text-primary"
        >
          <MaterialIcon name="arrow_back" className="text-base" />
          Clients
        </Link>
      </div>

      <header className="mb-10 flex flex-col justify-between gap-6 md:flex-row md:items-end">
        <div>
          <h1 className="font-headline text-4xl font-extrabold tracking-tight text-on-surface">
            {client.name || client.client_id}
          </h1>
          <p className="mt-1 font-mono text-sm text-on-surface-variant">{client.client_id}</p>
          <div className="mt-3 flex flex-wrap gap-2">
            <span className="rounded-md bg-tertiary-container px-2 py-0.5 text-[10px] font-bold text-on-tertiary-container">
              {client.is_public ? 'Public' : 'Confidential'}
            </span>
            {!client.is_public && client.client_secret_set ? (
              <span className="rounded-md bg-secondary-container px-2 py-0.5 text-[10px] font-bold text-on-secondary-container">
                Secret configured
              </span>
            ) : null}
            {client.grant_types?.map((grant) => (
              <span
                key={grant}
                className="rounded-md bg-secondary-container px-2 py-0.5 text-[10px] font-bold text-on-secondary-container"
              >
                {grant}
              </span>
            ))}
          </div>
          <div className="mt-3 flex flex-wrap gap-4 text-sm text-on-surface-variant">
            <span>Tenant: {client.tenant_id.slice(0, 8)}…</span>
            <span>Created {new Date(client.created_at).toLocaleDateString()}</span>
          </div>
        </div>
        <div className="flex flex-wrap gap-3">
          <button
            type="button"
            onClick={() => setEditOpen(true)}
            className="flex items-center gap-2 rounded-xl border border-outline-variant px-5 py-2.5 text-sm font-bold text-on-surface transition-colors hover:bg-surface-container"
          >
            <MaterialIcon name="edit" className="text-sm" />
            Edit
          </button>
          <button
            type="button"
            onClick={() => setDeleteOpen(true)}
            disabled={isDevClient}
            title={isDevClient ? 'The default development client cannot be deleted' : undefined}
            className="flex items-center gap-2 rounded-xl border border-error/30 px-5 py-2.5 text-sm font-bold text-error transition-colors hover:bg-error/5 disabled:cursor-not-allowed disabled:opacity-50"
          >
            <MaterialIcon name="delete" className="text-sm" />
            Delete
          </button>
        </div>
      </header>

      <section className="mb-8 rounded-xl bg-surface-container-lowest p-6 ghost-border">
        <h2 className="mb-4 font-headline text-lg font-bold text-on-surface">Configuration</h2>
        <dl className="grid gap-4 sm:grid-cols-2">
          <div>
            <dt className="text-[10px] font-bold uppercase tracking-wider text-on-surface-variant">Redirect URIs</dt>
            <dd className="mt-1 space-y-1">
              {client.redirect_uris?.length ? (
                client.redirect_uris.map((uri) => (
                  <p key={uri} className="break-all font-mono text-xs text-on-surface">
                    {uri}
                  </p>
                ))
              ) : (
                <p className="text-sm text-on-surface-variant">—</p>
              )}
            </dd>
          </div>
          <div>
            <dt className="text-[10px] font-bold uppercase tracking-wider text-on-surface-variant">Scopes</dt>
            <dd className="mt-2 flex flex-wrap gap-2">
              {client.scopes?.length ? (
                client.scopes.map((scope) => (
                  <span
                    key={scope}
                    className="rounded-md bg-surface-container-low px-2 py-0.5 font-mono text-xs text-on-surface"
                  >
                    {scope}
                  </span>
                ))
              ) : (
                <p className="text-sm text-on-surface-variant">—</p>
              )}
            </dd>
          </div>
        </dl>
      </section>

      <section className="rounded-xl bg-surface-container-lowest p-6 ghost-border">
        <h2 className="mb-4 font-headline text-lg font-bold text-on-surface">Usage</h2>
        <ClientUsagePanel clientId={client.id} />
      </section>

      <EditClientDialog open={editOpen} client={client} onOpenChange={setEditOpen} />
      <DeleteClientDialog
        open={deleteOpen}
        clientId={client.id}
        clientName={client.name || client.client_id}
        onOpenChange={setDeleteOpen}
        onDeleted={() => navigate('/console/clients')}
      />
    </div>
  )
}
