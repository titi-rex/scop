package main

import (
	"image"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func NewTexture(rgba *image.RGBA) uint32 {
	var id uint32
	gl.Enable(gl.TEXTURE_2D)
	gl.GenTextures(1, &id)
	gl.BindTexture(gl.TEXTURE_2D, id)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	gl.GenerateMipmap(id)
	return id
}
