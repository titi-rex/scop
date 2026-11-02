package vector

func (v *Vector) Add(o Vector) *Vector {
	*v = Add(*v, o)

	return v
}

func (v *Vector) Sub(o Vector) *Vector {
	*v = Sub(*v, o)

	return v
}

func (v *Vector) Mul(s float32) *Vector {
	*v = Mul(*v, s)

	return v
}

func (v *Vector) Cross(o Vector) *Vector {
	*v = Cross(*v, o)

	return v
}

func (v *Vector) Dot(o Vector) float32 {
	return Dot(*v, o)
}

func (v *Vector) Normalize() *Vector {
	*v = Normalize(*v)

	return v
}

func (v *Vector) Abs() *Vector {
	*v = Abs(*v)

	return v
}

func (v *Vector) Clamp0() *Vector {
	*v = Clamp0(*v)

	return v
}


func (v *Vector) DivVec(o Vector) *Vector {
	*v = DivVec(*v, o)

	return v
}
