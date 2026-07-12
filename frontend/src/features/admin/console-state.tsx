import type { ReactNode } from 'react'

import { MaterialIcon } from '@/components/icons/material-icon'
import { Alert, AlertDescription } from '@/components/ui/alert'

export function ConsoleLoadingState({ label = 'Loading data…' }: { label?: string }) {
  return (
    <div className="flex min-h-[240px] items-center justify-center rounded-xl bg-surface-container-lowest ghost-border">
      <div className="flex items-center gap-3 text-on-surface-variant">
        <MaterialIcon name="progress_activity" className="animate-spin text-primary" />
        <span className="text-sm font-medium">{label}</span>
      </div>
    </div>
  )
}

export function ConsoleErrorState({ message, action }: { message: string; action?: ReactNode }) {
  return (
    <Alert variant="destructive" className="rounded-xl">
      <AlertDescription className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <span>{message}</span>
        {action}
      </AlertDescription>
    </Alert>
  )
}

export function ConsoleEmptyState({ title, description }: { title: string; description: string }) {
  return (
    <div className="flex min-h-[240px] flex-col items-center justify-center rounded-xl bg-surface-container-lowest p-8 text-center ghost-border">
      <MaterialIcon name="inventory_2" className="mb-3 text-3xl text-on-surface-variant" />
      <h3 className="font-headline text-lg font-bold text-on-surface">{title}</h3>
      <p className="mt-1 max-w-md text-sm text-on-surface-variant">{description}</p>
    </div>
  )
}
