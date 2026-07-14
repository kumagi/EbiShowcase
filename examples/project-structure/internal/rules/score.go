// Copyright 2026 Ebi Showcase contributors. Licensed under Apache-2.0.
package rules

// Collect is pure game logic, so it can be tested without a window.
func Collect(score, gems int) (int, int) { return score + 10, gems + 1 }
