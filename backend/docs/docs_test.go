package docs

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/swaggo/swag"
)

func TestSwaggerInfoRegisteredOnInit(t *testing.T) {
	spec, err := swag.ReadDoc(SwaggerInfo.InstanceName())
	require.NoError(t, err)
	require.Contains(t, spec, `"basePath": "/api/v1"`)
	require.Equal(t, "swagger", SwaggerInfo.InstanceName())
}
