package parser

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"scop/internal/ogm/matrice"
	"scop/internal/ogm/vector"

	"strconv"
	"strings"
)

type token string

const (
	wordSep token = " "
	comment token = "#"
	objName token = "o"
	vertice token = "v"
	face    token = "f"
)

type Object struct {
	Name         string
	Vertices     []float32
	Indices      []uint32
	IndicesCount int32
	Matrice      matrice.Matrice

	tmpIndices []int32
	max        vector.Vector
	min        vector.Vector
}

func NewObj(path string) (Object, error) {
	file, err := os.Open(path)
	if err != nil {
		return Object{}, fmt.Errorf("can't load object %s: %w", path, err)
	}
	defer func() { _ = file.Close() }()

	log.Println("Load obj from", path)

	obj := Object{
		Matrice: *matrice.NewIdent(),
		max:     vector.Vector{-math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32},
		min:     vector.Vector{math.MaxFloat32, math.MaxFloat32, math.MaxFloat32},
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		words := strings.Split(scanner.Text(), " ")
		switch token(words[0]) {
		case comment:
			continue
		case objName:
			obj.Name = words[1]
		case vertice:
			if err := obj.loadVertice(words[1:]); err != nil {
				return Object{}, err
			}
		case face:
			if err := obj.loadFace(words[1:]); err != nil {
				return Object{}, err
			}
		default:
			log.Println("unknow line id", words[0])
		}
	}

	err = scanner.Err()
	if err != nil {
		return Object{}, err
	}

	obj.convertFacesToTriangle()
	obj.normalize()
	// fmt.Printf("\n\naft: %+v\n", obj.Vertices)
	// fmt.Printf("\n\nvef %+v\n", obj.Indices)
	obj.IndicesCount = int32(len(obj.Indices))

	return obj, nil
}

func (o *Object) loadVertice(raw []string) error {

	x, errX := strconv.ParseFloat(raw[0], 32)
	y, errY := strconv.ParseFloat(raw[1], 32)
	z, errZ := strconv.ParseFloat(raw[2], 32)
	err := errors.Join(errX, errY, errZ)
	if err != nil {
		return err
	}

	if len(raw) > 3 {
		log.Println("ignored garbage:", raw[3:])
	}

	o.max[0] = max(float32(x), o.max[0])
	o.min[0] = min(float32(x), o.min[0])

	o.max[1] = max(float32(y), o.max[1])
	o.min[1] = min(float32(y), o.min[1])

	o.max[2] = max(float32(z), o.max[2])
	o.min[2] = min(float32(z), o.min[2])

	o.Vertices = append(o.Vertices, float32(x))
	o.Vertices = append(o.Vertices, float32(y))
	o.Vertices = append(o.Vertices, float32(z))

	return nil
}

func (o *Object) loadFace(raw []string) error {
	var values []int32
	for _, s := range raw {
		v, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return err
		}
		if v == 0 {
			return errors.New("indice == 0")
		}
		values = append(values, int32(v))
	}

	// triangle fan
	var n int32
	for i := 1; i+1 < len(raw); i += 1 {
		o.tmpIndices = append(o.tmpIndices, int32(values[0]))
		o.tmpIndices = append(o.tmpIndices, int32(values[i]))
		o.tmpIndices = append(o.tmpIndices, int32(values[i+1]))
		n += 3
	}

	return nil
}

func (o *Object) normalize() {
	const targetScale = float32(1)

	center := vector.Add(o.min, o.max)
	center.Mul(0.5)

	size := vector.Sub(o.max, o.min)

	maxExtent := max(size.X(), size.Y(), size.Z())

	scaleFactor := float32(1)
	if maxExtent > 0 {
		scaleFactor = float32(targetScale / maxExtent)
	}

	// vertices := make([]float32, len(o.Indices))
	var vertices []float32
	var indices []uint32
	currentIdx := uint32(0)
	for i := 0; i < len(o.Indices); i += 3 {
		idx0 := o.Indices[i]
		idx1 := o.Indices[i+1]
		idx2 := o.Indices[i+2]

		v0 := vector.LoadFromSlice(o.Vertices, idx0*3)
		v1 := vector.LoadFromSlice(o.Vertices, idx1*3)
		v2 := vector.LoadFromSlice(o.Vertices, idx2*3)

		a := vector.Sub(v1, v0)
		b := vector.Sub(v2, v0)
		n := vector.Cross(a, b)
		n.Normalize().Abs()

		v0.Sub(center).Mul(scaleFactor)
		v1.Sub(center).Mul(scaleFactor)
		v2.Sub(center).Mul(scaleFactor)

		u, v := ProjTexture(v0, o.min, size, n)
		vertices = append(vertices, []float32{v0.X(), v0.Y(), v0.Z(), u, v}...)
		indices = append(indices, currentIdx)
		currentIdx++

		vertices = append(vertices, []float32{v1.X(), v1.Y(), v1.Z(), u, v}...)
		indices = append(indices, currentIdx)
		currentIdx++

		vertices = append(vertices, []float32{v2.X(), v2.Y(), v2.Z(), u, v}...)
		indices = append(indices, currentIdx)
		currentIdx++
	}

	o.Vertices = vertices
	o.Indices = indices
}

func ProjTexture(vertice, min, size, normal vector.Vector) (float32, float32) {
	uv := vector.Sub(vertice, min)
	uv.DivVec(size).Clamp0()

	// choose proj with the greater normal axis
	if normal.Z() >= normal.X() && normal.Z() >= normal.Y() { // plan XY
		return uv.X(), uv.Y()
	} else if normal.X() >= normal.Y() && normal.X() >= normal.Z() { // plan YZ
		return uv.Z(), uv.Y()
	} else { // plan XZ
		return uv.X(), uv.Z()
	}
}

// convert indices from 1..21, -1, -4,.. to >=0 indices
func (o *Object) convertFacesToTriangle() {
	nv := int32(len(o.Vertices))

	for _, v := range o.tmpIndices {
		if v > 0 {
			v -= 1
		} else if v < 0 {
			v = nv + v
		}
		o.Indices = append(o.Indices, uint32(v))
	}
}
