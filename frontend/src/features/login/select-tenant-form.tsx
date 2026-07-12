import { useState } from 'react'
import { useNavigate } from 'react-router-dom'

import { selectTenant } from '@/api/client'
import type { TenantSummary } from '@/api/types'
import { useAuth } from '@/hooks/use-auth'
import { Button } from '@/components/ui/button'
import { setTokens } from '@/auth/token-store'

type SelectTenantFormProps = {
  tenants: TenantSummary[]
  selectionToken: string
  rememberMe?: boolean
}

export function SelectTenantForm({ tenants, selectionToken, rememberMe = false }: SelectTenantFormProps) {
  const navigate = useNavigate()
  const { refreshProfile } = useAuth()
  const [selected, setSelected] = useState(tenants[0]?.id ?? '')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!selected) return
    setLoading(true)
    setError(null)
    try {
      const envelope = await selectTenant({
        selection_token: selectionToken,
        tenant_id: selected,
        remember_me: rememberMe,
      })
      setTokens(envelope.data, rememberMe)
      sessionStorage.removeItem('tenant_selection')
      await refreshProfile()
      navigate('/console')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to select tenant')
    } finally {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={onSubmit} className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-on-surface">Choose organization</h1>
        <p className="mt-2 text-sm text-on-surface-variant">
          Your account has access to multiple organizations. Select one to continue.
        </p>
      </div>

      <ul className="space-y-2">
        {tenants.map((tenant) => (
          <li key={tenant.id}>
            <label className="flex cursor-pointer items-center gap-3 rounded-lg border border-outline-variant p-4 hover:bg-surface-container">
              <input
                type="radio"
                name="tenant"
                value={tenant.id}
                checked={selected === tenant.id}
                onChange={() => setSelected(tenant.id)}
              />
              <span>
                <span className="block font-medium text-on-surface">{tenant.name || tenant.id}</span>
                {tenant.domain ? (
                  <span className="text-xs text-on-surface-variant">{tenant.domain}</span>
                ) : null}
              </span>
            </label>
          </li>
        ))}
      </ul>

      {error ? <p className="text-sm text-error">{error}</p> : null}

      <Button type="submit" className="w-full" disabled={loading || !selected}>
        {loading ? 'Continuing…' : 'Continue'}
      </Button>
    </form>
  )
}
