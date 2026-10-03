import { zodResolver } from '@hookform/resolvers/zod'
import { Loader2Icon, ShieldCheckIcon, ShieldIcon } from 'lucide-react'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { useNavigate } from 'react-router'

import { ApiError } from '@/api/types'
import { useAuth } from '@/hooks/use-auth'
import { mfaCodeSchema, type MfaCodeFormValues } from '@/auth/schemas'
import { authCardClassName } from '@/components/layout/auth-card'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

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
      <div className={authCardClassName}>
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
    <div className={authCardClassName}>
      <div className="space-y-6">
        <div className="flex items-center gap-3">
          <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-primary-container">
            <ShieldIcon className="h-6 w-6 text-primary" aria-hidden />
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
          <div className="space-y-2">
            <Label htmlFor="code">Authentication code</Label>
            <Input
              id="code"
              inputMode="numeric"
              autoComplete="one-time-code"
              placeholder="123456"
              className="text-center font-mono text-lg tracking-widest"
              {...form.register('code')}
            />
            {form.formState.errors.code ? (
              <p className="mt-1 text-sm text-destructive">{form.formState.errors.code.message}</p>
            ) : null}
            <p className="mt-2 text-xs text-on-surface-variant">Enter a TOTP code or recovery code from your authenticator app.</p>
          </div>

          <Button type="submit" className="h-11 w-full" disabled={isSubmitting}>
            {isSubmitting ? <Loader2Icon className="h-4 w-4 animate-spin" /> : <ShieldCheckIcon className="h-4 w-4" aria-hidden />}
            Verify
          </Button>
        </form>
      </div>
    </div>
  )
}
