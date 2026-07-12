import { useState } from 'react'

import { ApiError } from '@/api/types'
import { MaterialIcon } from '@/components/icons/material-icon'
import { ConsolePortal } from '@/components/layout/console-portal'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useAdminTenants, useCreateAdminClient } from '@/features/admin/use-admin-queries'

const defaultTenantId = import.meta.env.VITE_DEFAULT_TENANT_ID ?? ''

interface CreateClientDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onCreated?: (clientId: string) => void
}

function parseLines(value: string): string[] {
  return value
    .split('\n')
    .map((line) => line.trim())
    .filter(Boolean)
}

export function CreateClientDialog({ open, onOpenChange, onCreated }: CreateClientDialogProps) {
  const createClient = useCreateAdminClient()
  const tenantsQuery = useAdminTenants({ page: 1, page_size: 100 })
  const tenants = tenantsQuery.data?.data ?? []

  const [tenantId, setTenantId] = useState(defaultTenantId)
  const [name, setName] = useState('')
  const [clientId, setClientId] = useState('')
  const [isPublic, setIsPublic] = useState(true)
  const [redirectUris, setRedirectUris] = useState('http://localhost:5173/callback')
  const [grantTypes, setGrantTypes] = useState('authorization_code')
  const [scopes, setScopes] = useState('openid\nemail\nprofile')
  const [error, setError] = useState<string | null>(null)
  const [createdSecret, setCreatedSecret] = useState<string | null>(null)
  const [createdClientRecordId, setCreatedClientRecordId] = useState<string | null>(null)

  function resetForm() {
    setTenantId(defaultTenantId)
    setName('')
    setClientId('')
    setIsPublic(true)
    setRedirectUris('http://localhost:5173/callback')
    setGrantTypes('authorization_code')
    setScopes('openid\nemail\nprofile')
    setError(null)
    setCreatedSecret(null)
    setCreatedClientRecordId(null)
  }

  function handleClose() {
    onOpenChange(false)
    resetForm()
  }

  function handleSecretDismiss() {
    const recordId = createdClientRecordId
    handleClose()
    if (recordId && onCreated) {
      onCreated(recordId)
    }
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError(null)

    const uris = parseLines(redirectUris)
    if (uris.length === 0) {
      setError('At least one redirect URI is required.')
      return
    }

    try {
      const result = await createClient.mutateAsync({
        tenant_id: tenantId,
        client_id: clientId.trim() || undefined,
        name: name.trim(),
        is_public: isPublic,
        redirect_uris: uris,
        grant_types: parseLines(grantTypes),
        scopes: parseLines(scopes),
      })
      const recordId = result.data?.id
      const secret = result.data?.client_secret
      if (secret && recordId) {
        setCreatedSecret(secret)
        setCreatedClientRecordId(recordId)
        return
      }
      handleClose()
      if (recordId && onCreated) {
        onCreated(recordId)
      }
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Could not register client.')
    }
  }

  if (!open) {
    return null
  }

  if (createdSecret) {
    return (
      <ConsolePortal>
        <div className="console-modal-scrim fixed inset-0 z-50 flex items-center justify-center p-4">
        <div
          role="dialog"
          aria-modal="true"
          aria-labelledby="client-secret-title"
          className="console-modal-panel w-full max-w-md rounded-xl p-6"
        >
          <h2 id="client-secret-title" className="font-headline text-xl font-bold text-on-surface">
            Save your client secret
          </h2>
          <p className="mt-2 text-sm text-on-surface-variant">
            This secret is shown only once. Copy it now — you will not be able to view it again.
          </p>
          <div className="mt-4 rounded-lg bg-surface-container-low p-3">
            <p className="break-all font-mono text-sm text-on-surface">{createdSecret}</p>
          </div>
          <div className="mt-6 flex justify-end">
            <Button type="button" onClick={handleSecretDismiss}>
              I have saved the secret
            </Button>
          </div>
        </div>
      </div>
      </ConsolePortal>
    )
  }

  return (
    <ConsolePortal>
      <div className="console-modal-scrim fixed inset-0 z-50 flex items-center justify-center p-4">
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby="create-client-title"
        className="console-modal-panel w-full max-w-lg rounded-xl p-6"
      >
        <div className="mb-6 flex items-start justify-between">
          <div>
            <h2 id="create-client-title" className="font-headline text-xl font-bold text-on-surface">
              Register client
            </h2>
            <p className="mt-1 text-sm text-on-surface-variant">Create an OAuth 2.0 / OIDC application.</p>
          </div>
          <button
            type="button"
            onClick={handleClose}
            className="rounded-lg p-1 text-on-surface-variant hover:bg-surface-container"
            aria-label="Close"
          >
            <MaterialIcon name="close" />
          </button>
        </div>

        {error ? (
          <Alert variant="destructive" className="mb-4">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        ) : null}

        <form onSubmit={(e) => void handleSubmit(e)} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="create-client-tenant">Tenant</Label>
            <select
              id="create-client-tenant"
              required
              value={tenantId}
              onChange={(e) => setTenantId(e.target.value)}
              className="w-full rounded-lg border border-outline-variant bg-surface-container-lowest px-3 py-2 text-sm"
            >
              {tenants.length === 0 ? (
                <option value={defaultTenantId}>Default tenant</option>
              ) : (
                tenants.map((t) => (
                  <option key={t.id} value={t.id}>
                    {t.name || t.id}
                  </option>
                ))
              )}
            </select>
          </div>

          <div className="space-y-2">
            <Label htmlFor="create-client-name">Display name</Label>
            <Input
              id="create-client-name"
              required
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="My Application"
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="create-client-id">Client ID (optional)</Label>
            <Input
              id="create-client-id"
              value={clientId}
              onChange={(e) => setClientId(e.target.value)}
              placeholder="Auto-generated if empty"
            />
          </div>

          <div className="space-y-2">
            <Label>Client type</Label>
            <div className="flex gap-4">
              <label className="flex items-center gap-2 text-sm">
                <input
                  type="radio"
                  name="client-type"
                  checked={isPublic}
                  onChange={() => setIsPublic(true)}
                />
                Public (PKCE)
              </label>
              <label className="flex items-center gap-2 text-sm">
                <input
                  type="radio"
                  name="client-type"
                  checked={!isPublic}
                  onChange={() => setIsPublic(false)}
                />
                Confidential
              </label>
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="create-client-redirects">Redirect URIs (one per line)</Label>
            <textarea
              id="create-client-redirects"
              required
              rows={3}
              value={redirectUris}
              onChange={(e) => setRedirectUris(e.target.value)}
              className="w-full rounded-lg border border-outline-variant bg-surface-container-lowest px-3 py-2 font-mono text-sm"
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="create-client-grants">Grant types (one per line)</Label>
            <textarea
              id="create-client-grants"
              rows={2}
              value={grantTypes}
              onChange={(e) => setGrantTypes(e.target.value)}
              className="w-full rounded-lg border border-outline-variant bg-surface-container-lowest px-3 py-2 font-mono text-sm"
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="create-client-scopes">Scopes (one per line)</Label>
            <textarea
              id="create-client-scopes"
              rows={3}
              value={scopes}
              onChange={(e) => setScopes(e.target.value)}
              className="w-full rounded-lg border border-outline-variant bg-surface-container-lowest px-3 py-2 font-mono text-sm"
            />
          </div>

          <div className="flex justify-end gap-3 pt-2">
            <Button type="button" variant="outline" onClick={handleClose}>
              Cancel
            </Button>
            <Button type="submit" disabled={createClient.isPending || !name.trim() || !tenantId}>
              {createClient.isPending ? 'Registering…' : 'Register client'}
            </Button>
          </div>
        </form>
      </div>
    </div>
    </ConsolePortal>
  )
}
