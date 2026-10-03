import { Navigate, Outlet } from 'react-router'

import { GateForgeLoading } from '@/components/brand/gateforge-loading'
import { useAuth } from '@/hooks/use-auth'

export function ProtectedRoute() {
  const { isAuthenticated, isLoading } = useAuth()

  if (isLoading) {
    return <GateForgeLoading label="Loading session…" />
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  return <Outlet />
}

function signedInDestination(isPlatformAdmin: boolean | undefined) {
  return isPlatformAdmin ? '/console' : '/settings/profile'
}

export function GuestRoute() {
  const { isAuthenticated, isLoading, user } = useAuth()

  if (isLoading) {
    return <GateForgeLoading label="Loading session…" />
  }

  if (isAuthenticated) {
    return <Navigate to={signedInDestination(user?.is_platform_admin)} replace />
  }

  return <Outlet />
}

export function RootRedirect() {
  const { isAuthenticated, isLoading, user } = useAuth()

  if (isLoading) {
    return <GateForgeLoading label="Loading session…" />
  }

  return (
    <Navigate
      to={isAuthenticated ? signedInDestination(user?.is_platform_admin) : '/login'}
      replace
    />
  )
}

export function AdminRoute() {
  const { user, isLoading } = useAuth()

  if (isLoading) {
    return <GateForgeLoading label="Loading session…" />
  }

  if (!user?.is_platform_admin) {
    return <Navigate to="/settings/profile" replace />
  }

  return <Outlet />
}
