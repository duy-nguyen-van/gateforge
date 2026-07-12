import { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'

import { exchangeSessionLogin } from '@/api/client'
import { isTenantSelection } from '@/api/types'
import { setTokens } from '@/auth/token-store'
import { GateForgeLoading } from '@/components/brand/gateforge-loading'
import { useAuth } from '@/hooks/use-auth'
import { Alert, AlertDescription } from '@/components/ui/alert'

export function FederationCompletePage() {
  const navigate = useNavigate()
  const { refreshProfile } = useAuth()
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    void (async () => {
      try {
        const envelope = await exchangeSessionLogin()

        if (isTenantSelection(envelope.data)) {
          sessionStorage.setItem(
            'tenant_selection',
            JSON.stringify({
              selection_token: envelope.data.selection_token,
              tenants: envelope.data.tenants,
              remember_me: false,
            }),
          )
          navigate('/select-tenant', { replace: true })
          return
        }

        if ('access_token' in envelope.data) {
          setTokens(envelope.data)
          await refreshProfile()
          navigate('/console', { replace: true })
          return
        }

        setError('Unexpected sign-in response. Try again or use another method.')
      } catch {
        setError('Could not complete sign-in. Try again or use another method.')
      }
    })()
  }, [navigate, refreshProfile])

  if (error) {
    return (
      <div className="mx-auto max-w-md px-6 py-16 text-center">
        <Alert variant="destructive" className="mb-6 text-left">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
        <Link to="/login" className="text-sm font-semibold text-primary hover:underline underline-offset-4">
          Back to sign in
        </Link>
      </div>
    )
  }

  return <GateForgeLoading label="Completing sign-in…" className="min-h-[40vh] bg-transparent" />
}
