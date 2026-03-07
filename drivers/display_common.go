package drivers

import (
	"apodeiktikos.com/fbtest/core"
)

const (
	VW, VH = 320, 200
)

var GlobalDisplay *Display

func (d *Display) DrawSpriteRect(sprite *core.Sprite, rect core.Rect, position core.Point) {
	bitmap := sprite.GetBitmap()
	for sy := 0; sy < int(rect.Size.H); sy++ {
		for sx := 0; sx < int(rect.Size.W); sx++ {
			// Calculate original position inside bitmap
			origX := rect.Point.X + float64(sx)
			origY := rect.Point.Y + float64(sy)

			// Avoid drawing outside bounds
			if origX < 0 || origX >= float64(bitmap.W) || origY < 0 || origY >= float64(bitmap.H) {
				continue
			}

			srcOff := int((origY*float64(bitmap.W) + origX) * 4)
			color := bitmap.Pixels[srcOff : srcOff+4]

			// Skip transparencies
			if color[3] < 128 {
				continue
			}

			finalColor := sprite.ProcessColor(color)

			d.DrawPixel(int32(position.X)+int32(sx), int32(position.Y)+int32(sy), finalColor)
		}
	}

	if sprite.CurrentPalleteSwapPosition >= 1 {
		sprite.CurrentPalleteSwapPosition = 0
	}
	sprite.CurrentPalleteSwapPosition += sprite.CurrentPalleteSwapOffset
}
