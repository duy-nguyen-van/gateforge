import { zodResolver } from '@hookform/resolvers/zod'
import { Loader2Icon } from 'lucide-react'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { useNavigate } from 'react-router-dom'

import { ApiError } from '@/api/types'
import { MaterialIcon } from '@/components/icons/material-icon'
import { useAuth } from '@/hooks/use-auth'
import { mfaCodeSchema, type MfaCodeFormValues } from '@/auth/schemas'
import { Alert, AlertDescription } from '@/components/ui/alert'

export function MfaChallengeForm() {
  const { verifyMfa } = useAuth()
  const navigate = useNavigate()
  const [error, setError] = useState<string | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)

  const ticket = sessionStorage.getItem('mfa_ticket')
  const returnTo = sessionStorage.getItem('mfa_return_to') ?? undefined

  const form = useForm<MfaCodeFormValues>({
    resolver: zodResolver(mfaCodeSchema),
    defaultValues: { code: '' },
  })

  if (!ticket) {
    return (
      <div className="glass-panel rounded-xl p-8 ring-1 ring-outline-variant/15">
        <Alert>
          <AlertDescription>
            MFA session expired.{' '}
            <button type="button" className="font-bold text-primary underline" onClick={() => navigate('/login')}>
              Sign in again
            </button>
          </AlertDescription>
        </Alert>
      </div>
    )
  }

  async function onSubmit(values: MfaCodeFormValues) {
    setError(null)
    setIsSubmitting(true)
    try {
      await verifyMfa(ticket!, values.code, returnTo)
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Verification failed')
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <div className="glass-panel overflow-hidden rounded-xl shadow-2xl shadow-on-surface/5 ring-1 ring-outline-variant/15">
      <div className="space-y-6 p-8 md:p-10">
        <div className="flex items-center gap-3">
          <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-primary-container">
            <MaterialIcon name="security" filled className="text-primary text-2xl" />
          </div>
          <div>
            <h2 className="font-headline text-2xl font-bold text-on-surface">Verify Identity</h2>
            <p className="text-sm text-on-surface-variant">Complete two-factor authentication</p>
          </div>
        </div>

        {error ? (
          <Alert variant="destructive">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        ) : null}

        <form className="space-y-6" onSubmit={form.handleSubmit(onSubmit)}>
          <div>
            <label htmlFor="code" className="mb-2 block font-label text-xs font-bold uppercase tracking-wider text-on-surface-variant">
              Authentication Code
            </label>
            <input
              id="code"
              inputMode="numeric"
              autoComplete="one-time-code"
              placeholder="123456"
              className="w-full rounded-lg border-none bg-surface-container-low px-4 py-3 text-center font-mono text-lg tracking-widest text-on-surface focus:ring-2 focus:ring-primary"
              {...form.register('code')}
            />
            {form.formState.errors.code ? (
              <p className="mt-1 text-sm text-destructive">{form.formState.errors.code.message}</p>
            ) : null}
            <p className="mt-2 text-xs text-on-surface-variant">Enter a TOTP code or recovery code from your authenticator app.</p>
          </div>

          <button
            type="submit"
            disabled={isSubmitting}
            className="flex w-full items-center justify-center gap-2 rounded-lg bg-on-primary-fixed py-4 font-headline font-bold tracking-widest text-surface shadow-md transition-all hover:bg-inverse-surface disabled:opacity-60"
          >
            {isSubmitting ? <Loader2Icon className="h-4 w-4 animate-spin" /> : <MaterialIcon name="verified_user" />}
            VERIFY SESSION
          </button>
        </form>
      </div>
    </div>
  )
}
