package ogm

import (
	"math"
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

func Sincos(a float32) (float32, float32) {
	sin64, cos64 := math.Sincos(float64(a))
	return float32(sin64), float32(cos64)
}

func Cos(a float32) float32 {
	return float32(math.Cos(float64(a)))
}

func Sin(a float32) float32 {
	return float32(math.Sin(float64(a)))
}

func Tan(a float32) float32 {
	return float32(math.Tan(float64(a)))
}

func Sqrt(f float32) float32 {
	return float32(math.Sqrt(float64(f)))
}

func Abs(f float32) float32 {
	return float32(math.Abs(float64(f)))
}

func Clamp0(value float32) float32 {
	if value <= 0 {
		return 0
	}
	return value
}
