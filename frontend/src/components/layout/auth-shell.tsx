import type { ReactNode } from 'react'

import { GateForgeBrand } from '@/components/brand/gateforge-brand'

const footerLinks = [
  { label: 'Privacy Policy', href: '#' },
  { label: 'Terms of Service', href: '#' },
  { label: 'Security Compliance', href: '#' },
  { label: 'Support', href: '#' },
] as const

export function AuthFooter() {
  return (
    <footer className="fixed bottom-0 flex w-full items-center justify-center gap-8 border-t border-slate-200/15 bg-transparent py-6 dark:border-slate-800/15">
      <span className="text-xs uppercase tracking-wide text-slate-400 dark:text-slate-500">
        © 2026 GateForge. All rights reserved.
      </span>
      <div className="hidden gap-8 md:flex">
        {footerLinks.map(({ label, href }) => (
          <a
            key={label}
            href={href}
            className="text-xs uppercase tracking-wide text-slate-400 transition-all underline-offset-4 hover:text-primary hover:underline dark:text-slate-500 dark:hover:text-slate-300"
          >
            {label}
          </a>
        ))}
      </div>
    </footer>
  )
}

export function AuthTrustBadges() {
  return (
    <div className="flex items-center justify-center gap-6 opacity-60 grayscale transition-all hover:grayscale-0 hover:opacity-100">
      <div className="flex items-center gap-2">
        <span className="material-symbols-outlined text-sm">verified_user</span>
        <span className="font-label text-[10px] font-bold uppercase tracking-widest">SOC2 Type II</span>
      </div>
      <div className="h-3 w-px bg-outline-variant" />
      <div className="flex items-center gap-2">
        <span className="material-symbols-outlined text-sm">policy</span>
        <span className="font-label text-[10px] font-bold uppercase tracking-widest">GDPR Compliant</span>
      </div>
      <div className="h-3 w-px bg-outline-variant" />
      <div className="flex items-center gap-2">
        <span className="material-symbols-outlined text-sm">lock</span>
        <span className="font-label text-[10px] font-bold uppercase tracking-widest">ISO 27001</span>
      </div>
    </div>
  )
}

interface AuthShellProps {
  children: ReactNode
  showBrand?: boolean
  showTrustBadges?: boolean
  showFooter?: boolean
  brandVariant?: 'gateforge' | 'sovereign'
}

export function AuthShell({
  children,
  showBrand = true,
  showTrustBadges = true,
  showFooter = true,
  brandVariant = 'sovereign',
}: AuthShellProps) {
  return (
    <div className="relative flex min-h-screen flex-col bg-surface font-body text-on-surface selection:bg-primary-container selection:text-on-primary-container">
      <div className="fixed inset-0 -z-10 bg-grid-pattern opacity-40" />
      <div className="fixed inset-0 -z-20 bg-gradient-to-tr from-surface via-surface-container-low to-surface-container-high" />

      <main className="flex flex-grow items-center justify-center p-6 pb-28 md:p-12 md:pb-32">
        <div className="w-full max-w-[480px] space-y-8">
          {showBrand ? <GateForgeBrand variant={brandVariant} linkTo="/" /> : null}
          {children}
          {showTrustBadges ? <AuthTrustBadges /> : null}
        </div>
      </main>

      {showFooter ? <AuthFooter /> : null}
    </div>
  )
}

export function AuthCardFooter({ children }: { children: ReactNode }) {
  return (
    <div className="flex items-center justify-center border-t border-outline-variant/15 bg-surface-container-low/50 px-8 py-4">
      {children}
    </div>
  )
}
