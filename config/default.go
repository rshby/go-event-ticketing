package config

import "time"

var (
	DefaultDbConnMaxLifetime   = 20 * time.Minute
	DbDefaultDbConnMaxIdletime = 10 * time.Minute
	DefaultRedisTTL            = 30 * time.Minute

	// http server
	DefaultHttpServerReadHeaderTimeout = 3 * time.Second
	DefaultHttpServerReadTimeout       = 5 * time.Second
	DefaultHttpServerWriteTimeout      = 8 * time.Second
	DefaultHttpServerIdleTimeout       = 30 * time.Second
)
