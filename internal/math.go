package main

import (
	"math"

	mgl "github.com/go-gl/mathgl/mgl32"
)

const (
	rader  = math.Pi / 180
	degrer = 180 / math.Pi
)

type float interface {
	~float32 | ~float64
}

func Rad[T float](d T) T {
	return d * rader
}

func Deg[T float](r T) T {
	return r * degrer
}

func RotateY(d float32) mgl.Mat4 {
	r := Rad(float64(d))
	cos := float32(math.Cos(r))
	sin := float32(math.Sin(r))

	return mgl.Mat4{
		cos, 0, -sin, 0,
		0, 1, 0, 0,
		sin, 0, cos, 0,
		0, 0, 0, 1,
	}
}
