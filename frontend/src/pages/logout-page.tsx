import { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'

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

  return (
    <div className="flex min-h-screen items-center justify-center text-sm text-muted-foreground">
      Signing out…
    </div>
  )
}
