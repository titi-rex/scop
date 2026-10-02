package parser

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"

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
	Name     string
	Vertices []float32
	Indices  []uint32
	ICount int32

	tmpIndices []int32
}

func NewObj(path string) (Object, error) {
	file, err := os.Open(path)
	if err != nil {
		return Object{}, fmt.Errorf("can't load object %s: %w", path, err)
	}
	defer func() { _ = file.Close() }()

	log.Println("Load obj from", path)

	obj := Object{}
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

	obj.convertIndices()
	obj.ICount = int32(len(obj.Indices))

	return obj, nil
}

func (o *Object) loadVertice(raw []string) error {
	for _, s := range raw {
		v, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return err
		}
		o.Vertices = append(o.Vertices, float32(v))
	}

	return nil
}

func (o *Object) loadFace(raw []string) error {
	var values []int32
	for _, s := range(raw) {
		v, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return err
		}
		if v == 0 {
			return errors.New("indice == 0")
		}
		values = append(values, int32(v))
	}

	var n int32
	for i := 1; i + 1 < len(raw); i += 1 {
		o.tmpIndices = append(o.tmpIndices, int32(values[0]))
		o.tmpIndices = append(o.tmpIndices, int32(values[i]))
		o.tmpIndices = append(o.tmpIndices, int32(values[i+1]))
		n += 3
	}

	return nil
}


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
