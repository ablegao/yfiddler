package main

import (
	"yfiddler/hooks"
)

func main() {
}

type cbcDecrypter struct {
}

func (self *cbcDecrypter) Gen(args ...string) string {
	if len(args) == 2 {
		b, err := hooks.CBCDecrypter([]byte(args[0]), []byte(args[1]))
		if err == nil {
			return string(b)
		}
	}
	return ""
}

// exported as symbol named "Cbcdecrypter"
var Cbcdecrypter cbcDecrypter
