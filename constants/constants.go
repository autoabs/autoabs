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

// DB=mongodb://pritunl:s9T82kD8JVLm@db0.pritunl.net:27017/autoabs?ssl=true
