package core

type Palette []Color

func (palette Palette) GradientIndex(color []byte) int {
	for i, c := range palette {
		if color[0] == c.R && color[1] == c.G && color[2] == c.B {
			return i
		}
	}

	return -1
}

func (p Palette) ReplacePalette(color []byte, targetPalette *Palette, animationIndex int) []byte {
	gradientIndex := p.GradientIndex(color)

	if gradientIndex >= 0 {
		return (*targetPalette)[(gradientIndex+animationIndex)%len(*targetPalette)].Byte()
	}
	
	return color
}
