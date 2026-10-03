import { useState } from 'react'

import { ConsolePageHeader } from '@/components/layout/console-page-header'
import { ConsoleEmptyState, ConsoleErrorState, ConsoleLoadingState } from '@/features/admin/console-state'
import { ProviderConfigPanel } from '@/features/admin/provider-config-panel'
import { usePatchIdentityProvider, useTenantIdentityProviders } from '@/features/admin/use-admin-queries'
import { FederationProviderIcon } from '@/features/login/federation-provider-icons'

const defaultTenantId = import.meta.env.VITE_DEFAULT_TENANT_ID ?? '00000000-0000-0000-0000-000000000001'

export function IdentityProvidersPage() {
  const providersQuery = useTenantIdentityProviders(defaultTenantId)
  const patchProvider = usePatchIdentityProvider(defaultTenantId)
  const providers = providersQuery.data?.data ?? []
  const [editingProvider, setEditingProvider] = useState<string | null>(null)

  return (
    <div>
      <ConsolePageHeader
        title="Identity providers"
        description="Configure upstream OAuth credentials. Secrets stay encrypted in the database."
      />

      {providersQuery.isLoading ? (
        <ConsoleLoadingState />
      ) : providersQuery.isError ? (
        <ConsoleErrorState message="Could not load identity providers." />
      ) : providers.length === 0 ? (
        <ConsoleEmptyState title="No providers available" description="Supported identity providers will appear here for configuration." />
      ) : (
        <div className="grid gap-4 md:grid-cols-2">
          {providers.map((p) => (
            <div key={p.provider} className="rounded-xl bg-surface-container-lowest p-6 ghost-border">
              <div className="mb-4 flex items-start justify-between">
                <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-primary-container">
                  <FederationProviderIcon provider={p.provider} className="h-6 w-6" />
                </div>
                <div className="flex flex-col items-end gap-1">
                  <span
                    className={`rounded-full px-2 py-0.5 text-[10px] font-bold uppercase ${
                      p.enabled
                        ? 'bg-success-container text-on-success-container'
                        : 'bg-warning-container text-on-warning-container'
                    }`}
                  >
                    {p.enabled ? 'Enabled' : 'Disabled'}
                  </span>
                  <span
                    className={`rounded-full px-2 py-0.5 text-[10px] font-bold uppercase ${
                      p.configured
                        ? 'bg-primary-container text-on-primary-container'
                        : 'bg-surface-container text-on-surface-variant'
                    }`}
                  >
                    {p.configured ? 'Configured' : 'Not configured'}
                  </span>
                </div>
              </div>
              <h3 className="font-headline text-lg font-bold">{p.name}</h3>
              <span className="mt-1 inline-block rounded-md bg-tertiary-container px-2 py-0.5 text-[10px] font-bold text-on-tertiary-container">
                OAuth 2.0 / OIDC
              </span>

              {p.oauth_client_id ? (
                <p className="mt-3 truncate font-mono text-xs text-on-surface-variant">client: {p.oauth_client_id}</p>
              ) : null}

              <div className="mt-4 flex items-center justify-between">
                <p className="font-mono text-xs text-on-surface-variant">tenant: {p.tenant_id.slice(0, 8)}…</p>
                <button
                  type="button"
                  onClick={() => setEditingProvider((current) => (current === p.provider ? null : p.provider))}
                  className="rounded-lg bg-surface-container px-3 py-1.5 text-xs font-bold text-on-surface hover:bg-surface-container-high"
                >
                  {editingProvider === p.provider ? 'Close' : 'Configure'}
                </button>
              </div>

              {editingProvider === p.provider ? (
                <ProviderConfigPanel
                  provider={p}
                  isSaving={patchProvider.isPending}
                  onSave={async (body) => {
                    await patchProvider.mutateAsync({ provider: p.provider, body })
                    setEditingProvider(null)
                  }}
                  onCancel={() => setEditingProvider(null)}
                />
              ) : null}
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
