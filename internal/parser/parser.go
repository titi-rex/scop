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
	"slices"

	"strconv"
	"strings"
)

type token string

const (
	wordSep   token = " "
	comment   token = "#"
	objName   token = "o"
	vertice   token = "v"
	texture   token = "vt"
	face      token = "f"
	matLib    token = "mtllib"
	matUse    token = "usemtl"
	smoothing token = "s"
	group     token = "g"
	emptyLine token = ""
)

var (
	ErrUnknowToken     = errors.New("unknow ligne prefix")
	ErrMissingVertices = errors.New("missing vertices")
	ErrMissingIndices  = errors.New("missing indices")
)

// Object is the struct holding the vertices data
type Object struct {
	Name     string
	Vertices []float32
	Indices  []uint32
	Matrice  matrice.Matrice

	tmpIndices []int32
	max        vector.Vector
	min        vector.Vector
}

// NewObj load an object from an .obj file
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
		case smoothing, matLib, matUse, emptyLine, group, texture:
			continue
		default:
			return Object{}, fmt.Errorf("%s: %w", words[0], ErrUnknowToken)
		}
	}

	if err = scanner.Err(); err != nil {
		return Object{}, fmt.Errorf("scanner: %w", err)
	}

	obj.convertIndices()

	if obj.CountIndices() == 0 {
		return Object{}, fmt.Errorf("len check: %w", ErrMissingIndices)
	}

	if obj.CountVertices() == 0 {
		return Object{}, fmt.Errorf("len check: %w", ErrMissingVertices)
	}

	if slices.Max(obj.Indices) >= uint32(obj.CountVertices()) {
		return Object{}, ErrMissingVertices
	}

	if err = obj.normalize(); err != nil {
		return Object{}, fmt.Errorf("normalize: %w", err)
	}

	return obj, nil
}

// loadVertice parse a vertice line, it expect 3 points.
//
// example:
//
//	v 1.21513 2.213553 0.5654654
func (o *Object) loadVertice(raw []string) error {
	x, errX := strconv.ParseFloat(raw[0], 32)
	y, errY := strconv.ParseFloat(raw[1], 32)
	z, errZ := strconv.ParseFloat(raw[2], 32)
	if err := errors.Join(errX, errY, errZ); err != nil {
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

// loadFace parse a face line with an arbitrary number of vertices referred
// and make a triangle fan with. [https://en.wikipedia.org/wiki/Triangle_fan]
//
// example:
//
//	f 1 5 3 14
func (o *Object) loadFace(line []string) error {
	var values []int32
	for _, s := range line {
		raw := strings.Split(s, "/")

		value, err := strconv.ParseInt(raw[0], 10, 32)
		if err != nil {
			return err

		}
		if value == 0 {
			return errors.New("indice == 0")
		}

		values = append(values, int32(value))
	}

	// triangle fan
	var n int32
	for i := 1; i+1 < len(line); i += 1 {
		o.tmpIndices = append(o.tmpIndices, int32(values[0]))
		o.tmpIndices = append(o.tmpIndices, int32(values[i]))
		o.tmpIndices = append(o.tmpIndices, int32(values[i+1]))
		n += 3
	}

	return nil
}

// convertIndices convert indices from 1..21, -1, -4,.. to >=0 indices.
func (o *Object) convertIndices() {
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

// normalize center and scale the object.
// It also apply a planar projection to compute texture coordinates.
func (o *Object) normalize() error {
	const targetScale = float32(2)

	center := vector.Add(o.min, o.max)
	center.Mul(0.5)

	size := vector.Sub(o.max, o.min)
	if size.Norm() == 0 {
		return errors.New("all points at same position")
	}

	greaterSide := max(size.X(), size.Y(), size.Z())
	scaleFactor := float32(targetScale / greaterSide)

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

		for _, vec := range []vector.Vector{v0, v1, v2} {
			vec.Sub(center)

			// texture planar projection
			textVec := vec.Clone()
			textVec.Mul(1 / scaleFactor)
			u, v := ProjTexture(textVec, o.min, size, n)

			vec.Mul(scaleFactor)
			vertices = append(vertices, []float32{vec.X(), vec.Y(), vec.Z(), u, v}...)
			indices = append(indices, currentIdx)
			currentIdx++
		}
	}

	o.Vertices = vertices
	o.Indices = indices
	return nil
}

// ProjTexture project texture onto the nearer plan.
func ProjTexture(vertice, min, size, normal vector.Vector) (float32, float32) {
	duv := vector.Sub(vertice, min)
	uv := vector.Zip(duv, size, func(e1, e2 float32) float32 {
		if e2 < 0 {
			return 0
		}
		return e1 / e2
	})

	// choose proj with the greater normal axis
	if normal.Z() >= normal.X() && normal.Z() >= normal.Y() { // plan XY
		return uv.X(), uv.Y()
	} else if normal.X() >= normal.Y() && normal.X() >= normal.Z() { // plan YZ
		return uv.Z(), uv.Y()
	} else { // plan XZ
		return uv.X(), uv.Z()
	}
}

func (o *Object) CountVertices() int32 {
	return int32(len(o.Vertices) / 3)
}

func (o *Object) CountIndices() int32 {
	return int32(len(o.Indices))
}
