// Copyright 2026 Ebi Showcase contributors. Licensed under Apache-2.0.
//go:build js && wasm

package raycasterui

import (
	"strconv"
	"syscall/js"
)

func browserLanguage() string {
	if js.Global().Get("location").Get("search").String() == "?lang=ja" {
		return "ja"
	}
	return "en"
}

func storedBest(key string) int {
	v := js.Global().Get("localStorage").Call("getItem", key)
	if v.IsNull() || v.IsUndefined() {
		return 0
	}
	n, _ := strconv.Atoi(v.String())
	return n
}

func storeBest(key string, value int) { js.Global().Get("localStorage").Call("setItem", key, value) }
