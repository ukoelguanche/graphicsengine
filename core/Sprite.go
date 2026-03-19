package core

type Sprite struct {
	Name           string
	BitmapSource   string `json:"BitmapSource"`
	Bitmap         *Bitmap
	PaletteBitmaps []*Bitmap
	Frames         []Frame             `json:"Frames"`
	Sequences      map[string]Sequence `json:"Sequences"`
	// ToDO: Remove this property as it is intended to be used for texts, not for sprites
	Characters  map[string]int `json:"Characters"`
	PaletteSwap PaletteSwap    `json:"PaletteSwap"`

	RelativePaletteSwapSpeed   float32 `json:"RelativePaletteSwapSpeed"`
	CurrentPalleteSwapOffset   float32
	CurrentPalleteSwapPosition float32
}

type Sprites struct {
	BitmapSources map[string]string   `json:"BitmapSources"`
	Sprites       map[string]*Sprite  `json:"sprites"`
	Palettes      map[string]*Palette `json:"Palettes"`
}

func (s *Sprite) GetBitmap() *Bitmap {
	if len(s.PaletteBitmaps) == 0 {
		return s.Bitmap
	}

	index := s.CurrentSwapPaletteIndex()
	if index < 0 || index >= len(s.PaletteBitmaps) {
		return s.Bitmap
	}

	return s.PaletteBitmaps[index]
}
func (s *Sprite) GetFrame(index int32) Frame { return s.Frames[index] }

func (s *Sprite) CurrentSwapPaletteIndex() int {
	if s.PaletteSwap.TargetPalette == nil || len(*s.PaletteSwap.TargetPalette) == 0 {
		return 0
	}

	paletteLen := len(*s.PaletteSwap.TargetPalette)
	index := int(float32(paletteLen) * s.CurrentPalleteSwapPosition)
	if index >= paletteLen {
		return paletteLen - 1
	}

	return index
}
