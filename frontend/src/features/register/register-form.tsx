import { zodResolver } from '@hookform/resolvers/zod'
import { Loader2Icon } from 'lucide-react'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { Link } from 'react-router-dom'

import { ApiError } from '@/api/types'
import { MaterialIcon } from '@/components/icons/material-icon'
import { useAuth } from '@/hooks/use-auth'
import { registerSchema, type RegisterFormValues } from '@/auth/schemas'
import { Alert, AlertDescription } from '@/components/ui/alert'

export function RegisterForm() {
  const { register: registerUser } = useAuth()
  const [error, setError] = useState<string | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)

  const form = useForm<RegisterFormValues>({
    resolver: zodResolver(registerSchema),
    defaultValues: { email: '', password: '', first_name: '', last_name: '' },
  })

  async function onSubmit(values: RegisterFormValues) {
    setError(null)
    setIsSubmitting(true)
    try {
      await registerUser(values)
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Registration failed')
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <div className="glass-panel overflow-hidden rounded-xl shadow-2xl shadow-on-surface/5 ring-1 ring-outline-variant/15">
      <div className="space-y-6 p-8 md:p-10">
        <div className="text-center">
          <h2 className="font-headline text-2xl font-bold text-on-surface">Request Access</h2>
          <p className="mt-1 text-sm text-on-surface-variant">Create your GateForge identity credentials.</p>
        </div>

        {error ? (
          <Alert variant="destructive">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        ) : null}

        <form className="space-y-4" onSubmit={form.handleSubmit(onSubmit)}>
          <div className="grid gap-4 sm:grid-cols-2">
            <div>
              <label htmlFor="first_name" className="mb-2 block font-label text-xs font-bold uppercase tracking-wider text-on-surface-variant">
                First name
              </label>
              <input
                id="first_name"
                autoComplete="given-name"
                className="w-full rounded-lg border-none bg-surface-container-low px-4 py-3 text-on-surface focus:ring-2 focus:ring-primary"
                {...form.register('first_name')}
              />
            </div>
            <div>
              <label htmlFor="last_name" className="mb-2 block font-label text-xs font-bold uppercase tracking-wider text-on-surface-variant">
                Last name
              </label>
              <input
                id="last_name"
                autoComplete="family-name"
                className="w-full rounded-lg border-none bg-surface-container-low px-4 py-3 text-on-surface focus:ring-2 focus:ring-primary"
                {...form.register('last_name')}
              />
            </div>
          </div>

          <div>
            <label htmlFor="email" className="mb-2 block font-label text-xs font-bold uppercase tracking-wider text-on-surface-variant">
              Institutional Email
            </label>
            <input
              id="email"
              type="email"
              autoComplete="email"
              placeholder="architect@gateforge.io"
              className="w-full rounded-lg border-none bg-surface-container-low px-4 py-3 text-on-surface placeholder:text-outline/50 focus:ring-2 focus:ring-primary"
              {...form.register('email')}
            />
            {form.formState.errors.email ? (
              <p className="mt-1 text-sm text-destructive">{form.formState.errors.email.message}</p>
            ) : null}
          </div>

          <div>
            <label htmlFor="password" className="mb-2 block font-label text-xs font-bold uppercase tracking-wider text-on-surface-variant">
              Password
            </label>
            <input
              id="password"
              type="password"
              autoComplete="new-password"
              className="w-full rounded-lg border-none bg-surface-container-low px-4 py-3 text-on-surface focus:ring-2 focus:ring-primary"
              {...form.register('password')}
            />
            {form.formState.errors.password ? (
              <p className="mt-1 text-sm text-destructive">{form.formState.errors.password.message}</p>
            ) : null}
          </div>

          <button
            type="submit"
            disabled={isSubmitting}
            className="flex w-full items-center justify-center gap-2 rounded-lg bg-primary py-4 font-headline font-bold text-on-primary shadow-lg shadow-primary/20 transition-all hover:bg-primary-dim disabled:opacity-60"
          >
            {isSubmitting ? <Loader2Icon className="h-4 w-4 animate-spin" /> : <MaterialIcon name="person_add" filled />}
            Create Account
          </button>
        </form>

        <p className="text-center text-sm text-on-surface-variant">
          Already have access?{' '}
          <Link to="/login" className="font-bold text-primary hover:underline">
            Sign in
          </Link>
        </p>
      </div>
    </div>
  )
}
