import { ApiError } from '@/api/types'
import { MaterialIcon } from '@/components/icons/material-icon'
import { ConsolePortal } from '@/components/layout/console-portal'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { useDeleteAdminClient } from '@/features/admin/use-admin-queries'

interface DeleteClientDialogProps {
  open: boolean
  clientId: string
  clientName: string
  onOpenChange: (open: boolean) => void
  onDeleted?: () => void
}

export function DeleteClientDialog({
  open,
  clientId,
  clientName,
  onOpenChange,
  onDeleted,
}: DeleteClientDialogProps) {
  const deleteClient = useDeleteAdminClient()

  async function handleConfirm() {
    try {
      await deleteClient.mutateAsync(clientId)
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
    deleteClient.error instanceof ApiError
      ? deleteClient.error.message
      : deleteClient.error
        ? 'Could not delete client.'
        : null

  return (
    <ConsolePortal>
    <div className="console-modal-scrim fixed inset-0 z-50 flex items-center justify-center p-4">
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby="delete-client-title"
        className="console-modal-panel w-full max-w-md rounded-xl p-6"
      >
        <div className="mb-4 flex items-start justify-between">
          <h2 id="delete-client-title" className="font-headline text-xl font-bold text-on-surface">
            Delete client
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
          Permanently delete <span className="font-semibold text-on-surface">{clientName || clientId}</span>?
          Active OAuth sessions and tokens for this client will stop working.
        </p>

        <div className="mt-6 flex justify-end gap-3">
          <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button
            type="button"
            variant="destructive"
            disabled={deleteClient.isPending}
            onClick={() => void handleConfirm()}
          >
            {deleteClient.isPending ? 'Deleting…' : 'Delete client'}
          </Button>
        </div>
      </div>
    </div>
    </ConsolePortal>
  )
}
