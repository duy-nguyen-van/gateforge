package dtos

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPaginateSlice(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}

	page, pageable := PaginateSlice(items, &PageableRequest{Page: 2, PageSize: 2})
	require.Equal(t, []int{3, 4}, page)
	require.Equal(t, int64(5), pageable.Total)
	require.Equal(t, 2, pageable.Page)
	require.Equal(t, 2, pageable.PageSize)
}

func TestPaginateSlice_EmptyPage(t *testing.T) {
	items := []int{1, 2}

	page, pageable := PaginateSlice(items, &PageableRequest{Page: 5, PageSize: 2})
	require.Empty(t, page)
	require.Equal(t, int64(2), pageable.Total)
}

func TestPaginateSlice_EmptyItems(t *testing.T) {
	page, pageable := PaginateSlice([]int{}, nil)
	require.Empty(t, page)
	require.Equal(t, int64(0), pageable.Total)
}

func TestPaginateSlice_NilRequestUsesDefaults(t *testing.T) {
	page, pageable := PaginateSlice([]int{1, 2, 3}, nil)
	require.Len(t, page, 3)
	require.Equal(t, int64(3), pageable.Total)
}

func TestClampPageSize_NilRequest(t *testing.T) {
	require.NotPanics(t, func() { ClampPageSize(nil, 100) })
}

func TestClampPageSize_ZeroMaxSize(t *testing.T) {
	pr := &PageableRequest{Page: 0, PageSize: 500}
	ClampPageSize(pr, 0)
	require.Equal(t, 500, pr.PageSize)
}

func TestClampPageSize(t *testing.T) {
	pr := &PageableRequest{Page: 0, PageSize: 500}
	ClampPageSize(pr, 100)
	require.Equal(t, 100, pr.PageSize)
	require.Equal(t, 1, pr.Page)

	pr = &PageableRequest{Page: 2, PageSize: 0}
	ClampPageSize(pr, 50)
	require.Equal(t, 10, pr.PageSize)
}
