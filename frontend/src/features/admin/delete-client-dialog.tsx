import { XIcon } from 'lucide-react'

import { ApiError } from '@/api/types'
import { ConsoleDialog } from '@/components/layout/console-dialog'
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

  async function handleConfirm(close: () => void) {
    try {
      await deleteClient.mutateAsync(clientId)
      close()
      onDeleted?.()
    } catch {
      // error shown via mutation state below
    }
  }

  const errorMessage =
    deleteClient.error instanceof ApiError
      ? deleteClient.error.message
      : deleteClient.error
        ? 'Could not delete client.'
        : null

  return (
    <ConsoleDialog
      open={open}
      onOpenChange={onOpenChange}
      labelledBy="delete-client-title"
      panelClassName="w-full max-w-md"
    >
      {({ close }) => (
        <>
          <div className="mb-4 flex items-start justify-between">
            <h2 id="delete-client-title" className="font-headline text-xl font-bold text-on-surface">
              Delete client
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
            Permanently delete <span className="font-semibold text-on-surface">{clientName || clientId}</span>?
            Active OAuth sessions and tokens for this client will stop working.
          </p>

          <div className="mt-6 flex justify-end gap-3">
            <Button type="button" variant="outline" onClick={close}>
              Cancel
            </Button>
            <Button
              type="button"
              variant="destructive"
              disabled={deleteClient.isPending}
              onClick={() => void handleConfirm(close)}
            >
              {deleteClient.isPending ? 'Deleting…' : 'Delete client'}
            </Button>
          </div>
        </>
      )}
    </ConsoleDialog>
  )
}
