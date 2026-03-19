package drivers

import (
	"github.com/ukoelguanche/graphicsengine/core"
)

type PixelTransformer interface {
	Transform(pixels []byte)
}

type ManagedPixelTransformer interface {
	PixelTransformer
	IsFinished() bool
	Complete()
}

type SpriteColorProcessor interface {
	ProcessColor(color []byte) []byte
}

type SpriteScopedColorProcessor interface {
	SpriteColorProcessor
	AppliesTo(sprite *core.Sprite) bool
}

const (
	VW, VH = 320, 200
)

var GlobalDisplay *Display
var SpriteColorProcessors []SpriteColorProcessor = make([]SpriteColorProcessor, 0)

func (d *Display) DrawSpriteRect(sprite *core.Sprite, rect core.Frame, position core.Point) {
	bitmap := sprite.GetBitmap()
	drawX := int(position.X)
	drawY := int(position.Y)
	srcX0 := int(rect.Point.X)
	srcY0 := int(rect.Point.Y)
	width := int(rect.Size.W)
	height := int(rect.Size.H)

	if width <= 0 || height <= 0 {
		return
	}

	if drawX >= VW || drawY >= VH || drawX+width <= 0 || drawY+height <= 0 {
		return
	}

	if srcX0 < 0 {
		drawX -= srcX0
		width += srcX0
		srcX0 = 0
	}

	if srcY0 < 0 {
		drawY -= srcY0
		height += srcY0
		srcY0 = 0
	}

	if srcX0+width > int(bitmap.W) {
		width = int(bitmap.W) - srcX0
	}

	if srcY0+height > int(bitmap.H) {
		height = int(bitmap.H) - srcY0
	}

	if drawX < 0 {
		srcX0 -= drawX
		width += drawX
		drawX = 0
	}

	if drawY < 0 {
		srcY0 -= drawY
		height += drawY
		drawY = 0
	}

	if drawX+width > VW {
		width = VW - drawX
	}

	if drawY+height > VH {
		height = VH - drawY
	}

	if width <= 0 || height <= 0 {
		return
	}

	processors := applicableProcessors(sprite)

	for sy := 0; sy < height; sy++ {
		srcRowOffset := ((srcY0 + sy) * int(bitmap.W) * 4) + srcX0*4
		dstY := drawY + sy

		for sx := 0; sx < width; sx++ {
			srcOff := srcRowOffset + sx*4
			color := bitmap.Pixels[srcOff : srcOff+4]

			if color[3] < 128 {
				continue
			}

			if len(processors) > 0 {
				for _, processor := range processors {
					color = processor.ProcessColor(color)
				}
			}

			d.DrawPixel(int32(drawX+sx), int32(dstY), color)
		}
	}

	if sprite.CurrentPalleteSwapPosition >= 1 {
		sprite.CurrentPalleteSwapPosition = 0
	}
	sprite.CurrentPalleteSwapPosition += sprite.CurrentPalleteSwapOffset
}

func (d *Display) AddTransformer(t PixelTransformer) {
	d.transformers = append(d.transformers, t)
}

func (d *Display) RemoveTransformer(t PixelTransformer) {
	for i, transformer := range d.transformers {
		if transformer != t {
			continue
		}

		copy(d.transformers[i:], d.transformers[i+1:])
		d.transformers[len(d.transformers)-1] = nil
		d.transformers = d.transformers[:len(d.transformers)-1]
		return
	}
}

func applicableProcessors(sprite *core.Sprite) []SpriteColorProcessor {
	if len(SpriteColorProcessors) == 0 {
		return nil
	}

	processors := make([]SpriteColorProcessor, 0, len(SpriteColorProcessors))
	for _, processor := range SpriteColorProcessors {
		scoped, ok := processor.(SpriteScopedColorProcessor)
		if ok && !scoped.AppliesTo(sprite) {
			continue
		}
		processors = append(processors, processor)
	}

	return processors
}
