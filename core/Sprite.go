package core

type Sprite struct {
	Name         string
	BitmapSource string `json:"BitmapSource"`
	Bitmap       *Bitmap
	Frames       []Frame             `json:"Frames"`
	Sequences    map[string]Sequence `json:"Sequences"`
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

func (s *Sprite) GetBitmap() *Bitmap         { return s.Bitmap }
func (s *Sprite) GetFrame(index int32) Frame { return s.Frames[index] }

func (s *Sprite) CurrentSwapPaletteIndex() int {
	return int(float32(len(*s.PaletteSwap.TargetPalette)) * s.CurrentPalleteSwapPosition)
}
