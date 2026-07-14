// Copyright 2026 Ebi Showcase contributors. Licensed under Apache-2.0.
//go:build !js || !wasm

package raycasterui

func browserLanguage() string { return "en" }
func storedBest(string) int   { return 0 }
func storeBest(string, int)   {}
