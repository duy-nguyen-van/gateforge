//go:build embedfrontend

package static

import "embed"

//go:embed all:dist
var Dist embed.FS
