import { ApiError } from '@/api/types'

export function isAdminMfaRequired(error: unknown): boolean {
  return error instanceof ApiError && error.errorCode === 'ADMIN_MFA_REQUIRED'
}
