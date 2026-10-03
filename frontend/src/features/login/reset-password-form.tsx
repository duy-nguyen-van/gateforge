import { zodResolver } from '@hookform/resolvers/zod'
import { Loader2Icon } from 'lucide-react'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { Link, useNavigate, useSearchParams } from 'react-router'

import { resetPassword } from '@/api/client'
import { ApiError } from '@/api/types'
import { resetPasswordSchema, type ResetPasswordValues } from '@/auth/schemas'
import { authCardClassName } from '@/components/layout/auth-card'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

export function ResetPasswordForm() {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const token = searchParams.get('token') ?? ''
  const [error, setError] = useState<string | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const form = useForm<ResetPasswordValues>({
    resolver: zodResolver(resetPasswordSchema),
    defaultValues: { password: '' },
  })

  async function onSubmit(values: ResetPasswordValues) {
    setError(null)
    setIsSubmitting(true)
    try {
      await resetPassword(token, values.password)
      navigate('/login')
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Could not reset password')
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <div className={authCardClassName}>
      <div className="mb-8 text-center">
        <h2 className="font-headline text-2xl font-semibold text-on-surface">Choose a new password</h2>
        <p className="mt-1 text-sm text-on-surface-variant">Use at least 12 characters.</p>
      </div>
      <form className="space-y-4" onSubmit={form.handleSubmit(onSubmit)}>
        {error ? (
          <Alert variant="destructive">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        ) : null}
        {!token ? (
          <Alert variant="destructive">
            <AlertDescription>This reset link is missing a token.</AlertDescription>
          </Alert>
        ) : null}
        <div className="space-y-2">
          <Label htmlFor="password">New password</Label>
          <Input id="password" type="password" autoComplete="new-password" {...form.register('password')} />
          {form.formState.errors.password ? (
            <p className="text-sm text-destructive">{form.formState.errors.password.message}</p>
          ) : null}
        </div>
        <Button type="submit" className="h-11 w-full" disabled={isSubmitting || token === ''}>
          {isSubmitting ? <Loader2Icon className="h-4 w-4 animate-spin" /> : null}
          Update password
        </Button>
      </form>
      <p className="mt-6 text-center text-sm text-on-surface-variant">
        <Link to="/login" className="font-semibold text-primary hover:underline underline-offset-4">
          Back to sign in
        </Link>
      </p>
    </div>
  )
}
