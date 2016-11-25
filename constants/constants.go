package constants

import (
	"time"
)

const (
	StaticRoot  = "www"
	StaticLive  = true
	StaticCache = false
	PackageExt  = "pkg.tar.xz"
	BuildImage  = "builder"
	RetryDelay  = 3 * time.Second
)
