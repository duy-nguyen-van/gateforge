import { startAuthentication } from '@simplewebauthn/browser'
import { Loader2Icon } from 'lucide-react'
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'

import { webauthnLoginFinish, webauthnLoginStart } from '@/api/client'
import { isMfaChallenge, isTenantSelection } from '@/api/types'
import { formatWebAuthnError, unwrapWebAuthnOptions } from '@/lib/webauthn-error'
import { useAuth } from '@/hooks/use-auth'
import { setTokens } from '@/auth/token-store'
import { MaterialIcon } from '@/components/icons/material-icon'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { cn } from '@/lib/utils'

interface PasskeyLoginButtonProps {
  email: string
  rememberMe?: boolean
  returnTo?: string
  variant?: 'primary' | 'secondary' | 'revamp'
  onEmailRequired?: () => void
}

export function PasskeyLoginButton({
  email,
  rememberMe,
  returnTo,
  variant = 'secondary',
  onEmailRequired,
}: PasskeyLoginButtonProps) {
  const navigate = useNavigate()
  const { refreshProfile } = useAuth()
  const [error, setError] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(false)

  async function handlePasskeyLogin() {
    if (!email) {
      onEmailRequired?.()
      setError('Enter your email in the form below before using a passkey.')
      return
    }

    setError(null)
    setIsLoading(true)

    try {
      const tenantId = import.meta.env.VITE_DEFAULT_TENANT_ID
      const start = await webauthnLoginStart({ email, tenant_id: tenantId })
      const credential = await startAuthentication({
        optionsJSON: unwrapWebAuthnOptions(start.data.options),
      })
      const finish = await webauthnLoginFinish({
        email,
        tenant_id: tenantId,
        session_token: start.data.session_token,
        credential,
        remember_me: rememberMe,
        return_to: returnTo,
      })

      if (isMfaChallenge(finish.data)) {
        sessionStorage.setItem('mfa_ticket', finish.data.mfa_ticket)
        sessionStorage.setItem('mfa_remember_me', String(rememberMe ?? false))
        if (returnTo) {
          sessionStorage.setItem('mfa_return_to', returnTo)
        }
        navigate('/mfa/challenge')
        return
      }

      if (isTenantSelection(finish.data)) {
        sessionStorage.setItem(
          'tenant_selection',
          JSON.stringify({
            selection_token: finish.data.selection_token,
            tenants: finish.data.tenants,
            remember_me: rememberMe ?? false,
          }),
        )
        navigate('/select-tenant')
        return
      }

      setTokens(finish.data, rememberMe)
      await refreshProfile()
      if (returnTo) {
        window.location.href = returnTo
        return
      }
      navigate('/console')
    } catch (err) {
      setError(formatWebAuthnError(err, 'Passkey sign-in failed'))
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="space-y-2">
      {error ? (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      ) : null}
      <button
        type="button"
        disabled={isLoading}
        onClick={() => void handlePasskeyLogin()}
        className={cn(
          'flex w-full items-center justify-center gap-3 font-semibold transition-all active:scale-[0.98] disabled:opacity-60',
          variant === 'revamp' &&
            'mb-6 h-12 rounded-lg bg-gradient-to-r from-primary to-primary-dim text-on-primary shadow-lg shadow-primary/20 hover:opacity-90',
          variant === 'primary' &&
            'rounded-lg bg-primary px-6 py-4 text-on-primary shadow-lg shadow-primary/20 hover:bg-primary-dim',
          variant === 'secondary' &&
            'rounded-lg border border-outline-variant/20 bg-surface-container-highest px-6 py-4 text-on-surface hover:bg-surface-variant',
        )}
      >
        {isLoading ? (
          <Loader2Icon className="h-4 w-4 animate-spin" />
        ) : (
          <MaterialIcon name={variant === 'revamp' ? 'passkey' : 'fingerprint'} filled={variant !== 'secondary'} />
        )}
        <span className={variant === 'revamp' ? undefined : 'font-headline tracking-wide'}>Sign in with Passkey</span>
      </button>
    </div>
  )
}
