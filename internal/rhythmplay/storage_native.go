// Copyright 2026 Ebi Showcase contributors. Licensed under Apache-2.0.
//go:build !js || !wasm

package rhythmplay

func browserLanguage() string { return "en" }
func storedInt(string) int    { return 0 }
func storeInt(string, int)    {}
