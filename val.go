package deno

import (
	"encoding/json"
	"math"
	"bytes"
)

func fromJSValue(val []byte) (goVal interface{}, err error) {
	if bytes.HasPrefix(val, nan) {
		goVal = math.NaN()
		return
	}
	if bytes.HasPrefix(val, undefined) {
		return
	}
	if e := json.Unmarshal(val, &goVal); e != nil {
		v := jsonNameRE.ReplaceAll(val, jsonNameRepl)
		if e = json.Unmarshal(v, &goVal); e != nil {
			err = e
		}
	}
	return
}
