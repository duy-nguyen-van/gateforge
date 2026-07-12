import { QRCodeSVG } from 'qrcode.react'
import { CheckCircle2Icon, CopyIcon, KeyRoundIcon, Loader2Icon } from 'lucide-react'
import { useState } from 'react'

import { generateRecoveryCodes, setupTotp, verifyTotpEnrollment } from '@/api/client'
import { ApiError } from '@/api/types'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useAuth } from '@/hooks/use-auth'

export function TotpSetupPanel() {
  const { user, refreshProfile } = useAuth()
  const totpEnabled = user?.mfa_enabled ?? false
  const [secret, setSecret] = useState<string | null>(null)
  const [otpauthUri, setOtpauthUri] = useState<string | null>(null)
  const [code, setCode] = useState('')
  const [recoveryCodes, setRecoveryCodes] = useState<string[] | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(false)

  async function handleSetup() {
    setError(null)
    setIsLoading(true)
    try {
      const envelope = await setupTotp()
      setSecret(envelope.data.secret)
      setOtpauthUri(envelope.data.otpauth_uri)
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Failed to start TOTP setup')
    } finally {
      setIsLoading(false)
    }
  }

  async function handleVerify() {
    setError(null)
    setIsLoading(true)
    try {
      await verifyTotpEnrollment({ code })
      await refreshProfile()
      setSecret(null)
      setOtpauthUri(null)
      setCode('')
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Invalid verification code')
    } finally {
      setIsLoading(false)
    }
  }

  async function handleRecoveryCodes() {
    setError(null)
    setIsLoading(true)
    try {
      const envelope = await generateRecoveryCodes()
      setRecoveryCodes(envelope.data.codes)
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Failed to generate recovery codes')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="space-y-4 rounded-xl border p-4">
      <div>
        <h3 className="font-medium">Authenticator app (TOTP)</h3>
        <p className="text-sm text-muted-foreground">Use Google Authenticator, 1Password, or similar apps.</p>
      </div>

      {error ? (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      ) : null}

      {totpEnabled ? (
        <Alert variant="success">
          <AlertDescription className="flex items-center gap-2">
            <CheckCircle2Icon className="h-4 w-4 shrink-0 text-green-600 dark:text-green-400" />
            Already set up — your authenticator app is active on this account.
          </AlertDescription>
        </Alert>
      ) : !secret ? (
        <Button type="button" disabled={isLoading} onClick={() => void handleSetup()}>
          {isLoading ? <Loader2Icon className="h-4 w-4 animate-spin" /> : <KeyRoundIcon className="h-4 w-4" />}
          Set up authenticator
        </Button>
      ) : (
        <div className="space-y-4">
          {otpauthUri ? (
            <div className="flex flex-col items-center gap-3 rounded-lg bg-muted/40 p-4">
              <QRCodeSVG value={otpauthUri} size={180} />
              <p className="break-all text-center text-xs text-muted-foreground">{secret}</p>
            </div>
          ) : null}

          <div className="space-y-2">
            <Label htmlFor="totp_code">Verification code</Label>
            <Input
              id="totp_code"
              inputMode="numeric"
              placeholder="123456"
              value={code}
              onChange={(event) => setCode(event.target.value)}
            />
          </div>

          <Button type="button" disabled={isLoading || code.length < 6} onClick={() => void handleVerify()}>
            Verify and enable
          </Button>
        </div>
      )}

      {totpEnabled ? (
        <div className="border-t pt-4">
          <h4 className="mb-2 text-sm font-medium">Recovery codes</h4>
          <Button type="button" variant="outline" disabled={isLoading} onClick={() => void handleRecoveryCodes()}>
            Generate recovery codes
          </Button>
          {recoveryCodes ? (
            <div className="mt-3 space-y-2 rounded-lg bg-muted/40 p-3 font-mono text-sm">
              {recoveryCodes.map((item) => (
                <div key={item} className="flex items-center justify-between gap-2">
                  <span>{item}</span>
                  <button
                    type="button"
                    className="text-muted-foreground hover:text-foreground"
                    onClick={() => void navigator.clipboard.writeText(item)}
                    aria-label="Copy code"
                  >
                    <CopyIcon className="h-4 w-4" />
                  </button>
                </div>
              ))}
            </div>
          ) : null}
        </div>
      ) : null}
    </div>
  )
}
