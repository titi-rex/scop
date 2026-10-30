package matrice

import (
	"bytes"
	"fmt"
	"text/tabwriter"
)

type Matrice [4 * 4]float32

// Data access the internal storage
func (m *Matrice) Data() *float32 {
	return &m[0]
}

// Add perform matrice addition
func (m *Matrice) Add(o *Matrice) *Matrice {
	*m = *Add(m, o)

	return m
}

// Mul perform matrice substraction
func (m *Matrice) Mul(o *Matrice) *Matrice {
	*m = *Mul(m, o)

	return m
}

// Scale perform matrice substraction
func (m *Matrice) Scale(s float32) *Matrice {
	*m = *Scale(m, s)

	return m
}

type v4 [4]float32

// row return a vector representing the corresponding row (staring from 0)
func (m Matrice) row(id int) v4 {
	return v4{m[id+0], m[id+4], m[id+8], m[id+12]}
}

// Pretty print a matrice
func (m Matrice) String() string {
	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 4, 4, 1, ' ', tabwriter.AlignRight)
	for i := range 4 {
		for _, col := range m.row(i) {
			_, _ = fmt.Fprintf(w, "%f\t", col)
		}

		_, _ = fmt.Fprintln(w, "")
	}
	_ = w.Flush()
	return buf.String()
}
