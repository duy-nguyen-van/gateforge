import { Navigate, Outlet } from 'react-router-dom'

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

export function GuestRoute() {
  const { isAuthenticated, isLoading, user } = useAuth()

  if (isLoading) {
    return <GateForgeLoading label="Loading session…" />
  }

  if (isAuthenticated) {
    return <Navigate to={user?.is_platform_admin ? '/console' : '/settings/profile'} replace />
  }

  return <Outlet />
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
