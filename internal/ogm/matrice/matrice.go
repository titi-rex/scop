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

// Det return the matrice determinant
func (m Matrice) Det() float32 {
	return m[0]*m[5]*m[10]*m[15] - m[0]*m[5]*m[11]*m[14] - m[0]*m[6]*m[9]*m[15] +
		m[0]*m[6]*m[11]*m[13] + m[0]*m[7]*m[9]*m[14] - m[0]*m[7]*m[10]*m[13] -
		m[1]*m[4]*m[10]*m[15] + m[1]*m[4]*m[11]*m[14] + m[1]*m[6]*m[8]*m[15] -
		m[1]*m[6]*m[11]*m[12] - m[1]*m[7]*m[8]*m[14] + m[1]*m[7]*m[10]*m[12] +
		m[2]*m[4]*m[9]*m[15] - m[2]*m[4]*m[11]*m[13] - m[2]*m[5]*m[8]*m[15] +
		m[2]*m[5]*m[11]*m[12] + m[2]*m[7]*m[8]*m[13] - m[2]*m[7]*m[9]*m[12] -
		m[3]*m[4]*m[9]*m[14] + m[3]*m[4]*m[10]*m[13] + m[3]*m[5]*m[8]*m[14] -
		m[3]*m[5]*m[10]*m[12] - m[3]*m[6]*m[8]*m[13] + m[3]*m[6]*m[9]*m[12]
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
