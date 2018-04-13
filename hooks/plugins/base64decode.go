package main

import "encoding/base64"

func main() {
}

type base64Decode struct {
}

func (self *base64Decode) Gen(args ...string) string {
	b, err := base64.StdEncoding.DecodeString(args[0])
	if err == nil {
		return string(b)
	}

	return ""
}

// exported as symbol named "Base64eecode"
var Base64decode base64Decode
