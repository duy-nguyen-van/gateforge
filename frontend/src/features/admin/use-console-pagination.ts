import { useCallback, useState } from 'react'

import { DEFAULT_CONSOLE_PAGE_SIZE } from '@/features/admin/admin-utils'

export function useConsolePagination(pageSize = DEFAULT_CONSOLE_PAGE_SIZE) {
  const [page, setPage] = useState(1)

  const resetPage = useCallback(() => {
    setPage(1)
  }, [])

  return {
    page,
    setPage,
    pageSize,
    resetPage,
    queryParams: { page, page_size: pageSize },
  }
}
