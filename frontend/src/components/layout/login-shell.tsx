import type { ReactNode } from 'react'

import { GateForgeBrand } from '@/components/brand/gateforge-brand'
import { MaterialIcon } from '@/components/icons/material-icon'

export function LoginBrand() {
  return (
    <header className="mb-10 text-center">
      <GateForgeBrand size="md" layout="stacked" showTagline linkTo="/" />
    </header>
  )
}

export function LoginTrustBadges() {
  const badges = [
    { icon: 'verified_user', label: 'FIDO2 Certified' },
    { icon: 'lock', label: 'AES-256 Encryption' },
    { icon: 'shield', label: 'SOC2 Type II' },
  ] as const

  return (
    <div className="mt-8 flex flex-nowrap items-center justify-center gap-4 opacity-60 grayscale transition-all duration-500 sm:gap-6 hover:grayscale-0 hover:opacity-100">
      {badges.map(({ icon, label }) => (
        <div key={label} className="flex shrink-0 items-center gap-2">
          <MaterialIcon name={icon} className="shrink-0 text-base leading-none" />
          <span className="whitespace-nowrap text-[10px] font-bold uppercase leading-none tracking-widest">
            {label}
          </span>
        </div>
      ))}
    </div>
  )
}

export function LoginDecorations() {
  return (
    <>
      <div className="fixed bottom-0 left-0 hidden p-8 md:block">
        <div className="flex flex-col gap-1">
          <div className="flex gap-1">
            <div className="h-1 w-1 bg-primary/20" />
            <div className="h-1 w-4 bg-primary/20" />
          </div>
          <div className="h-1 w-8 bg-primary/20" />
        </div>
      </div>
      <div className="fixed right-0 top-0 hidden p-8 md:block">
        <div className="flex flex-col items-end gap-1">
          <div className="h-1 w-12 bg-primary/10" />
          <div className="flex gap-1">
            <div className="h-1 w-4 bg-primary/10" />
            <div className="h-1 w-1 bg-primary/10" />
          </div>
        </div>
      </div>
    </>
  )
}

export function LoginShell({ children }: { children: ReactNode }) {
  return (
    <div className="dot-grid flex min-h-screen flex-col items-center justify-center bg-background p-6 font-body text-on-background">
      <LoginDecorations />
      <main className="flex w-full max-w-[420px] flex-col items-center">
        <LoginBrand />
        {children}
        <LoginTrustBadges />
      </main>
    </div>
  )
}
