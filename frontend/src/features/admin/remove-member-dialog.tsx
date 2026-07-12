import { ApiError } from '@/api/types'
import { MaterialIcon } from '@/components/icons/material-icon'
import { ConsolePortal } from '@/components/layout/console-portal'
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

  async function handleConfirm() {
    try {
      await removeMember.mutateAsync({ tenantId, userId })
      onOpenChange(false)
    } catch {
      // error shown via mutation state below
    }
  }

  if (!open) {
    return null
  }

  const errorMessage =
    removeMember.error instanceof ApiError
      ? removeMember.error.message
      : removeMember.error
        ? 'Could not remove member.'
        : null

  return (
    <ConsolePortal>
    <div className="console-modal-scrim fixed inset-0 z-50 flex items-center justify-center p-4">
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby="remove-member-title"
        className="console-modal-panel w-full max-w-md rounded-xl p-6"
      >
        <div className="mb-4 flex items-start justify-between">
          <h2 id="remove-member-title" className="font-headline text-xl font-bold text-on-surface">
            Remove from tenant
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
          Remove <span className="font-semibold text-on-surface">{email}</span> from tenant{' '}
          <span className="font-mono text-xs">{tenantId.slice(0, 8)}…</span>? This action cannot be undone.
        </p>

        <div className="mt-6 flex justify-end gap-3">
          <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button
            type="button"
            variant="destructive"
            disabled={removeMember.isPending}
            onClick={() => void handleConfirm()}
          >
            {removeMember.isPending ? 'Removing…' : 'Remove'}
          </Button>
        </div>
      </div>
    </div>
    </ConsolePortal>
  )
}
