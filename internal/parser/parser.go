package parser

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math"
	"os"

	"strconv"
	"strings"

	mgl "github.com/go-gl/mathgl/mgl32"
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
	Matrice      mgl.Mat4

	tmpIndices []int32
	max        []float64
	min        []float64
}

func NewObj(path string) (Object, error) {
	file, err := os.Open(path)
	if err != nil {
		return Object{}, fmt.Errorf("can't load object %s: %w", path, err)
	}
	defer func() { _ = file.Close() }()

	log.Println("Load obj from", path)

	obj := Object{
		Matrice: mgl.Ident4(),
		max:     []float64{-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64},
		min:     []float64{math.MaxFloat64, math.MaxFloat64, math.MaxFloat64},
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

	// fmt.Printf("vert bef: %v\n", obj.Vertices)
	obj.normalize()
	// fmt.Printf("vert aft: %v\n", obj.Vertices)
	obj.convertIndices()
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

	o.max[0] = max(x, o.max[0])
	o.min[0] = min(x, o.min[0])

	o.max[1] = max(y, o.max[1])
	o.min[1] = min(y, o.min[1])

	o.max[2] = max(z, o.max[2])
	o.min[2] = min(z, o.min[2])

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

	// triangle truc
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
	// log.Println("maxs:", o.max)
	// log.Println("mins:", o.min)

	// On centre et met le modele a l'echelle
	centerX := float32((o.min[0] + o.max[0]) / 2)
	centerY := float32((o.min[1] + o.max[1]) / 2)
	centerZ := float32((o.min[2] + o.max[2]) / 2)

	sizeX := o.max[0] - o.min[0]
	sizeY := o.max[1] - o.min[1]
	sizeZ := o.max[2] - o.min[2]

	maxExtent := max(sizeX, sizeY, sizeZ)

	// scaleFactor := (normalize && maxExtent > 0) ? (targetScale / maxExtent) : 1;
	log.Println("max extend", maxExtent)

	targetScale := float64(1)
	scaleFactor := float32(1)
	if maxExtent > 0 {
		scaleFactor = float32(targetScale / maxExtent)
	}

	vertices := make([]float32, len(o.Vertices))
	for i := 0; i < len(o.Vertices); i += 3 {
		vertices[i] = (o.Vertices[i] - centerX) * scaleFactor
		vertices[i+1] = (o.Vertices[i+1] - centerY) * scaleFactor
		vertices[i+2] = (o.Vertices[i+2] - centerZ) * scaleFactor
	}

	o.Vertices = vertices
}

// convert indices from 1..21, -1, -4,.. to >=0 indices
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
