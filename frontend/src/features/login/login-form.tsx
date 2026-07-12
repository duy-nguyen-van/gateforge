import { zodResolver } from '@hookform/resolvers/zod'
import { useQuery } from '@tanstack/react-query'
import { Loader2Icon } from 'lucide-react'
import { useState } from 'react'
import { useForm, useWatch } from 'react-hook-form'
import { Link, useSearchParams } from 'react-router-dom'

import { federationStartUrl, federationCompleteReturnTo, listFederationProviders } from '@/api/client'
import { ApiError } from '@/api/types'
import { MaterialIcon } from '@/components/icons/material-icon'
import { useAuth } from '@/hooks/use-auth'
import { loginSchema, type LoginFormValues } from '@/auth/schemas'
import { Alert, AlertDescription } from '@/components/ui/alert'
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
    <div className="rounded-xl bg-surface-container-lowest p-8 shadow-[0_32px_64px_-12px_rgba(42,52,57,0.08)] ring-1 ring-on-surface/5 md:p-10">
      <div className="mb-8 text-center">
        <h2 className="mb-2 font-headline text-2xl font-bold text-on-surface">Welcome Back</h2>
        <p className="text-sm text-on-surface-variant">Secure biometric authentication is recommended.</p>
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

      <div className="relative mb-2 flex items-center py-4">
        <div className="flex-grow border-t border-surface-container" />
        <span className="mx-4 flex-shrink text-[10px] font-bold uppercase tracking-widest text-on-surface-variant">
          Or continue with
        </span>
        <div className="flex-grow border-t border-surface-container" />
      </div>

      <div className="space-y-4">
        <button
          type="button"
          onClick={() => setShowEmailForm((open) => !open)}
          className="flex h-11 w-full items-center justify-center gap-3 rounded-lg border border-transparent bg-surface-container-low font-medium text-on-surface transition-colors hover:border-outline-variant/30"
        >
          <MaterialIcon name="business_center" className="text-primary" />
          Institutional Email / SSO
        </button>

        {showEmailForm ? (
          <form className="space-y-4 border-t border-surface-container pt-4" onSubmit={form.handleSubmit(onSubmit)}>
            <div>
              <label htmlFor="email" className="mb-2 block text-xs font-bold uppercase tracking-wider text-on-surface-variant">
                Institutional Email
              </label>
              <input
                id="email"
                type="email"
                autoComplete="email"
                placeholder="architect@gateforge.io"
                className="w-full rounded-lg border-none bg-surface-container-low px-4 py-3 text-on-surface placeholder:text-outline/50 focus:bg-surface-container-lowest focus:ring-2 focus:ring-primary"
                {...form.register('email')}
              />
              {form.formState.errors.email ? (
                <p className="mt-1 text-sm text-destructive">{form.formState.errors.email.message}</p>
              ) : null}
            </div>
            <div>
              <label htmlFor="password" className="mb-2 block text-xs font-bold uppercase tracking-wider text-on-surface-variant">
                Password
              </label>
              <input
                id="password"
                type="password"
                autoComplete="current-password"
                placeholder="••••••••••••"
                className="w-full rounded-lg border-none bg-surface-container-low px-4 py-3 text-on-surface placeholder:text-outline/50 focus:bg-surface-container-lowest focus:ring-2 focus:ring-primary"
                {...form.register('password')}
              />
              {form.formState.errors.password ? (
                <p className="mt-1 text-sm text-destructive">{form.formState.errors.password.message}</p>
              ) : null}
            </div>
            <button
              type="submit"
              disabled={isSubmitting}
              className="flex h-11 w-full items-center justify-center gap-2 rounded-lg bg-on-primary-fixed font-headline text-sm font-bold tracking-wide text-surface shadow-md transition-all hover:bg-inverse-surface disabled:opacity-60"
            >
              {isSubmitting ? <Loader2Icon className="h-4 w-4 animate-spin" /> : null}
              Sign in
            </button>
          </form>
        ) : null}
      </div>

      <SocialLoginGrid providers={federationLinks} />

      <div className="mt-10 flex flex-col items-center gap-4">
        <a href="#" className="text-sm font-medium text-primary hover:underline underline-offset-4">
          Need help accessing your account?
        </a>
        <p className="max-w-[280px] text-center text-[12px] text-on-surface-variant">
          By continuing, you agree to GateForge&apos;s <br />
          <a href="#" className="font-semibold text-on-surface transition-colors hover:text-primary">
            Terms of Service
          </a>{' '}
          &amp;{' '}
          <a href="#" className="font-semibold text-on-surface transition-colors hover:text-primary">
            Privacy Policy
          </a>
        </p>
        <p className="text-center text-[11px] text-on-surface-variant">
          New deployment?{' '}
          <Link to="/register" className="font-bold text-primary hover:underline underline-offset-4">
            Request access
          </Link>
        </p>
      </div>
    </div>
  )
}
