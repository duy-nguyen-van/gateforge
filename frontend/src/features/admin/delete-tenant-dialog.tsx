import { XIcon } from 'lucide-react'

import { ApiError } from '@/api/types'
import { ConsoleDialog } from '@/components/layout/console-dialog'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { useDeleteAdminTenant } from '@/features/admin/use-admin-queries'

interface DeleteTenantDialogProps {
  open: boolean
  tenantId: string
  tenantName: string
  onOpenChange: (open: boolean) => void
  onDeleted?: () => void
}

export function DeleteTenantDialog({
  open,
  tenantId,
  tenantName,
  onOpenChange,
  onDeleted,
}: DeleteTenantDialogProps) {
  const deleteTenant = useDeleteAdminTenant()

  async function handleConfirm(close: () => void) {
    try {
      await deleteTenant.mutateAsync(tenantId)
      close()
      onDeleted?.()
    } catch {
      // error shown via mutation state below
    }
  }

  const errorMessage =
    deleteTenant.error instanceof ApiError
      ? deleteTenant.error.message
      : deleteTenant.error
        ? 'Could not delete tenant.'
        : null

  return (
    <ConsoleDialog
      open={open}
      onOpenChange={onOpenChange}
      labelledBy="delete-tenant-title"
      panelClassName="w-full max-w-md"
    >
      {({ close }) => (
        <>
          <div className="mb-4 flex items-start justify-between">
            <h2 id="delete-tenant-title" className="font-headline text-xl font-bold text-on-surface">
              Delete tenant
            </h2>
            <button
              type="button"
              onClick={close}
              className="rounded-lg p-1 text-on-surface-variant hover:bg-surface-container"
              aria-label="Close"
            >
              <XIcon className="h-5 w-5" aria-hidden />
            </button>
          </div>

          {errorMessage ? (
            <Alert variant="destructive" className="mb-4">
              <AlertDescription>{errorMessage}</AlertDescription>
            </Alert>
          ) : null}

          <p className="text-sm text-on-surface-variant">
            Permanently delete <span className="font-semibold text-on-surface">{tenantName || tenantId}</span>?
            Members will lose access to this tenant.
          </p>

          <div className="mt-6 flex justify-end gap-3">
            <Button type="button" variant="outline" onClick={close}>
              Cancel
            </Button>
            <Button
              type="button"
              variant="destructive"
              disabled={deleteTenant.isPending}
              onClick={() => void handleConfirm(close)}
            >
              {deleteTenant.isPending ? 'Deleting…' : 'Delete tenant'}
            </Button>
          </div>
        </>
      )}
    </ConsoleDialog>
  )
}
