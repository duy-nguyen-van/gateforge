import { MaterialIcon } from '@/components/icons/material-icon'
import { TotpSetupPanel } from '@/features/mfa/totp-setup'
import { PasskeyRegisterPanel } from '@/features/webauthn/passkey-register'

export function SecurityPage() {
  return (
    <div>
      <header className="mb-10">
        <h1 className="font-headline text-4xl font-extrabold tracking-tight text-on-surface">Security Settings</h1>
        <p className="mt-1 text-on-surface-variant">Manage MFA, recovery codes, and passkeys for your account.</p>
      </header>

      <div className="space-y-6">
        <section className="overflow-hidden rounded-xl bg-surface-container-lowest ghost-border">
          <div className="flex items-center gap-3 border-b border-surface-container px-8 py-5">
            <MaterialIcon name="security" className="rounded-lg bg-primary-container p-2 text-primary" />
            <div>
              <h2 className="font-headline text-lg font-bold">Authentication Methods</h2>
              <p className="text-sm text-on-surface-variant">Configure TOTP and passkey credentials</p>
            </div>
          </div>
          <div className="space-y-6 p-8">
            <TotpSetupPanel />
            <div className="border-t border-surface-container pt-6">
              <PasskeyRegisterPanel />
            </div>
          </div>
        </section>
      </div>
    </div>
  )
}
