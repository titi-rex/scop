package matrice

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		m1       *Matrice
		m2       *Matrice
		expected *Matrice
	}{
		"1. add identity": {
			m1: &Matrice{
				1, 1, 1, 0,
				0, 0, 0, 1,
				0, 0, 0, 0,
				1, 0, 0, 1,
			},
			m2: NewIdent(),
			expected: &Matrice{
				2, 1, 1, 0,
				0, 1, 0, 1,
				0, 0, 1, 0,
				1, 0, 0, 2,
			},
		},
		"2. add Zero": {
			m1: NewIdent(),
			m2: NewZero(),
			expected:NewIdent(),
		},
		"2. add negative": {
			m1: NewIdent(),
			m2: &Matrice{
				-1, 1, -1, 0,
				0, -1, 0, 1,
				0, 0, 1, 0,
				-1, -3, 0, 2,
			},
			expected: &Matrice{
				0, 1, -1, 0,
				0, 0, 0, 1,
				0, 0, 2, 0,
				-1, -3, 0, 3,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// Given - two Matrice

			// Whan - Add operator
			r := Add(test.m1, test.m2)

			// Then - Matrice addition !
			assert.True(t, reflect.DeepEqual(r, test.expected))
		})
	}
}

func TestSelfAdd(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		m1       *Matrice
		m2       *Matrice
		m3       *Matrice
		expected *Matrice
	}{
		"1. add identity": {
			m1: &Matrice{
				1, 1, 1, 0,
				0, 0, 0, 1,
				0, 0, 0, 0,
				1, 0, 0, 1,
			},
			m2: NewIdent(),
			m3: NewIdent(),
			expected: &Matrice{
				3, 1, 1, 0,
				0, 2, 0, 1,
				0, 0, 2, 0,
				1, 0, 0, 3,
			},
		},
		"2. add Zero": {
			m1: NewIdent(),
			m2: NewZero(),
			m3: NewZero(),
			expected:NewIdent(),
		},
		"2. add negative": {
			m1: NewIdent(),
			m2: &Matrice{
				-1, 1, -1, 0,
				0, -1, 0, 1,
				0, 0, 1, 0,
				-1, -3, 0, 2,
			},
			m3: NewIdent(),
			expected: &Matrice{
				1, 1, -1, 0,
				0, 1, 0, 1,
				0, 0, 3, 0,
				-1, -3, 0, 4,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// Given - two Matrice

			// Whan - Add operator
			test.m1.Add(test.m2).Add(test.m3)

			assert.True(t, reflect.DeepEqual(test.m1, test.expected))
		})
	}
}
