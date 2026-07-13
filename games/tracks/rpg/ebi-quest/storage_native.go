//go:build !js || !wasm

package main

func storageGet(string) (string, bool) { return "", false }
func storageSet(string, string)        {}
func storageRemove(string)             {}
