import { useEffect } from 'react'
import { useNavigate } from 'react-router'

import { GateForgeLoading } from '@/components/brand/gateforge-loading'
import { useAuth } from '@/hooks/use-auth'

export function LogoutPage() {
  const { logout } = useAuth()
  const navigate = useNavigate()

  useEffect(() => {
    void (async () => {
      try {
        await logout()
      } catch {
        navigate('/login')
      }
    })()
  }, [logout, navigate])

  return <GateForgeLoading label="Signing out…" />
}
