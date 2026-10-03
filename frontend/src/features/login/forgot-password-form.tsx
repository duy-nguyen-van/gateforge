import { zodResolver } from '@hookform/resolvers/zod'
import { Loader2Icon } from 'lucide-react'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { Link } from 'react-router'

import { forgotPassword } from '@/api/client'
import { ApiError } from '@/api/types'
import { forgotPasswordSchema, type ForgotPasswordValues } from '@/auth/schemas'
import { authCardClassName } from '@/components/layout/auth-card'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

export function ForgotPasswordForm() {
  const [error, setError] = useState<string | null>(null)
  const [sent, setSent] = useState(false)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const form = useForm<ForgotPasswordValues>({
    resolver: zodResolver(forgotPasswordSchema),
    defaultValues: { email: '' },
  })

  async function onSubmit(values: ForgotPasswordValues) {
    setError(null)
    setIsSubmitting(true)
    try {
      await forgotPassword(values.email)
      setSent(true)
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Could not send reset email')
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <div className={authCardClassName}>
      <div className="mb-8 text-center">
        <h2 className="font-headline text-2xl font-semibold text-on-surface">Reset password</h2>
        <p className="mt-1 text-sm text-on-surface-variant">We will email a link if the account exists.</p>
      </div>
      {sent ? (
        <Alert>
          <AlertDescription>If an account exists, a reset link has been sent.</AlertDescription>
        </Alert>
      ) : (
        <form className="space-y-4" onSubmit={form.handleSubmit(onSubmit)}>
          {error ? (
            <Alert variant="destructive">
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          ) : null}
          <div className="space-y-2">
            <Label htmlFor="email">Work email</Label>
            <Input id="email" type="email" autoComplete="email" placeholder="name@company.com" {...form.register('email')} />
            {form.formState.errors.email ? (
              <p className="text-sm text-destructive">{form.formState.errors.email.message}</p>
            ) : null}
          </div>
          <Button type="submit" className="h-11 w-full" disabled={isSubmitting}>
            {isSubmitting ? <Loader2Icon className="h-4 w-4 animate-spin" /> : null}
            Send reset link
          </Button>
        </form>
      )}
      <p className="mt-6 text-center text-sm text-on-surface-variant">
        <Link to="/login" className="font-semibold text-primary hover:underline underline-offset-4">
          Back to sign in
        </Link>
      </p>
    </div>
  )
}
