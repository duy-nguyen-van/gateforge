import { cn } from '@/lib/utils'

import {
  DEFAULT_AVATAR_VARIANTS,
  getAvatarInitials,
  getDefaultAvatarIndex,
} from '@/components/avatars/default-avatar-utils'

const sizeClasses = {
  xs: 'h-6 w-6 text-[8px]',
  sm: 'h-8 w-8 text-[10px]',
  md: 'h-10 w-10 text-sm',
  lg: 'h-16 w-16 text-xl',
} as const

interface DefaultAvatarProps {
  seed: string
  name?: string
  variant?: number
  size?: keyof typeof sizeClasses
  className?: string
  title?: string
}

export function DefaultAvatar({ seed, name, variant, size = 'sm', className, title }: DefaultAvatarProps) {
  const index = variant ?? getDefaultAvatarIndex(seed)
  const style = DEFAULT_AVATAR_VARIANTS[index]
  const initials = getAvatarInitials(seed, name)

  return (
    <div
      role="img"
      aria-label={title ?? `${name ?? seed} avatar`}
      className={cn(
        'flex shrink-0 items-center justify-center rounded-circle font-bold',
        sizeClasses[size],
        style.bg,
        style.text,
        className,
      )}
    >
      {initials}
    </div>
  )
}

/** Stitch dashboard MFA card — two overlapping color chips, no initials */
export function MfaAvatarPreviewStack({ className }: { className?: string }) {
  return (
    <div className={cn('flex -space-x-2', className)}>
      <div className="h-6 w-6 rounded-circle border-2 border-surface-container-lowest bg-surface-container-highest" />
      <div className="h-6 w-6 rounded-circle border-2 border-surface-container-lowest bg-secondary-container" />
    </div>
  )
}
