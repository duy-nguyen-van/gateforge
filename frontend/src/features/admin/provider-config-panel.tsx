import { Loader2Icon } from 'lucide-react'
import { useState } from 'react'

import type { AdminIdentityProviderResponse, PatchIdentityProviderRequest } from '@/api/types'
import { ApiError } from '@/api/types'
import { MaterialIcon } from '@/components/icons/material-icon'
import { Alert, AlertDescription } from '@/components/ui/alert'

interface ProviderConfigPanelProps {
  provider: AdminIdentityProviderResponse
  isSaving: boolean
  onSave: (body: PatchIdentityProviderRequest) => Promise<void>
  onCancel: () => void
}

export function ProviderConfigPanel({ provider, isSaving, onSave, onCancel }: ProviderConfigPanelProps) {
  const [clientId, setClientId] = useState(provider.oauth_client_id ?? '')
  const [clientSecret, setClientSecret] = useState('')
  const [enabled, setEnabled] = useState(provider.enabled)
  const [error, setError] = useState<string | null>(null)

  const displayName = provider.name
  const idPrefix = provider.provider

  async function handleSubmit(event: React.FormEvent) {
    event.preventDefault()
    setError(null)

    const body: PatchIdentityProviderRequest = { enabled }
    if (clientId.trim()) {
      body.oauth_client_id = clientId.trim()
    }
    if (clientSecret.trim()) {
      body.oauth_client_secret = clientSecret.trim()
    }

    if (enabled && !clientId.trim() && !provider.oauth_client_secret_set) {
      setError(`Client ID and client secret are required before enabling ${displayName} sign-in.`)
      return
    }
    if (enabled && !clientSecret.trim() && !provider.oauth_client_secret_set) {
      setError(`Client secret is required before enabling ${displayName} sign-in.`)
      return
    }

    try {
      await onSave(body)
    } catch (err) {
      setError(err instanceof ApiError ? err.message : `Could not save ${displayName} provider settings.`)
    }
  }

  return (
    <form onSubmit={(event) => void handleSubmit(event)} className="mt-4 space-y-4 border-t border-surface-container pt-4">
      {error ? (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      ) : null}

      <div>
        <label
          htmlFor={`${idPrefix}-client-id`}
          className="mb-1 block text-xs font-bold uppercase tracking-wider text-on-surface-variant"
        >
          OAuth Client ID
        </label>
        <input
          id={`${idPrefix}-client-id`}
          type="text"
          value={clientId}
          onChange={(event) => setClientId(event.target.value)}
          placeholder="OAuth client ID"
          className="w-full rounded-lg border-none bg-surface-container-low px-3 py-2 font-mono text-sm text-on-surface focus:ring-2 focus:ring-primary"
        />
      </div>

      <div>
        <label
          htmlFor={`${idPrefix}-client-secret`}
          className="mb-1 block text-xs font-bold uppercase tracking-wider text-on-surface-variant"
        >
          OAuth Client Secret
        </label>
        <input
          id={`${idPrefix}-client-secret`}
          type="password"
          value={clientSecret}
          onChange={(event) => setClientSecret(event.target.value)}
          placeholder={provider.oauth_client_secret_set ? 'Leave blank to keep existing secret' : 'Enter client secret'}
          className="w-full rounded-lg border-none bg-surface-container-low px-3 py-2 font-mono text-sm text-on-surface focus:ring-2 focus:ring-primary"
        />
      </div>

      {provider.redirect_uri ? (
        <div>
          <p className="mb-1 text-xs font-bold uppercase tracking-wider text-on-surface-variant">Authorized redirect URI</p>
          <p className="rounded-lg bg-surface-container-low px-3 py-2 font-mono text-xs text-on-surface">{provider.redirect_uri}</p>
          <p className="mt-1 text-xs text-on-surface-variant">
            Add this exact URI in your {displayName} OAuth app settings.
            {provider.setup_console_url ? (
              <>
                {' '}
                <a
                  href={provider.setup_console_url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="font-semibold text-primary hover:underline"
                >
                  Open setup console
                </a>
              </>
            ) : null}
          </p>
        </div>
      ) : null}

      <label className="flex items-center gap-2 text-sm text-on-surface">
        <input
          type="checkbox"
          checked={enabled}
          onChange={(event) => setEnabled(event.target.checked)}
          className="rounded border-outline-variant"
        />
        Enable {displayName} sign-in for this tenant
      </label>

      <div className="flex items-center gap-2">
        <button
          type="submit"
          disabled={isSaving}
          className="flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-bold text-on-primary disabled:opacity-60"
        >
          {isSaving ? <Loader2Icon className="h-4 w-4 animate-spin" /> : <MaterialIcon name="save" className="text-sm" />}
          Save configuration
        </button>
        <button
          type="button"
          onClick={onCancel}
          className="rounded-lg px-4 py-2 text-sm font-semibold text-on-surface-variant hover:text-on-surface"
        >
          Cancel
        </button>
      </div>
    </form>
  )
}
