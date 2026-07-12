package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertAccented(t *testing.T) {
	require.Equal(t, "nguyen", ConvertAccented("Nguyễn"))
	require.Equal(t, "dac biet", ConvertAccented("Đặc Biệt"))
	require.Equal(t, "hello", ConvertAccented("HELLO"))
}
