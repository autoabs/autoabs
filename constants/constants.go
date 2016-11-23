package constants

import (
	"time"
)

const (
	StaticLive  = true
	StaticRoot  = "www"
	StaticCache = false
	PackageExt  = "pkg.tar.xz"
	BuildImage  = "builder"
	RetryDelay  = 3 * time.Second
)
