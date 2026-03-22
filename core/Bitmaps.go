package core

type Bitmap struct {
	Name          string
	W             int32
	H             int32
	Pixels        []byte
	IndexedPixels []byte
	Palette       *Palette
}

type Bitmaps map[string]*Bitmap
