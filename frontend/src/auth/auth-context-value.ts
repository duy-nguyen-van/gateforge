import { createContext } from 'react'

import type { UserResponse } from '@/api/types'
import type { LoginFormValues, RegisterFormValues } from '@/auth/schemas'

export interface AuthContextValue {
  user: UserResponse | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (values: LoginFormValues, returnTo?: string) => Promise<'mfa' | 'done'>
  register: (values: RegisterFormValues) => Promise<void>
  verifyMfa: (ticket: string, code: string, returnTo?: string) => Promise<void>
  logout: () => Promise<void>
  refreshProfile: () => Promise<void>
}

export const AuthContext = createContext<AuthContextValue | null>(null)
