import { useState } from 'react'

import { ApiError } from '@/api/types'
import { MaterialIcon } from '@/components/icons/material-icon'
import { ConsolePortal } from '@/components/layout/console-portal'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useCreateAdminTenant } from '@/features/admin/use-admin-queries'

interface CreateTenantDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onCreated?: (tenantId: string) => void
}

export function CreateTenantDialog({ open, onOpenChange, onCreated }: CreateTenantDialogProps) {
  const createTenant = useCreateAdminTenant()
  const [name, setName] = useState('')
  const [domain, setDomain] = useState('')
  const [error, setError] = useState<string | null>(null)

  function handleClose() {
    onOpenChange(false)
    setName('')
    setDomain('')
    setError(null)
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError(null)

    try {
      const result = await createTenant.mutateAsync({
        name: name.trim(),
        domain: domain.trim() || undefined,
      })
      const tenantId = result.data?.id
      handleClose()
      if (tenantId && onCreated) {
        onCreated(tenantId)
      }
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Could not create tenant.')
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
        aria-labelledby="create-tenant-title"
        className="console-modal-panel w-full max-w-md rounded-xl p-6"
      >
        <div className="mb-6 flex items-start justify-between">
          <div>
            <h2 id="create-tenant-title" className="font-headline text-xl font-bold text-on-surface">
              New tenant
            </h2>
            <p className="mt-1 text-sm text-on-surface-variant">Create an organization boundary.</p>
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

        <form onSubmit={(e) => void handleSubmit(e)} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="create-tenant-name">Organization name</Label>
            <Input
              id="create-tenant-name"
              required
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Acme Corp"
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="create-tenant-domain">Domain (optional)</Label>
            <Input
              id="create-tenant-domain"
              value={domain}
              onChange={(e) => setDomain(e.target.value)}
              placeholder="acme"
            />
            <p className="text-xs text-on-surface-variant">Used for subdomain-based tenant resolution.</p>
          </div>

          <div className="flex justify-end gap-3 pt-2">
            <Button type="button" variant="outline" onClick={handleClose}>
              Cancel
            </Button>
            <Button type="submit" disabled={createTenant.isPending || !name.trim()}>
              {createTenant.isPending ? 'Creating…' : 'Create tenant'}
            </Button>
          </div>
        </form>
      </div>
    </div>
    </ConsolePortal>
  )
}
