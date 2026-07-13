//go:build !js || !wasm

package main

// Native builds use a no-op store so pure logic can compile and be tested
// without a browser. Persistence remains a browser-only adapter.
func storageGet(string) (string, bool) { return "", false }
func storageSet(string, string)        {}
func storageRemove(string)             {}
