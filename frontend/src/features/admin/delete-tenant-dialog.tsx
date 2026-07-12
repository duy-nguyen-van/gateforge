import { ApiError } from '@/api/types'
import { MaterialIcon } from '@/components/icons/material-icon'
import { ConsolePortal } from '@/components/layout/console-portal'
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

  async function handleConfirm() {
    try {
      await deleteTenant.mutateAsync(tenantId)
      onOpenChange(false)
      onDeleted?.()
    } catch {
      // error shown via mutation state below
    }
  }

  if (!open) {
    return null
  }

  const errorMessage =
    deleteTenant.error instanceof ApiError
      ? deleteTenant.error.message
      : deleteTenant.error
        ? 'Could not delete tenant.'
        : null

  return (
    <ConsolePortal>
    <div className="console-modal-scrim fixed inset-0 z-50 flex items-center justify-center p-4">
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby="delete-tenant-title"
        className="console-modal-panel w-full max-w-md rounded-xl p-6"
      >
        <div className="mb-4 flex items-start justify-between">
          <h2 id="delete-tenant-title" className="font-headline text-xl font-bold text-on-surface">
            Delete tenant
          </h2>
          <button
            type="button"
            onClick={() => onOpenChange(false)}
            className="rounded-lg p-1 text-on-surface-variant hover:bg-surface-container"
            aria-label="Close"
          >
            <MaterialIcon name="close" />
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
          <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button
            type="button"
            variant="destructive"
            disabled={deleteTenant.isPending}
            onClick={() => void handleConfirm()}
          >
            {deleteTenant.isPending ? 'Deleting…' : 'Delete tenant'}
          </Button>
        </div>
      </div>
    </div>
    </ConsolePortal>
  )
}
