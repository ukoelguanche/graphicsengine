package core

type Bitmap struct {
	Name   string
	W      int32
	H      int32
	Pixels []byte
}

type Bitmaps map[string]*Bitmap
