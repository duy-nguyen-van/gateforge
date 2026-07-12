import gateforgeIcon from '@/assets/brand/gateforge-icon.webp'
import { cn } from '@/lib/utils'

interface GateForgeIconMarkProps {
  className?: string
  alt?: string
}

export function GateForgeIconMark({ className, alt = 'GateForge' }: Readonly<GateForgeIconMarkProps>) {
  return (
    <img
      src={gateforgeIcon}
      alt={alt}
      decoding="async"
      className={cn('h-full w-full object-contain', className)}
    />
  )
}
