//go:build !js || !wasm

package reversiui

func storedBest(string) int   { return 0 }
func storeBest(string, int)   {}
func browserLanguage() string { return "en" }
