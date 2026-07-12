export const DEFAULT_AVATAR_COUNT = 6

/** Stitch user-management / audit-logs M3 avatar palette */
export const DEFAULT_AVATAR_VARIANTS = [
  { bg: 'bg-primary-container', text: 'text-on-primary-container' },
  { bg: 'bg-secondary-container', text: 'text-on-secondary-container' },
  { bg: 'bg-tertiary-container', text: 'text-on-tertiary-container' },
  { bg: 'bg-surface-variant', text: 'text-on-surface-variant' },
  { bg: 'bg-primary-container/30', text: 'text-on-primary-container' },
  { bg: 'bg-surface-container-highest', text: 'text-primary' },
] as const

export function getDefaultAvatarIndex(seed: string): number {
  if (!seed) {
    return 0
  }

  let hash = 2166136261
  for (let i = 0; i < seed.length; i++) {
    hash ^= seed.charCodeAt(i)
    hash = Math.imul(hash, 16777619)
  }

  return Math.abs(hash) % DEFAULT_AVATAR_COUNT
}

export function getAvatarInitials(seed: string, name?: string): string {
  if (name?.trim()) {
    const parts = name.trim().split(/\s+/).filter(Boolean)
    if (parts.length >= 2) {
      return `${parts[0][0] ?? ''}${parts[parts.length - 1][0] ?? ''}`.toUpperCase()
    }
    if (parts[0].length >= 2) {
      return parts[0].slice(0, 2).toUpperCase()
    }
    return (parts[0][0] ?? '?').toUpperCase()
  }

  const local = seed.split('@')[0] ?? seed
  const tokens = local.split(/[._-]+/).filter(Boolean)
  if (tokens.length >= 2) {
    return `${tokens[0][0] ?? ''}${tokens[1][0] ?? ''}`.toUpperCase()
  }
  if (local.length >= 2) {
    return local.slice(0, 2).toUpperCase()
  }
  return (local[0] ?? '?').toUpperCase()
}
