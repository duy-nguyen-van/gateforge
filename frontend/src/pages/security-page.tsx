import { ShieldIcon } from 'lucide-react'

import { ConsolePageHeader } from '@/components/layout/console-page-header'
import { TotpSetupPanel } from '@/features/mfa/totp-setup'
import { PasskeyRegisterPanel } from '@/features/webauthn/passkey-register'

export function SecurityPage() {
  return (
    <div>
      <ConsolePageHeader
        title="Security"
        description="Add an authenticator app or a passkey. The console stays closed until one of those is on this account."
      />

      <div className="space-y-6">
        <section className="overflow-hidden rounded-xl bg-surface-container-lowest ghost-border">
          <div className="flex items-center gap-3 border-b border-surface-container px-8 py-5">
            <span className="rounded-lg bg-primary-container p-2 text-primary">
              <ShieldIcon className="h-5 w-5" aria-hidden />
            </span>
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
