package core

type Color struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

func (c Color) Byte() []byte {
	return []byte{c.R, c.G, c.B, c.A}
}
