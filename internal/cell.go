package main

import (
	"math/rand"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Cell struct {
	drawable uint32

	x, y             int
	alive, aliveNext bool
}

type World [][]*Cell

const (
	rows      = 15
	columns   = 15
	threshold = 0.25
)

var (
	gCount = 0
	square = []float32{
		-0.5, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,

		-0.5, 0.5, 0,
		0.5, 0.5, 0,
		0.5, -0.5, 0,
	}
)

func makeWorld() World {
	// r := rand.New(rand.NewSource(12))
	cells := make([][]*Cell, rows, columns)
	for x := range rows {
		for y := range columns {
			c := newCell(x, y)

			c.alive = rand.Float64() < threshold
			c.aliveNext = c.alive

			cells[x] = append(cells[x], c)
		}
	}
	return cells
}

func newCell(x, y int) *Cell {
	points := make([]float32, len(square))
	copy(points, square)

	for i := range len(points) {
		var pos, size float32
		switch i % 3 {
		case 0:
			size = 1.0 / float32(columns)
			pos = float32(x) * size
		case 1:
			size = 1.0 / float32(rows)
			pos = float32(y) * size
		default:
			continue
		}

		if points[i] < 0 {
			points[i] = (pos * 2) - 1
		} else {
			points[i] = ((pos + size) * 2) - 1
		}
	}
	return &Cell{
		drawable: makeVao(points),
		x:        x,
		y:        y,
	}
}

func (c *Cell) draw() {
	if !c.alive {
		return
	}

	gl.BindVertexArray(c.drawable)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square)/3))
}

func (c *Cell) checkState(w World) {
	c.alive = c.aliveNext
	c.aliveNext = c.alive

	liveCount := c.liveNeighbors(w)
	if c.alive {
		// 1. Any live cell with fewer than two live neighbours dies, as if caused by underpopulation.
		if liveCount < 2 {
			c.aliveNext = false
		}

		// 2. Any live cell with two or three live neighbours lives on to the next generation.
		if liveCount == 2 || liveCount == 3 {
			c.aliveNext = true
		}

		// 3. Any live cell with more than three live neighbours dies, as if by overpopulation.
		if liveCount > 3 {
			c.aliveNext = false
		}
	} else {
		// 4. Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.
		if liveCount == 3 {
			c.aliveNext = true
		}
	}

	if c.alive {
		gCount++
	}
}

func (c *Cell) liveNeighbors(w World) int {
	var count int
	add := func(x, y int) {
		// If we're at an edge, check the other side of the board.
		if x == len(w) {
			x = 0
		} else if x == -1 {
			x = len(w) - 1
		}
		if y == len(w[x]) {
			y = 0
		} else if y == -1 {
			y = len(w[x]) - 1
		}

		if w[x][y].alive {
			count++
		}
	}
	add(c.x-1, c.y)
	add(c.x+1, c.y)
	add(c.x, c.y+1)
	add(c.x, c.y-1)
	add(c.x-1, c.y+1)
	add(c.x+1, c.y+1)
	add(c.x-1, c.y-1)
	add(c.x+1, c.y-1)
	return count
}
