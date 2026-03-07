package core

type Sprite struct {
	Name         string
	BitmapSource string `json:"BitmapSource"`
	Bitmap       *Bitmap
	Frames       []Rect           `json:"Frames"`
	Sequences    map[string][]int `json:"Sequences"`
	Characters   map[string]int   `json:"Characters"`
	PaletteSwap  PaletteSwap      `json:"PaletteSwap"`

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
	return s.Bitmap
}

func (s *Sprite) CurrentSwapPaletteIndex() int {
	return int(float32(len(*s.PaletteSwap.TargetPalette)) * s.CurrentPalleteSwapPosition)
}

func (s *Sprite) ProcessColor(color []byte) []byte {
	if s.PaletteSwap.TargetPalette == nil {
		return color
	}

	index := s.CurrentSwapPaletteIndex()
	color = s.PaletteSwap.SourcePalette.ReplacePalette(color, s.PaletteSwap.TargetPalette, index)

	return color
}
