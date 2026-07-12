import { startRegistration } from '@simplewebauthn/browser'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { CheckCircle2Icon, FingerprintIcon, Loader2Icon } from 'lucide-react'
import { useState } from 'react'

import { listWebauthnCredentials, webauthnRegisterFinish, webauthnRegisterStart } from '@/api/client'
import type { WebauthnCredentialResponse } from '@/api/types'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { formatWebAuthnError, unwrapWebAuthnOptions } from '@/lib/webauthn-error'

const passkeysQueryKey = ['webauthn', 'credentials'] as const

function formatPasskeyDate(iso: string) {
  return new Date(iso).toLocaleDateString(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}

export function PasskeyRegisterPanel() {
  const queryClient = useQueryClient()
  const passkeysQuery = useQuery({
    queryKey: passkeysQueryKey,
    queryFn: async () => {
      const envelope = await listWebauthnCredentials({ page_size: 100 })
      return envelope.data
    },
  })
  const passkeys: WebauthnCredentialResponse[] = passkeysQuery.data ?? []
  const isLoadingList = passkeysQuery.isLoading
  const [deviceName, setDeviceName] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(false)
  const [showAddForm, setShowAddForm] = useState(false)

  async function handleRegister() {
    setError(null)
    setIsLoading(true)

    try {
      const start = await webauthnRegisterStart({ device_name: deviceName || undefined })
      const credential = await startRegistration({
        optionsJSON: unwrapWebAuthnOptions(start.data.options),
      })
      await webauthnRegisterFinish({
        session_token: start.data.session_token,
        credential,
      })
      setDeviceName('')
      setShowAddForm(false)
      await queryClient.invalidateQueries({ queryKey: passkeysQueryKey })
    } catch (err) {
      setError(formatWebAuthnError(err, 'Passkey registration failed'))
    } finally {
      setIsLoading(false)
    }
  }

  const hasPasskeys = passkeys.length > 0

  return (
    <div className="space-y-4 rounded-xl border p-4">
      <div>
        <h3 className="font-medium">Passkeys</h3>
        <p className="text-sm text-muted-foreground">Register a passkey for passwordless sign-in.</p>
      </div>

      {error ? (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      ) : null}

      {isLoadingList ? (
        <div className="flex items-center gap-2 text-sm text-muted-foreground">
          <Loader2Icon className="h-4 w-4 animate-spin" />
          Loading passkeys…
        </div>
      ) : hasPasskeys ? (
        <>
          <Alert variant="success">
            <AlertDescription className="flex items-center gap-2">
              <CheckCircle2Icon className="h-4 w-4 shrink-0 text-green-600 dark:text-green-400" />
              {passkeys.length === 1
                ? '1 passkey registered on this account.'
                : `${passkeys.length} passkeys registered on this account.`}
            </AlertDescription>
          </Alert>

          <ul className="space-y-2">
            {passkeys.map((passkey) => (
              <li
                key={passkey.id}
                className="flex items-center justify-between gap-3 rounded-lg bg-muted/40 px-3 py-2 text-sm"
              >
                <span className="flex items-center gap-2 font-medium">
                  <FingerprintIcon className="h-4 w-4 text-muted-foreground" />
                  {passkey.device_name || 'Passkey'}
                </span>
                <span className="text-muted-foreground">{formatPasskeyDate(passkey.created_at)}</span>
              </li>
            ))}
          </ul>
        </>
      ) : null}

      {hasPasskeys && !showAddForm ? (
        <Button type="button" variant="outline" onClick={() => setShowAddForm(true)}>
          <FingerprintIcon className="h-4 w-4" />
          Add another passkey
        </Button>
      ) : null}

      {!hasPasskeys || showAddForm ? (
        <div className={hasPasskeys ? 'space-y-4 border-t pt-4' : 'space-y-4'}>
          {hasPasskeys ? (
            <p className="text-sm font-medium">Register another device</p>
          ) : null}

          <div className="space-y-2">
            <Label htmlFor="device_name">Device name (optional)</Label>
            <Input
              id="device_name"
              placeholder="Work MacBook"
              value={deviceName}
              onChange={(event) => setDeviceName(event.target.value)}
            />
          </div>

          <div className="flex flex-wrap gap-2">
            <Button type="button" disabled={isLoading} onClick={() => void handleRegister()}>
              {isLoading ? (
                <Loader2Icon className="h-4 w-4 animate-spin" />
              ) : (
                <FingerprintIcon className="h-4 w-4" />
              )}
              {hasPasskeys ? 'Register passkey' : 'Add passkey'}
            </Button>
            {hasPasskeys ? (
              <Button
                type="button"
                variant="ghost"
                disabled={isLoading}
                onClick={() => {
                  setShowAddForm(false)
                  setDeviceName('')
                  setError(null)
                }}
              >
                Cancel
              </Button>
            ) : null}
          </div>
        </div>
      ) : null}
    </div>
  )
}
