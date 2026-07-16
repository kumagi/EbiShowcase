// Command gen-favicon creates the shared site icon from Ebi Tenjiroh's
// original character artwork.
//
// SPDX-License-Identifier: Apache-2.0
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"

	xdraw "golang.org/x/image/draw"
)

const iconSize = 512

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	sourcePath := filepath.Join(root, "assets", "characters", "ebi-boy.png")
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		panic(err)
	}
	source, err := png.Decode(sourceFile)
	_ = sourceFile.Close()
	if err != nil {
		panic(err)
	}

	icon := image.NewNRGBA(image.Rect(0, 0, iconSize, iconSize))
	center := iconSize / 2
	radius := 248
	background := color.NRGBA{R: 14, G: 27, B: 62, A: 255}
	for y := 0; y < iconSize; y++ {
		for x := 0; x < iconSize; x++ {
			dx, dy := x-center, y-center
			if dx*dx+dy*dy <= radius*radius {
				icon.SetNRGBA(x, y, background)
			}
		}
	}

	// The source is 390×904. Its head occupies the top 390-pixel square.
	// Keeping this crop here makes the favicon reproducible from the licensed
	// character source rather than maintaining a second hand-edited portrait.
	head := source.(interface {
		SubImage(image.Rectangle) image.Image
	}).SubImage(image.Rect(0, 0, 390, 390))
	xdraw.CatmullRom.Scale(icon, image.Rect(12, 12, 500, 500), head, head.Bounds(), xdraw.Over, nil)

	outputPath := filepath.Join(root, "web", "assets", "favicon.png")
	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	if err := png.Encode(outputFile, icon); err != nil {
		_ = outputFile.Close()
		panic(err)
	}
	if err := outputFile.Close(); err != nil {
		panic(err)
	}
	fmt.Printf("Generated %s from %s\n", outputPath, sourcePath)
}
