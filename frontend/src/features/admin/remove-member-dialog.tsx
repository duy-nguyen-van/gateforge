import { XIcon } from 'lucide-react'

import { ApiError } from '@/api/types'
import { ConsoleDialog } from '@/components/layout/console-dialog'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { useRemoveTenantMember } from '@/features/admin/use-admin-queries'

interface RemoveMemberDialogProps {
  open: boolean
  email: string
  tenantId: string
  userId: string
  onOpenChange: (open: boolean) => void
}

export function RemoveMemberDialog({
  open,
  email,
  tenantId,
  userId,
  onOpenChange,
}: RemoveMemberDialogProps) {
  const removeMember = useRemoveTenantMember()

  async function handleConfirm(close: () => void) {
    try {
      await removeMember.mutateAsync({ tenantId, userId })
      close()
    } catch {
      // error shown via mutation state below
    }
  }

  const errorMessage =
    removeMember.error instanceof ApiError
      ? removeMember.error.message
      : removeMember.error
        ? 'Could not remove member.'
        : null

  return (
    <ConsoleDialog
      open={open}
      onOpenChange={onOpenChange}
      labelledBy="remove-member-title"
      panelClassName="w-full max-w-md"
    >
      {({ close }) => (
        <>
          <div className="mb-4 flex items-start justify-between">
            <h2 id="remove-member-title" className="font-headline text-xl font-bold text-on-surface">
              Remove from tenant
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
            Remove <span className="font-semibold text-on-surface">{email}</span> from tenant{' '}
            <span className="font-mono text-xs">{tenantId.slice(0, 8)}…</span>? This action cannot be undone.
          </p>

          <div className="mt-6 flex justify-end gap-3">
            <Button type="button" variant="outline" onClick={close}>
              Cancel
            </Button>
            <Button
              type="button"
              variant="destructive"
              disabled={removeMember.isPending}
              onClick={() => void handleConfirm(close)}
            >
              {removeMember.isPending ? 'Removing…' : 'Remove'}
            </Button>
          </div>
        </>
      )}
    </ConsoleDialog>
  )
}
