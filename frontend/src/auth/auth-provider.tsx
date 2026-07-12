import { useCallback, useEffect, useMemo, useState, type ReactNode } from 'react'
import { useNavigate } from 'react-router-dom'

import {
  getMe,
  loginUser,
  logoutUser,
  registerUser,
  verifyMfaChallenge,
} from '@/api/client'
import { ApiError, isMfaChallenge, isTenantSelection, type LoginRequest } from '@/api/types'
import { AuthContext, type AuthContextValue } from '@/auth/auth-context-value'
import type { LoginFormValues, RegisterFormValues } from '@/auth/schemas'
import { clearTokens, getAccessToken, setTokens } from '@/auth/token-store'

export function AuthProvider({ children }: { children: ReactNode }) {
  const navigate = useNavigate()
  const [user, setUser] = useState<AuthContextValue['user']>(null)
  const [isLoading, setIsLoading] = useState(true)

  const refreshProfile = useCallback(async () => {
    try {
      const envelope = await getMe()
      setUser(envelope.data)
    } catch {
      setUser(null)
      clearTokens()
    }
  }, [])

  useEffect(() => {
    void (async () => {
      setIsLoading(true)
      const { hydrateRefreshTokenFromStorage, refreshTokens: doRefresh } = await import(
        '@/auth/token-store'
      )
      hydrateRefreshTokenFromStorage()
      if (!getAccessToken()) {
        await doRefresh()
      }
      await refreshProfile()
      setIsLoading(false)
    })()
  }, [refreshProfile])

  const login = useCallback(
    async (values: LoginFormValues, returnTo?: string): Promise<'mfa' | 'done'> => {
      const body: LoginRequest = {
        email: values.email,
        password: values.password,
        remember_me: values.remember_me,
        return_to: returnTo,
      }

      const envelope = await loginUser(body)
      if (isMfaChallenge(envelope.data)) {
        sessionStorage.setItem('mfa_ticket', envelope.data.mfa_ticket)
        sessionStorage.setItem('mfa_remember_me', String(values.remember_me))
        if (returnTo) {
          sessionStorage.setItem('mfa_return_to', returnTo)
        }
        navigate('/mfa/challenge')
        return 'mfa'
      }

      if (isTenantSelection(envelope.data)) {
        sessionStorage.setItem(
          'tenant_selection',
          JSON.stringify({
            selection_token: envelope.data.selection_token,
            tenants: envelope.data.tenants,
            remember_me: values.remember_me,
          }),
        )
        navigate('/select-tenant')
        return 'done'
      }

      setTokens(envelope.data, values.remember_me)
      await refreshProfile()

      if (returnTo) {
        window.location.href = returnTo
        return 'done'
      }

      navigate('/console')
      return 'done'
    },
    [navigate, refreshProfile],
  )

  const register = useCallback(
    async (values: RegisterFormValues) => {
      await registerUser({
        email: values.email,
        password: values.password,
        first_name: values.first_name,
        last_name: values.last_name,
      })
      navigate('/login')
    },
    [navigate],
  )

  const verifyMfa = useCallback(
    async (ticket: string, code: string, returnTo?: string) => {
      const rememberMe = sessionStorage.getItem('mfa_remember_me') === 'true'
      const envelope = await verifyMfaChallenge({ mfa_ticket: ticket, code })
      setTokens(envelope.data, rememberMe)
      sessionStorage.removeItem('mfa_ticket')
      sessionStorage.removeItem('mfa_return_to')
      sessionStorage.removeItem('mfa_remember_me')
      await refreshProfile()

      if (returnTo) {
        window.location.href = returnTo
        return
      }

      navigate('/console')
    },
    [navigate, refreshProfile],
  )

  const logout = useCallback(async () => {
    try {
      await logoutUser()
    } catch (error) {
      if (!(error instanceof ApiError) || error.status !== 401) {
        throw error
      }
    } finally {
      clearTokens()
      setUser(null)
      navigate('/login')
    }
  }, [navigate])

  const value = useMemo<AuthContextValue>(
    () => ({
      user,
      isAuthenticated: Boolean(user),
      isLoading,
      login,
      register,
      verifyMfa,
      logout,
      refreshProfile,
    }),
    [user, isLoading, login, register, verifyMfa, logout, refreshProfile],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}
