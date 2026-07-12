import { useContext } from 'react'

import { AuthContext } from '@/auth/auth-context-value'

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider')
  }
  return context
}
