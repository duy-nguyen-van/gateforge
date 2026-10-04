import type { ReactNode } from 'react'

import { LoginShell } from '@/components/layout/login-shell'

export function AuthLayout({ children }: { children: ReactNode }) {
  return <LoginShell>{children}</LoginShell>
}
