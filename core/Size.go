package core

type Size struct {
	W, H float64
}

func (s Size) SetW(w float64) Size      { return Size{W: w, H: s.H} }
func (s Size) SetH(h float64) Size      { return Size{W: s.W, H: h} }
func (s Size) Mult(factor float64) Size { return Size{W: s.W * factor, H: s.H * factor} }
