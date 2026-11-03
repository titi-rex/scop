package texture

import (
	"fmt"
)

const (
	mixSpeed = 0.05
	minRatio = 0.0
	maxRatio = 1.0
)

type mixState int

const (
	Increasing mixState = iota
	Up
	Decreasing
	Down
)

type MixRatio struct {
	ratio float32
	state mixState
}

func NewMixRatio(initialRatio float32) MixRatio {
	return MixRatio{
		ratio: initialRatio,
		state: Down,
	}
}

func (m *MixRatio) SetState(s mixState) {
	m.state = s
}

func (m *MixRatio) Toggle() {
	switch m.state {
	case Down, Decreasing:
		m.state = Increasing
	case Up, Increasing:
		m.state = Decreasing
	}
}

func (m *MixRatio) Update() {
	switch m.state {
	case Down, Up:
		return
	case Decreasing:
		m.ratio -= mixSpeed
		if m.ratio <= minRatio {
			m.ratio = minRatio
			m.state = Down
		}
	case Increasing:
		m.ratio += mixSpeed
		if m.ratio >= maxRatio {
			m.ratio = maxRatio
			m.state = Up
		}
	}
}

func (m *MixRatio) Ratio() float32 {
	return m.ratio
}

func (m *MixRatio) String() string {
	return fmt.Sprintf("ratio: %f, %d", m.ratio, m.state)
}
