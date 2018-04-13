package main

import (
	"encoding/base64"
)

func main() {
}

type base64Encode struct {
}

func (self *base64Encode) Gen(args ...string) string {
	return base64.StdEncoding.EncodeToString([]byte(args[0]))

}

// exported as symbol named "Base64encode"
var Base64encode base64Encode
