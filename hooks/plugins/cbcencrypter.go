package main

import (
	"yfiddler/hooks"
)

func main() {
}

type cbcEncrypter struct {
}

func (self *cbcEncrypter) Gen(args ...string) string {
	if len(args) == 2 {
		b, err := hooks.CBCEncrypter([]byte(args[0]), []byte(args[1]))
		if err == nil {
			return string(b)
		}
	}
	return ""
}

// exported as symbol named "Cbcencrypter"
var Cbcencrypter cbcEncrypter
