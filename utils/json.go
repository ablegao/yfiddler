package utils

import "encoding/json"

func Json(val interface{}) string {

	b, err := json.Marshal(val)
	if err != nil {
		panic(err)
	}
	return string(b)
}
