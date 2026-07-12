import { MaterialIcon } from '@/components/icons/material-icon'
import { paginationRange } from '@/features/admin/admin-utils'

interface ConsolePaginationProps {
  page: number
  pageSize: number
  total: number
  onPageChange: (page: number) => void
}

export function ConsolePagination({ page, pageSize, total, onPageChange }: ConsolePaginationProps) {
  if (total <= 0) {
    return null
  }

  const { start, end, totalPages } = paginationRange(page, pageSize, total)
  const canGoPrev = page > 1
  const canGoNext = page < totalPages

  return (
    <div className="flex items-center justify-between border-t border-surface-container px-6 py-4">
      <p className="text-xs font-medium text-on-surface-variant">
        Showing{' '}
        <span className="font-bold text-on-surface">
          {start}–{end}
        </span>{' '}
        of <span className="font-bold text-on-surface">{total.toLocaleString()}</span>
        {totalPages > 1 ? (
          <>
            {' '}
            · Page {page} of {totalPages}
          </>
        ) : null}
      </p>
      <div className="flex items-center gap-2">
        <button
          type="button"
          disabled={!canGoPrev}
          onClick={() => onPageChange(page - 1)}
          className="rounded-lg p-1.5 transition-colors hover:bg-surface-container disabled:opacity-30"
          aria-label="Previous page"
        >
          <MaterialIcon name="chevron_left" className="text-lg" />
        </button>
        <button
          type="button"
          disabled={!canGoNext}
          onClick={() => onPageChange(page + 1)}
          className="rounded-lg p-1.5 transition-colors hover:bg-surface-container disabled:opacity-30"
          aria-label="Next page"
        >
          <MaterialIcon name="chevron_right" className="text-lg" />
        </button>
      </div>
    </div>
  )
}
