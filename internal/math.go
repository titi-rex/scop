package main

import "math"

const (
	 rader = math.Pi / 180
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
