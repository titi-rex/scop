package main

import (
	"fmt"
	"image"
	"image/draw"
	"os"
)

func NewImg(file string) (*image.RGBA, error) {

	imgFile, err := os.Open("./resources/wall.png")
	if err != nil {
		return nil, fmt.Errorf("texture %q not found on disk: %w", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	return rgba, nil
}
