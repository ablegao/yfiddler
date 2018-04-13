package main

import (
	"testing"
)

func TestBase64encode(t *testing.T) {
	m := mode{}
	t.Log(m.Gen("aaaaa.com?sss=ddd&"))
}
