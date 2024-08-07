package main

import (
	"os"
)

var redisConnStr = defaultString(os.Getenv("REDIS_CONN_STR"), "redis://localhost:6379")
