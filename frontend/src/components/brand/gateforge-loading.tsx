import { GateForgeIconMark } from '@/components/brand/gateforge-icon-mark'
import { cn } from '@/lib/utils'

interface GateForgeLoadingProps {
  label?: string
  className?: string
}

export function GateForgeLoading({
  label = 'Loading…',
  className,
}: Readonly<GateForgeLoadingProps>) {
  return (
    <div
      className={cn('flex min-h-[100dvh] items-center justify-center bg-background', className)}
      role="status"
      aria-live="polite"
      aria-label={label}
    >
      <div className="h-24 w-24 p-[18%]">
        <GateForgeIconMark />
      </div>
    </div>
  )
}
