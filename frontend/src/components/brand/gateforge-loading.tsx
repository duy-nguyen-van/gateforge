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
      className={cn('flex min-h-screen items-center justify-center bg-transparent', className)}
      role="status"
      aria-live="polite"
      aria-label={label}
    >
      <div className="gateforge-loader relative h-14 w-14">
        <span className="gateforge-loader-ring" aria-hidden />
        <div className="gateforge-loader-mark relative z-10 h-full w-full">
          <GateForgeIconMark />
        </div>
      </div>
    </div>
  )
}
