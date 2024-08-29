package rl

import (
	"go.uber.org/ratelimit"
	"time"
)

var Limiter = ratelimit.New(5, ratelimit.Per(time.Minute))
