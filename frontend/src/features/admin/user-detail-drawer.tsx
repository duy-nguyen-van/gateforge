import { ApiError } from '@/api/types'
import { DefaultAvatar } from '@/components/avatars/default-avatar'
import { MaterialIcon } from '@/components/icons/material-icon'
import { ConsolePortal } from '@/components/layout/console-portal'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { displayUserName, formatUserStatus } from '@/features/admin/admin-utils'
import { ConsoleErrorState, ConsoleLoadingState } from '@/features/admin/console-state'
import {
  useAdminUser,
  useDisableAdminUser,
  useForceLogoutAdminUser,
  useResetAdminUserMFA,
  useResetAdminUserPasskeys,
} from '@/features/admin/use-admin-queries'
import { useAuth } from '@/hooks/use-auth'
import { useState } from 'react'

type ConfirmAction = 'disable' | 'force-logout' | 'reset-passkey' | 'reset-mfa' | null

interface UserDetailDrawerProps {
  userId: string | null
  onClose: () => void
}

export function UserDetailDrawer({ userId, onClose }: UserDetailDrawerProps) {
  const { user: currentUser } = useAuth()
  const userQuery = useAdminUser(userId)
  const disableUser = useDisableAdminUser()
  const forceLogout = useForceLogoutAdminUser()
  const resetPasskeys = useResetAdminUserPasskeys()
  const resetMFA = useResetAdminUserMFA()
  const [confirmAction, setConfirmAction] = useState<ConfirmAction>(null)

  if (!userId) {
    return null
  }

  const user = userQuery.data?.data
  const isSelf = currentUser?.id === userId
  const isDisabled = user?.status === 'disabled'

  const pendingMutation = disableUser.isPending || forceLogout.isPending || resetPasskeys.isPending || resetMFA.isPending
  const mutationError = disableUser.error ?? forceLogout.error ?? resetPasskeys.error ?? resetMFA.error
  const errorMessage =
    mutationError instanceof ApiError
      ? mutationError.message
      : mutationError
        ? 'Action failed.'
        : null

  async function runConfirmedAction() {
    if (!userId || !confirmAction) return
    try {
      if (confirmAction === 'disable') {
        await disableUser.mutateAsync(userId)
      } else if (confirmAction === 'force-logout') {
        await forceLogout.mutateAsync(userId)
      } else if (confirmAction === 'reset-passkey') {
        await resetPasskeys.mutateAsync(userId)
      } else if (confirmAction === 'reset-mfa') {
        await resetMFA.mutateAsync(userId)
      }
      setConfirmAction(null)
      void userQuery.refetch()
    } catch {
      // error surfaced via mutation state
    }
  }

  const name = user ? displayUserName(user.first_name, user.last_name, user.email) : ''
  const disableBlocked = isSelf || isDisabled

  return (
    <ConsolePortal>
    <aside className="console-drawer-panel fixed right-0 top-0 z-50 flex h-full w-full max-w-md flex-col border-l border-outline-variant/20">
        <div className="flex items-center justify-between border-b border-surface-container px-6 py-4">
          <h2 className="font-headline text-lg font-bold text-on-surface">User details</h2>
          <button
            type="button"
            onClick={onClose}
            className="rounded-lg p-1 text-on-surface-variant hover:bg-surface-container"
            aria-label="Close"
          >
            <MaterialIcon name="close" />
          </button>
        </div>

        <div className="flex-1 overflow-y-auto px-6 py-6">
          {userQuery.isLoading ? (
            <ConsoleLoadingState label="Loading user…" />
          ) : userQuery.isError ? (
            <ConsoleErrorState message="Could not load user details." />
          ) : user ? (
            <div className="space-y-6">
              <div className="flex items-center gap-4">
                <DefaultAvatar seed={user.email} name={name} size="lg" title={name} />
                <div>
                  <p className="font-headline text-xl font-bold">{name}</p>
                  <p className="text-sm text-on-surface-variant">{user.email}</p>
                  <div className="mt-2 flex flex-wrap gap-2">
                    <span
                      className={`rounded-full px-2 py-0.5 text-xs font-bold ${user.status === 'active' ? 'bg-green-50 text-green-700' : 'bg-error-container/30 text-error'
                        }`}
                    >
                      {formatUserStatus(user.status)}
                    </span>
                    {user.is_platform_admin ? (
                      <span className="rounded-full bg-primary-container px-2 py-0.5 text-xs font-bold text-on-primary-container">
                        Platform admin
                      </span>
                    ) : null}
                  </div>
                </div>
              </div>

              <dl className="grid grid-cols-2 gap-3 text-sm">
                <div className="rounded-lg bg-surface-container-low p-3">
                  <dt className="text-xs font-bold uppercase text-on-surface-variant">MFA</dt>
                  <dd className="mt-1 font-semibold">{user.mfa_enabled ? 'Enabled' : 'Off'}</dd>
                </div>
                <div className="rounded-lg bg-surface-container-low p-3">
                  <dt className="text-xs font-bold uppercase text-on-surface-variant">Passkeys</dt>
                  <dd className="mt-1 font-semibold">{user.passkey_count}</dd>
                </div>
                <div className="rounded-lg bg-surface-container-low p-3">
                  <dt className="text-xs font-bold uppercase text-on-surface-variant">Active sessions</dt>
                  <dd className="mt-1 font-semibold">{user.active_sessions}</dd>
                </div>
                <div className="rounded-lg bg-surface-container-low p-3">
                  <dt className="text-xs font-bold uppercase text-on-surface-variant">Joined</dt>
                  <dd className="mt-1 font-semibold">{new Date(user.created_at).toLocaleDateString()}</dd>
                </div>
              </dl>

              {user.memberships.length > 0 ? (
                <div>
                  <h3 className="mb-2 text-xs font-bold uppercase tracking-wider text-on-surface-variant">
                    Tenant memberships
                  </h3>
                  <ul className="space-y-2">
                    {user.memberships.map((m) => (
                      <li key={m.tenant_id} className="rounded-lg bg-surface-container-low px-3 py-2 text-sm">
                        <p className="font-semibold">{m.tenant_name || m.tenant_id.slice(0, 8)}</p>
                        <p className="text-xs text-on-surface-variant">
                          {m.role} · {m.status}
                        </p>
                      </li>
                    ))}
                  </ul>
                </div>
              ) : null}

              {errorMessage ? (
                <Alert variant="destructive">
                  <AlertDescription>{errorMessage}</AlertDescription>
                </Alert>
              ) : null}

              {confirmAction ? (
                <div className="rounded-xl border border-surface-container bg-surface-container-low p-4">
                  <p className="text-sm text-on-surface-variant">
                    {confirmAction === 'disable' &&
                      'Disable this account? The user will be signed out and cannot sign in until re-enabled manually.'}
                    {confirmAction === 'force-logout' &&
                      'Force logout? All sessions and refresh tokens for this user will be revoked.'}
                    {confirmAction === 'reset-passkey' &&
                      'Remove all passkeys? The user must register new passkeys to sign in with WebAuthn.'}
                    {confirmAction === 'reset-mfa' &&
                      'Reset MFA? TOTP and recovery codes will be removed. The user must set up MFA again from their security settings.'}
                  </p>
                  <div className="mt-4 flex gap-2">
                    <Button type="button" variant="outline" size="sm" onClick={() => setConfirmAction(null)}>
                      Cancel
                    </Button>
                    <Button
                      type="button"
                      variant={confirmAction === 'disable' ? 'destructive' : 'default'}
                      size="sm"
                      disabled={pendingMutation}
                      onClick={() => void runConfirmedAction()}
                    >
                      Confirm
                    </Button>
                  </div>
                </div>
              ) : (
                <div className="flex flex-col gap-2">
                  <Button
                    type="button"
                    variant="destructive"
                    disabled={disableBlocked}
                    title={
                      isSelf
                        ? 'You cannot disable your own account'
                        : isDisabled
                          ? 'User is already disabled'
                          : undefined
                    }
                    onClick={() => setConfirmAction('disable')}
                  >
                    Disable account
                  </Button>
                  <Button type="button" variant="outline" onClick={() => setConfirmAction('force-logout')}>
                    Force logout
                  </Button>
                  <Button type="button" variant="outline" onClick={() => setConfirmAction('reset-passkey')}>
                    Reset passkeys
                  </Button>
                  <Button type="button" variant="outline" onClick={() => setConfirmAction('reset-mfa')}>
                    Reset MFA
                  </Button>
                </div>
              )}
            </div>
          ) : null}
        </div>
      </aside>
    </ConsolePortal>
  )
}
