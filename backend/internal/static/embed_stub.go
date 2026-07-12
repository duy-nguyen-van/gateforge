//go:build !embedfrontend

package static

import "embed"

// Dist is populated when built with -tags embedfrontend.
var Dist embed.FS
