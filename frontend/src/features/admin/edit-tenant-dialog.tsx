import { useState } from 'react'

import type { AdminTenantResponse } from '@/api/types'
import { ApiError } from '@/api/types'
import { MaterialIcon } from '@/components/icons/material-icon'
import { ConsolePortal } from '@/components/layout/console-portal'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useUpdateAdminTenant } from '@/features/admin/use-admin-queries'

interface EditTenantDialogProps {
  open: boolean
  tenant: AdminTenantResponse | null
  onOpenChange: (open: boolean) => void
}

function EditTenantForm({
  tenant,
  onClose,
}: {
  tenant: AdminTenantResponse
  onClose: () => void
}) {
  const updateTenant = useUpdateAdminTenant(tenant.id)
  const [name, setName] = useState(tenant.name)
  const [domain, setDomain] = useState(tenant.domain ?? '')
  const [error, setError] = useState<string | null>(null)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError(null)

    try {
      await updateTenant.mutateAsync({
        name: name.trim(),
        domain: domain.trim(),
      })
      onClose()
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Could not update tenant.')
    }
  }

  return (
    <>
      {error ? (
        <Alert variant="destructive" className="mb-4">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      ) : null}

      <form onSubmit={(e) => void handleSubmit(e)} className="space-y-4">
        <div className="space-y-2">
          <Label htmlFor="edit-tenant-name">Organization name</Label>
          <Input
            id="edit-tenant-name"
            required
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
        </div>

        <div className="space-y-2">
          <Label htmlFor="edit-tenant-domain">Domain</Label>
          <Input
            id="edit-tenant-domain"
            value={domain}
            onChange={(e) => setDomain(e.target.value)}
            placeholder="acme"
          />
        </div>

        <div className="flex justify-end gap-3 pt-2">
          <Button type="button" variant="outline" onClick={onClose}>
            Cancel
          </Button>
          <Button type="submit" disabled={updateTenant.isPending || !name.trim()}>
            {updateTenant.isPending ? 'Saving…' : 'Save changes'}
          </Button>
        </div>
      </form>
    </>
  )
}

export function EditTenantDialog({ open, tenant, onOpenChange }: EditTenantDialogProps) {
  if (!open || !tenant) {
    return null
  }

  function handleClose() {
    onOpenChange(false)
  }

  return (
    <ConsolePortal>
    <div className="console-modal-scrim fixed inset-0 z-50 flex items-center justify-center p-4">
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby="edit-tenant-title"
        className="console-modal-panel w-full max-w-md rounded-xl p-6"
      >
        <div className="mb-6 flex items-start justify-between">
          <h2 id="edit-tenant-title" className="font-headline text-xl font-bold text-on-surface">
            Edit tenant
          </h2>
          <button
            type="button"
            onClick={handleClose}
            className="rounded-lg p-1 text-on-surface-variant hover:bg-surface-container"
            aria-label="Close"
          >
            <MaterialIcon name="close" />
          </button>
        </div>

        <EditTenantForm key={tenant.id} tenant={tenant} onClose={handleClose} />
      </div>
    </div>
    </ConsolePortal>
  )
}
