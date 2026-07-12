package dtos

import "github.com/gateforge-iam/gateforge-iam/internal/constants"

// ClampPageSize caps page size to prevent abuse.
func ClampPageSize(pr *PageableRequest, maxSize int) {
	if pr == nil || maxSize <= 0 {
		return
	}
	if pr.PageSize <= 0 {
		pr.PageSize = constants.DefaultPageSize
	}
	if pr.PageSize > maxSize {
		pr.PageSize = maxSize
	}
	if pr.Page <= 0 {
		pr.Page = constants.DefaultPage
	}
}

// PaginateSlice returns a page of items and pageable metadata for in-memory lists.
func PaginateSlice[T any](items []T, pr *PageableRequest) ([]T, *Pageable) {
	if pr == nil {
		pr = NewPageableRequest()
	}
	ClampPageSize(pr, constants.MaxPageSize)

	total := int64(len(items))
	if total == 0 {
		return []T{}, &Pageable{Page: pr.Page, PageSize: pr.PageSize, Total: 0}
	}

	offset := pr.GetOffset()
	if offset >= len(items) {
		return []T{}, &Pageable{Page: pr.Page, PageSize: pr.PageSize, Total: total}
	}

	end := offset + pr.GetLimit()
	if end > len(items) {
		end = len(items)
	}
	return items[offset:end], &Pageable{Page: pr.Page, PageSize: pr.PageSize, Total: total}
}
