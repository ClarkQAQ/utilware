// Use modernc.org/sqlite on all supported platforms unless Cgo driver
// was explicitly requested.
//
// See also https://pkg.go.dev/modernc.org/sqlite#hdr-Supported_platforms_and_architectures

package sqliteshim

import "utilware/sqlite"

const (
	hasDriver  = true
	driverName = "sqlite"
)

var shimDriver = &sqlite.Driver{}
