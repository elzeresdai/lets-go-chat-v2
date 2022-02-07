package cache

import (
	"github.com/maxchagin/go-memorycache-example"
	"time"
)

var Cache *memorycache.Cache

func init() {
	cache := memorycache.New(5*time.Minute, 10*time.Minute)
	Cache = cache
}
