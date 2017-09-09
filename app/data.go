package app

import (
	"time"

	"github.com/elBroom/highloadCup/app/cache"
)

var Phase = 1
var LifeTime = 5 * time.Minute
var MemoryCache = cache.InitCache()
