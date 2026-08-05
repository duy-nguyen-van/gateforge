import { useState, type AnimationEvent, type ReactNode } from 'react'

import { ConsolePortal } from '@/components/layout/console-portal'
import { cn } from '@/lib/utils'

function prefersReducedMotion() {
  return typeof window !== 'undefined' && window.matchMedia('(prefers-reduced-motion: reduce)').matches
}

export type ConsoleDialogApi = {
  close: () => void
}

interface ConsoleDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  labelledBy: string
  panelClassName?: string
  /** Defaults to true. Set false for one-time secret / irreversible confirmations. */
  closeOnOutsideClick?: boolean
  onExitComplete?: () => void
  children: ReactNode | ((api: ConsoleDialogApi) => ReactNode)
}

export function ConsoleDialog({
  open,
  onOpenChange,
  labelledBy,
  panelClassName,
  closeOnOutsideClick = true,
  onExitComplete,
  children,
}: ConsoleDialogProps) {
  const [leaving, setLeaving] = useState(false)
  const [prevOpen, setPrevOpen] = useState(open)

  if (open !== prevOpen) {
    setPrevOpen(open)
    if (open) {
      setLeaving(false)
    } else if (prefersReducedMotion()) {
      const done = onExitComplete
      queueMicrotask(() => done?.())
    } else {
      setLeaving(true)
    }
  }

  const mounted = open || leaving

  function close() {
    if (!open || leaving) return
    onOpenChange(false)
  }

  function handlePanelAnimationEnd(event: AnimationEvent<HTMLDivElement>) {
    if (!leaving || event.animationName !== 'console-modal-out') return
    setLeaving(false)
    onExitComplete?.()
  }

  if (!mounted) {
    return null
  }

  const content = typeof children === 'function' ? children({ close }) : children

  return (
    <ConsolePortal>
      <div
        className={cn(
          'console-modal-scrim fixed inset-0 z-50 flex items-center justify-center p-4',
          leaving && 'is-leaving',
        )}
        onClick={closeOnOutsideClick ? close : undefined}
      >
        <div
          role="dialog"
          aria-modal="true"
          aria-labelledby={labelledBy}
          className={cn('console-modal-panel rounded-xl p-6', panelClassName, leaving && 'is-leaving')}
          onClick={(event) => event.stopPropagation()}
          onAnimationEnd={handlePanelAnimationEnd}
        >
          {content}
        </div>
      </div>
    </ConsolePortal>
  )
}
