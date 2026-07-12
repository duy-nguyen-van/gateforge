import { Navigate } from 'react-router-dom'

import { SelectTenantForm } from '@/features/login/select-tenant-form'

export function SelectTenantPage() {
  const raw = sessionStorage.getItem('tenant_selection')
  if (!raw) {
    return <Navigate to="/login" replace />
  }

  let payload: {
    selection_token: string
    tenants: { id: string; name: string; domain: string; role: string }[]
    remember_me?: boolean
  }
  try {
    payload = JSON.parse(raw) as typeof payload
  } catch {
    return <Navigate to="/login" replace />
  }

  return (
    <SelectTenantForm
      tenants={payload.tenants}
      selectionToken={payload.selection_token}
      rememberMe={payload.remember_me}
    />
  )
}
