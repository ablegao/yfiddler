package main

import (
	"strings"
	"time"
)

type cacheBusting struct {
}

func (self *cacheBusting) Gen(args ...string) string {

	return strings.Replace(args[0], "{CACHE_TIME}", time.Now().Format(time.RFC3339Nano), -1)

}

// exported as symbol named "Cachebusting"
var Cachebusting cacheBusting
