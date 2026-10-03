import { zodResolver } from '@hookform/resolvers/zod'
import { useQuery } from '@tanstack/react-query'
import { Loader2Icon, MailIcon } from 'lucide-react'
import { useState } from 'react'
import { Controller, useForm, useWatch } from 'react-hook-form'
import { Link, useSearchParams } from 'react-router'

import { federationStartUrl, federationCompleteReturnTo, listFederationProviders } from '@/api/client'
import { ApiError } from '@/api/types'
import { useAuth } from '@/hooks/use-auth'
import { loginSchema, type LoginFormValues } from '@/auth/schemas'
import { authCardClassName } from '@/components/layout/auth-card'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { SocialLoginGrid } from '@/features/login/social-provider-icons'
import { PasskeyLoginButton } from '@/features/webauthn/passkey-login'

const defaultTenantId = import.meta.env.VITE_DEFAULT_TENANT_ID ?? '00000000-0000-0000-0000-000000000001'

export function LoginForm() {
  const { login } = useAuth()
  const [searchParams] = useSearchParams()
  const returnTo = searchParams.get('return_to') ?? undefined
  const [error, setError] = useState<string | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [showEmailForm, setShowEmailForm] = useState(false)

  const form = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: { email: '', password: '', remember_me: false },
  })

  const email = useWatch({ control: form.control, name: 'email' })
  const federationProvidersQuery = useQuery({
    queryKey: ['federation', 'providers', defaultTenantId],
    queryFn: () => listFederationProviders(defaultTenantId),
  })
  const providers = federationProvidersQuery.data?.data ?? []
  const federationReturnTo = returnTo ?? federationCompleteReturnTo()
  const federationLinks = providers.map((p) => ({
    provider: p.provider,
    name: p.name,
    href: federationStartUrl(p.provider, federationReturnTo),
  }))

  async function onSubmit(values: LoginFormValues) {
    setError(null)
    setIsSubmitting(true)
    try {
      await login(values, returnTo)
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Sign in failed')
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <div className={authCardClassName}>
      <div className="mb-8 text-center">
        <h2 className="mb-2 font-headline text-2xl font-semibold text-on-surface">Sign in</h2>
        <p className="text-sm text-on-surface-variant">Use a passkey, or sign in with your work email.</p>
      </div>

      {error ? (
        <Alert variant="destructive" className="mb-6">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      ) : null}

      <PasskeyLoginButton
        email={email ?? ''}
        returnTo={returnTo}
        variant="revamp"
        onEmailRequired={() => setShowEmailForm(true)}
      />

      <div className="mt-4 space-y-4">
        <Button
          type="button"
          variant="outline"
          className="h-11 w-full"
          onClick={() => setShowEmailForm((open) => !open)}
        >
          <MailIcon className="h-5 w-5" aria-hidden />
          Sign in with email
        </Button>

        {showEmailForm ? (
          <form className="space-y-4 border-t border-surface-container pt-4" onSubmit={form.handleSubmit(onSubmit)}>
            <div className="space-y-2">
              <Label htmlFor="email">Work email</Label>
              <Input
                id="email"
                type="email"
                autoComplete="email"
                placeholder="name@company.com"
                {...form.register('email')}
              />
              {form.formState.errors.email ? (
                <p className="text-sm text-destructive">{form.formState.errors.email.message}</p>
              ) : null}
            </div>
            <div className="space-y-2">
              <Label htmlFor="password">Password</Label>
              <Input
                id="password"
                type="password"
                autoComplete="current-password"
                {...form.register('password')}
              />
              {form.formState.errors.password ? (
                <p className="text-sm text-destructive">{form.formState.errors.password.message}</p>
              ) : null}
            </div>
            <div className="flex items-center gap-2">
              <Controller
                name="remember_me"
                control={form.control}
                render={({ field }) => (
                  <Checkbox
                    id="remember_me"
                    checked={field.value}
                    onCheckedChange={(checked) => field.onChange(checked === true)}
                  />
                )}
              />
              <Label htmlFor="remember_me">Remember me</Label>
            </div>
            <Button type="submit" className="h-11 w-full" disabled={isSubmitting}>
              {isSubmitting ? <Loader2Icon className="h-4 w-4 animate-spin" /> : null}
              Sign in
            </Button>
          </form>
        ) : null}
      </div>

      {federationLinks.length > 0 ? (
        <div className="mt-6">
          <div className="relative mb-4 flex items-center">
            <div className="flex-grow border-t border-surface-container" />
            <span className="mx-4 text-xs font-medium text-on-surface-variant">Or continue with</span>
            <div className="flex-grow border-t border-surface-container" />
          </div>
          <SocialLoginGrid providers={federationLinks} />
        </div>
      ) : null}

      <div className="mt-8 flex flex-col items-center gap-3 text-sm">
        <Link to="/forgot-password" className="font-medium text-primary hover:underline underline-offset-4">
          Forgot password
        </Link>
        <p className="text-on-surface-variant">
          Need an account?{' '}
          <Link to="/register" className="font-semibold text-primary hover:underline underline-offset-4">
            Create account
          </Link>
        </p>
      </div>
    </div>
  )
}
