import type { ReactNode } from 'react'

import { AuthShell } from '@/components/layout/auth-shell'
import { LoginShell } from '@/components/layout/login-shell'

export function AuthLayout({
  children,
  variant = 'default',
}: {
  children: ReactNode
  title?: string
  subtitle?: string
  variant?: 'default' | 'revamp'
}) {
  if (variant === 'revamp') {
    return <LoginShell>{children}</LoginShell>
  }

  return <AuthShell>{children}</AuthShell>
}
