import { cn } from '@/lib/utils'

interface MaterialIconProps {
  name: string
  filled?: boolean
  className?: string
}

export function MaterialIcon({ name, filled, className }: MaterialIconProps) {
  return (
    <span
      className={cn(
        'material-symbols-outlined inline-flex items-center justify-center',
        filled && 'filled',
        className,
      )}
      aria-hidden
    >
      {name}
    </span>
  )
}
