package main

import (
	"sync/atomic"
	"time"
)

var initialSnowflake = time.Now().Unix()

func Snowflake() int64 {
	// in production we have to use more complex approach
	// and / or existing library for snowflakes or alternatives
	return atomic.AddInt64(&initialSnowflake, 1)
}
