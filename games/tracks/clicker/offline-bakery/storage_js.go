//go:build js && wasm

package main

import "syscall/js"

func storageGet(key string) (string, bool) {
	value := js.Global().Get("localStorage").Call("getItem", key)
	if value.Type() != js.TypeString {
		return "", false
	}
	return value.String(), true
}

func storageSet(key, value string) {
	js.Global().Get("localStorage").Call("setItem", key, value)
}

func storageRemove(key string) {
	js.Global().Get("localStorage").Call("removeItem", key)
}
