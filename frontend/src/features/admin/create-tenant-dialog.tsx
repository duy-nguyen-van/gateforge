import { useState } from 'react'

import { XIcon } from 'lucide-react'

import { ApiError } from '@/api/types'
import { ConsoleDialog } from '@/components/layout/console-dialog'
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

  function resetForm() {
    setName('')
    setDomain('')
    setError(null)
  }

  async function handleSubmit(e: React.FormEvent, close: () => void) {
    e.preventDefault()
    setError(null)

    try {
      const result = await createTenant.mutateAsync({
        name: name.trim(),
        domain: domain.trim() || undefined,
      })
      const tenantId = result.data?.id
      close()
      if (tenantId && onCreated) {
        onCreated(tenantId)
      }
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Could not create tenant.')
    }
  }

  return (
    <ConsoleDialog
      open={open}
      onOpenChange={onOpenChange}
      labelledBy="create-tenant-title"
      panelClassName="w-full max-w-md"
      onExitComplete={resetForm}
    >
      {({ close }) => (
        <>
          <div className="mb-6 flex items-start justify-between">
            <div>
              <h2 id="create-tenant-title" className="font-headline text-xl font-bold text-on-surface">
                New tenant
              </h2>
              <p className="mt-1 text-sm text-on-surface-variant">Create an organization boundary.</p>
            </div>
            <button
              type="button"
              onClick={close}
              className="rounded-lg p-1 text-on-surface-variant hover:bg-surface-container"
              aria-label="Close"
            >
              <XIcon className="h-5 w-5" aria-hidden />
            </button>
          </div>

          {error ? (
            <Alert variant="destructive" className="mb-4">
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          ) : null}

          <form onSubmit={(e) => void handleSubmit(e, close)} className="space-y-4">
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
              <Button type="button" variant="outline" onClick={close}>
                Cancel
              </Button>
              <Button type="submit" disabled={createTenant.isPending || !name.trim()}>
                {createTenant.isPending ? 'Creating…' : 'Create tenant'}
              </Button>
            </div>
          </form>
        </>
      )}
    </ConsoleDialog>
  )
}
