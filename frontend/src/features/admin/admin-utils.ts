export const DEFAULT_CONSOLE_PAGE_SIZE = 25

export function paginationRange(page: number, pageSize: number, total: number) {
  const totalPages = Math.max(1, Math.ceil(total / pageSize))
  const safePage = Math.min(Math.max(page, 1), totalPages)
  const start = total === 0 ? 0 : (safePage - 1) * pageSize + 1
  const end = Math.min(safePage * pageSize, total)
  return { start, end, totalPages, safePage }
}

export function formatUserStatus(status: string): string {
  if (!status) return 'Unknown'
  return status.charAt(0).toUpperCase() + status.slice(1)
}

export function displayUserName(firstName: string, lastName: string, email: string): string {
  const full = `${firstName ?? ''} ${lastName ?? ''}`.trim()
  return full || email
}

export function formatAuditAction(action: string): string {
  return action.replaceAll('.', ' · ')
}

const auditResultStyles: Record<string, string> = {
  success: 'bg-success-container text-on-success-container',
  failure: 'bg-error-container text-error',
  denied: 'bg-warning-container text-on-warning-container',
}

export function auditResultBadgeClass(result: string): string {
  return auditResultStyles[result] ?? 'bg-surface-container-highest text-on-surface-variant'
}

export function formatAuditTimestamp(iso: string): string {
  const date = new Date(iso)
  if (Number.isNaN(date.getTime())) return iso
  return date.toLocaleString()
}
