import { zodResolver } from '@hookform/resolvers/zod'
import { useQuery } from '@tanstack/react-query'
import { Loader2Icon } from 'lucide-react'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { useNavigate } from 'react-router'
import { z } from 'zod'

import { acceptMemberInvite, previewMemberInvite } from '@/api/client'
import { ApiError, isMfaChallenge } from '@/api/types'
import { setTokens } from '@/auth/token-store'
import { authCardClassName } from '@/components/layout/auth-card'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useAuth } from '@/hooks/use-auth'

const inviteAcceptSchema = z.object({
  first_name: z.string().trim().min(1, 'First name is required').max(100),
  last_name: z.string().max(100),
  password: z.string().min(12, 'Password must be at least 12 characters').max(128),
})

type InviteAcceptValues = z.infer<typeof inviteAcceptSchema>

export function InviteAcceptForm({ token }: { token: string }) {
  const { refreshProfile } = useAuth()
  const navigate = useNavigate()
  const [error, setError] = useState<string | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const preview = useQuery({
    queryKey: ['member-invite', token],
    queryFn: () => previewMemberInvite(token),
  })
  const form = useForm<InviteAcceptValues>({
    resolver: zodResolver(inviteAcceptSchema),
    defaultValues: { first_name: '', last_name: '', password: '' },
  })

  async function onSubmit(values: InviteAcceptValues) {
    setError(null)
    setIsSubmitting(true)
    try {
      const envelope = await acceptMemberInvite({
        token,
        password: values.password,
        first_name: values.first_name.trim() || undefined,
        last_name: values.last_name.trim() || undefined,
      })
      if (isMfaChallenge(envelope.data)) {
        sessionStorage.setItem('mfa_ticket', envelope.data.mfa_ticket)
        sessionStorage.setItem('mfa_remember_me', 'false')
        navigate('/mfa/challenge')
        return
      }
      if (!('access_token' in envelope.data)) {
        setError('Could not sign in')
        return
      }
      setTokens(envelope.data, false)
      await refreshProfile()
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Could not sign in')
    } finally {
      setIsSubmitting(false)
    }
  }

  const invite = preview.data?.data
  const orgName = invite?.organization_name || 'this organization'

  return (
    <div className={authCardClassName}>
      <div className="mb-8 text-center">
        <h2 className="mb-2 font-headline text-2xl font-semibold text-on-surface">Join {orgName}</h2>
        <p className="text-sm text-on-surface-variant">
          Sign in to continue. If you do not have an account yet, choose a password and we will create one.
        </p>
      </div>

      {preview.isError ? (
        <Alert variant="destructive" className="mb-6">
          <AlertDescription>
            {preview.error instanceof ApiError ? preview.error.message : 'This invite is invalid or expired.'}
          </AlertDescription>
        </Alert>
      ) : null}

      {error ? (
        <Alert variant="destructive" className="mb-6">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      ) : null}

      <form className="space-y-4" onSubmit={form.handleSubmit(onSubmit)}>
        <div className="space-y-2">
          <Label htmlFor="invite-email">Email</Label>
          <Input id="invite-email" value={invite?.email ?? ''} readOnly disabled />
        </div>
        <div className="grid gap-4 sm:grid-cols-2">
          <div className="space-y-2">
            <Label htmlFor="invite-first-name">First name</Label>
            <Input id="invite-first-name" autoComplete="given-name" {...form.register('first_name')} />
            {form.formState.errors.first_name ? (
              <p className="text-sm text-destructive">{form.formState.errors.first_name.message}</p>
            ) : null}
          </div>
          <div className="space-y-2">
            <Label htmlFor="invite-last-name">Last name</Label>
            <Input id="invite-last-name" autoComplete="family-name" {...form.register('last_name')} />
          </div>
        </div>
        <div className="space-y-2">
          <Label htmlFor="invite-password">Password</Label>
          <Input
            id="invite-password"
            type="password"
            autoComplete="new-password"
            {...form.register('password')}
          />
          {form.formState.errors.password ? (
            <p className="text-sm text-destructive">{form.formState.errors.password.message}</p>
          ) : null}
        </div>
        <Button type="submit" className="h-11 w-full" disabled={isSubmitting || preview.isError}>
          {isSubmitting ? <Loader2Icon className="h-4 w-4 animate-spin" /> : null}
          Sign in
        </Button>
      </form>
    </div>
  )
}
