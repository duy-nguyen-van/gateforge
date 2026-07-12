import { useState } from 'react'

import { ApiError } from '@/api/types'
import { MaterialIcon } from '@/components/icons/material-icon'
import { ConsolePortal } from '@/components/layout/console-portal'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useAddTenantMember, useAdminTenants } from '@/features/admin/use-admin-queries'

const defaultTenantIdEnv = import.meta.env.VITE_DEFAULT_TENANT_ID ?? ''

interface AddMemberDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  /** When set, locks the tenant and hides the tenant picker. */
  defaultTenantId?: string
}

export function AddMemberDialog({ open, onOpenChange, defaultTenantId }: AddMemberDialogProps) {
  const tenantsQuery = useAdminTenants({ page_size: 100 })
  const addMember = useAddTenantMember()
  const [email, setEmail] = useState('')
  const [tenantId, setTenantId] = useState(defaultTenantId ?? defaultTenantIdEnv)
  const [role, setRole] = useState<'member' | 'admin'>('member')
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState(false)

  const tenants = tenantsQuery.data?.data ?? []
  const lockedTenant = Boolean(defaultTenantId)
  const effectiveTenantId = lockedTenant
    ? defaultTenantId!
    : tenantId || tenants[0]?.id || ''

  function handleClose() {
    onOpenChange(false)
    setEmail('')
    setTenantId(defaultTenantId ?? defaultTenantIdEnv)
    setRole('member')
    setError(null)
    setSuccess(false)
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError(null)
    setSuccess(false)

    if (!effectiveTenantId) {
      setError('Select a tenant.')
      return
    }

    try {
      await addMember.mutateAsync({
        tenantId: effectiveTenantId,
        body: { email: email.trim(), role },
      })
      setSuccess(true)
      setTimeout(handleClose, 1200)
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Could not add member.')
    }
  }

  if (!open) {
    return null
  }

  return (
    <ConsolePortal>
    <div className="console-modal-scrim fixed inset-0 z-50 flex items-center justify-center p-4">
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby="add-member-title"
        className="console-modal-panel w-full max-w-md rounded-xl p-6"
      >
        <div className="mb-6 flex items-start justify-between">
          <div>
            <h2 id="add-member-title" className="font-headline text-xl font-bold text-on-surface">
              Add member
            </h2>
            <p className="mt-1 text-sm text-on-surface-variant">
              The user must already be registered in the system.
            </p>
          </div>
          <button
            type="button"
            onClick={handleClose}
            className="rounded-lg p-1 text-on-surface-variant hover:bg-surface-container"
            aria-label="Close"
          >
            <MaterialIcon name="close" />
          </button>
        </div>

        {error ? (
          <Alert variant="destructive" className="mb-4">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        ) : null}

        {success ? (
          <Alert variant="success" className="mb-4">
            <AlertDescription>Member added successfully.</AlertDescription>
          </Alert>
        ) : null}

        <form onSubmit={(e) => void handleSubmit(e)} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="add-member-email">Email</Label>
            <Input
              id="add-member-email"
              type="email"
              required
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="user@example.com"
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="add-member-tenant">Tenant</Label>
            {lockedTenant ? (
              <p className="text-sm text-on-surface">
                {tenants.find((t) => t.id === defaultTenantId)?.name ||
                  defaultTenantId?.slice(0, 8) + '…'}
              </p>
            ) : tenantsQuery.isLoading ? (
              <p className="text-sm text-on-surface-variant">Loading tenants…</p>
            ) : tenants.length === 0 ? (
              <p className="text-sm text-on-surface-variant">No tenants available.</p>
            ) : (
              <select
                id="add-member-tenant"
                value={effectiveTenantId}
                onChange={(e) => setTenantId(e.target.value)}
                className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
              >
                {tenants.map((t) => (
                  <option key={t.id} value={t.id}>
                    {t.name || t.domain || t.id.slice(0, 8)}
                  </option>
                ))}
              </select>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="add-member-role">Role</Label>
            <select
              id="add-member-role"
              value={role}
              onChange={(e) => setRole(e.target.value as 'member' | 'admin')}
              className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
            >
              <option value="member">Member</option>
              <option value="admin">Admin</option>
            </select>
          </div>

          <div className="flex justify-end gap-3 pt-2">
            <Button type="button" variant="outline" onClick={handleClose}>
              Cancel
            </Button>
            <Button type="submit" disabled={addMember.isPending || !effectiveTenantId}>
              {addMember.isPending ? 'Adding…' : 'Add member'}
            </Button>
          </div>
        </form>
      </div>
    </div>
    </ConsolePortal>
  )
}
