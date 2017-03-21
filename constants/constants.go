package constants

import (
	"time"
)

const (
	Production        = false
	BuildTest         = true
	StaticRoot        = "www/dist"
	StaticTestingRoot = "www"
	StaticCache       = false
	PackageExt        = "pkg.tar.xz"
	BuildImage        = "builder"
	RetryDelay        = 3 * time.Second
)
