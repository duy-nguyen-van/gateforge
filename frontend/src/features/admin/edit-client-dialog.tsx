import { useState } from 'react'

import type { AdminClientResponse } from '@/api/types'
import { ApiError } from '@/api/types'
import { MaterialIcon } from '@/components/icons/material-icon'
import { ConsolePortal } from '@/components/layout/console-portal'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useUpdateAdminClient } from '@/features/admin/use-admin-queries'

interface EditClientDialogProps {
  open: boolean
  client: AdminClientResponse | null
  onOpenChange: (open: boolean) => void
}

function parseLines(value: string): string[] {
  return value
    .split('\n')
    .map((line) => line.trim())
    .filter(Boolean)
}

function EditClientForm({
  client,
  onClose,
}: {
  client: AdminClientResponse
  onClose: () => void
}) {
  const updateClient = useUpdateAdminClient(client.id)
  const [name, setName] = useState(client.name)
  const [isPublic, setIsPublic] = useState(client.is_public)
  const [redirectUris, setRedirectUris] = useState(client.redirect_uris?.join('\n') ?? '')
  const [grantTypes, setGrantTypes] = useState(client.grant_types?.join('\n') ?? '')
  const [scopes, setScopes] = useState(client.scopes?.join('\n') ?? '')
  const [clientSecret, setClientSecret] = useState('')
  const [error, setError] = useState<string | null>(null)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError(null)

    const uris = parseLines(redirectUris)
    if (uris.length === 0) {
      setError('At least one redirect URI is required.')
      return
    }

    try {
      await updateClient.mutateAsync({
        name: name.trim(),
        is_public: isPublic,
        redirect_uris: uris,
        grant_types: parseLines(grantTypes),
        scopes: parseLines(scopes),
        client_secret: clientSecret.trim() || undefined,
      })
      onClose()
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Could not update client.')
    }
  }

  return (
    <>
      {error ? (
        <Alert variant="destructive" className="mb-4">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      ) : null}

      <form onSubmit={(e) => void handleSubmit(e)} className="space-y-4">
        <div className="space-y-2">
          <Label htmlFor="edit-client-name">Display name</Label>
          <Input
            id="edit-client-name"
            required
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-2">
            <Label>Client ID</Label>
            <Input value={client.client_id} readOnly className="font-mono text-xs" />
          </div>
          <div className="space-y-2">
            <Label>Tenant ID</Label>
            <Input value={client.tenant_id} readOnly className="font-mono text-xs" />
          </div>
        </div>

        <div className="space-y-2">
          <Label>Client type</Label>
          <div className="flex gap-4">
            <label className="flex items-center gap-2 text-sm">
              <input
                type="radio"
                name="edit-client-type"
                checked={isPublic}
                onChange={() => setIsPublic(true)}
              />
              Public (PKCE)
            </label>
            <label className="flex items-center gap-2 text-sm">
              <input
                type="radio"
                name="edit-client-type"
                checked={!isPublic}
                onChange={() => setIsPublic(false)}
              />
              Confidential
            </label>
          </div>
        </div>

        <div className="space-y-2">
          <Label htmlFor="edit-client-redirects">Redirect URIs (one per line)</Label>
          <textarea
            id="edit-client-redirects"
            required
            rows={3}
            value={redirectUris}
            onChange={(e) => setRedirectUris(e.target.value)}
            className="w-full rounded-lg border border-outline-variant bg-surface-container-lowest px-3 py-2 font-mono text-sm"
          />
        </div>

        <div className="space-y-2">
          <Label htmlFor="edit-client-grants">Grant types (one per line)</Label>
          <textarea
            id="edit-client-grants"
            rows={2}
            value={grantTypes}
            onChange={(e) => setGrantTypes(e.target.value)}
            className="w-full rounded-lg border border-outline-variant bg-surface-container-lowest px-3 py-2 font-mono text-sm"
          />
        </div>

        <div className="space-y-2">
          <Label htmlFor="edit-client-scopes">Scopes (one per line)</Label>
          <textarea
            id="edit-client-scopes"
            rows={3}
            value={scopes}
            onChange={(e) => setScopes(e.target.value)}
            className="w-full rounded-lg border border-outline-variant bg-surface-container-lowest px-3 py-2 font-mono text-sm"
          />
        </div>

        {!isPublic ? (
          <div className="space-y-2">
            <Label htmlFor="edit-client-secret">New client secret (optional)</Label>
            <Input
              id="edit-client-secret"
              type="password"
              value={clientSecret}
              onChange={(e) => setClientSecret(e.target.value)}
              placeholder={client.client_secret_set ? 'Leave blank to keep current' : 'Set a new secret'}
            />
          </div>
        ) : null}

        <div className="flex justify-end gap-3 pt-2">
          <Button type="button" variant="outline" onClick={onClose}>
            Cancel
          </Button>
          <Button type="submit" disabled={updateClient.isPending || !name.trim()}>
            {updateClient.isPending ? 'Saving…' : 'Save changes'}
          </Button>
        </div>
      </form>
    </>
  )
}

export function EditClientDialog({ open, client, onOpenChange }: EditClientDialogProps) {
  if (!open || !client) {
    return null
  }

  function handleClose() {
    onOpenChange(false)
  }

  return (
    <ConsolePortal>
    <div className="console-modal-scrim fixed inset-0 z-50 flex items-center justify-center p-4">
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby="edit-client-title"
        className="console-modal-panel w-full max-w-lg rounded-xl p-6"
      >
        <div className="mb-6 flex items-start justify-between">
          <h2 id="edit-client-title" className="font-headline text-xl font-bold text-on-surface">
            Edit client
          </h2>
          <button
            type="button"
            onClick={handleClose}
            className="rounded-lg p-1 text-on-surface-variant hover:bg-surface-container"
            aria-label="Close"
          >
            <MaterialIcon name="close" />
          </button>
        </div>

        <EditClientForm key={client.id} client={client} onClose={handleClose} />
      </div>
    </div>
    </ConsolePortal>
  )
}
