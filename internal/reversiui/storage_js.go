//go:build js && wasm

package reversiui

import "syscall/js"

func storedBest(key string) int {
	v := js.Global().Get("localStorage").Call("getItem", key)
	if v.Type() != js.TypeString {
		return 0
	}
	value := 0
	for _, r := range v.String() {
		if r < '0' || r > '9' {
			return 0
		}
		value = value*10 + int(r-'0')
	}
	return value
}

func storeBest(key string, value int) {
	js.Global().Get("localStorage").Call("setItem", key, value)
}

func browserLanguage() string {
	if js.Global().Get("location").Get("search").String() == "?lang=ja" {
		return "ja"
	}
	return "en"
}
