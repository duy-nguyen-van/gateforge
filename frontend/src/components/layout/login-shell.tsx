import type { ReactNode } from 'react'

import { GateForgeBrand } from '@/components/brand/gateforge-brand'

export function LoginBrand() {
  return (
    <header className="mb-10 text-center">
      <GateForgeBrand size="md" layout="stacked" showTagline linkTo="/" />
    </header>
  )
}

export function LoginShell({ children }: { children: ReactNode }) {
  return (
    <div className="flex min-h-[100dvh] flex-col items-center justify-center bg-background p-6 text-on-surface">
      <main className="flex w-full max-w-[420px] flex-col items-center">
        <LoginBrand />
        {children}
      </main>
    </div>
  )
}
