import { useQueryClient } from '@tanstack/react-query'
import { ShieldCheckIcon } from 'lucide-react'

import { TotpSetupPanel } from '@/features/mfa/totp-setup'
import { PasskeyRegisterPanel } from '@/features/webauthn/passkey-register'

const steps = [
  'Open an authenticator app, or use a passkey on this device.',
  'Scan the QR code. The new account is named GateForge.',
  'Enter the 6-digit code. The console opens after that code is accepted.',
]

export function AdminMfaEnrollment() {
  const queryClient = useQueryClient()

  function handleEnrolled() {
    void queryClient.invalidateQueries({ queryKey: ['admin'] })
  }

  return (
    <section className="overflow-hidden rounded-xl bg-surface-container-lowest ghost-border">
      <div className="border-b border-surface-container px-8 py-6">
        <div className="flex items-start gap-4">
          <span className="rounded-lg bg-primary-container p-2 text-primary">
            <ShieldCheckIcon className="h-5 w-5" aria-hidden />
          </span>
          <div>
            <h1 className="font-headline text-2xl font-extrabold tracking-tight text-on-surface">
              Add a second factor to open the console
            </h1>
            <p className="mt-2 max-w-2xl text-sm leading-relaxed text-on-surface-variant">
              Platform admin pages stay closed until this account has an authenticator app or a passkey.
              Sign-in already worked. This step only unlocks administration.
            </p>
          </div>
        </div>
        <ol className="mt-6 grid gap-3 sm:grid-cols-3">
          {steps.map((step, index) => (
            <li key={step} className="rounded-xl bg-surface-container-low px-4 py-3">
              <p className="font-label text-[10px] font-bold uppercase tracking-widest text-primary">
                Step {index + 1}
              </p>
              <p className="mt-1 text-sm leading-relaxed text-on-surface">{step}</p>
            </li>
          ))}
        </ol>
      </div>
      <div className="space-y-6 p-8">
        <TotpSetupPanel onEnrolled={handleEnrolled} />
        <div className="border-t border-surface-container pt-6">
          <PasskeyRegisterPanel onEnrolled={handleEnrolled} />
        </div>
      </div>
    </section>
  )
}
