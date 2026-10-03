import type { ReactNode } from 'react'

export function ConsolePageHeader({
  title,
  description,
  actions,
}: {
  title: string
  description?: string
  actions?: ReactNode
}) {
  return (
    <header className="mb-8 flex flex-col justify-between gap-4 sm:flex-row sm:items-end">
      <div>
        <h1 className="font-headline text-2xl font-semibold tracking-tight text-on-surface">{title}</h1>
        {description ? <p className="mt-1 max-w-2xl text-sm text-on-surface-variant">{description}</p> : null}
      </div>
      {actions ? <div className="flex shrink-0 flex-wrap gap-3">{actions}</div> : null}
    </header>
  )
}
