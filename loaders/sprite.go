package loaders

import (
	"log"

	"github.com/ukoelguanche/graphicsengine/core"
)

func LoadSprites(definitionPath string, sprites *core.Sprites) {
	LoadJson(definitionPath, sprites)

	bitmaps := loadBitmapSources(sprites.BitmapSources)
	assignSpritesAndPalettes(sprites, bitmaps)
}

func loadBitmapSources(bitmapSources map[string]string) core.Bitmaps {
	bitmaps := make(core.Bitmaps)
	for name, path := range bitmapSources {
		log.Printf("Loading bitmap source: %s %s", name, path)
		bitmaps[name] = LoadBitmap(path)
		bitmaps[name].Name = name
	}
	return bitmaps
}

func assignSpritesAndPalettes(sprites *core.Sprites, bitmaps core.Bitmaps) {
	for name, sprite := range sprites.Sprites {
		sprite.Name = name
		sprite.Bitmap = bitmaps[sprite.BitmapSource]
		log.Printf("Bitmap [%s] assigned to sprite: [%s]", sprite.BitmapSource, sprite.Name)
	}

	for _, sprite := range sprites.Sprites {
		if sprite.PaletteSwap.SourcePaletteName != "" {
			sprite.PaletteSwap.SourcePalette = sprites.Palettes[sprite.PaletteSwap.SourcePaletteName]
			log.Printf("Source palette [%s] assigned to sprite: [%s]", sprite.PaletteSwap.SourcePaletteName, sprite.Name)
		}
		if sprite.PaletteSwap.TargetPaletteName != "" {
			sprite.PaletteSwap.TargetPalette = sprites.Palettes[sprite.PaletteSwap.TargetPaletteName]
			log.Printf("Target palette [%s] assigned to sprite: [%s]", sprite.PaletteSwap.TargetPaletteName, sprite.Name)

			sprite.CurrentPalleteSwapOffset = 1 / float32(len(*sprite.PaletteSwap.TargetPalette)) * sprite.RelativePaletteSwapSpeed
			sprite.CurrentPalleteSwapPosition = 0.0
			sprite.PaletteBitmaps = buildPaletteBitmaps(sprite)
		}
	}
}

func buildPaletteBitmaps(sprite *core.Sprite) []*core.Bitmap {
	if sprite.Bitmap == nil || sprite.PaletteSwap.SourcePalette == nil || sprite.PaletteSwap.TargetPalette == nil {
		return nil
	}

	targetPalette := *sprite.PaletteSwap.TargetPalette
	if len(targetPalette) == 0 {
		return nil
	}

	if len(sprite.Bitmap.IndexedPixels) > 0 && sprite.Bitmap.Palette != nil {
		bitmaps := make([]*core.Bitmap, len(targetPalette))
		for animationIndex := range targetPalette {
			shiftedPalette := buildShiftedPalette(*sprite.Bitmap.Palette, sprite.PaletteSwap.SourcePalette, sprite.PaletteSwap.TargetPalette, animationIndex)
			bitmaps[animationIndex] = &core.Bitmap{
				Name:          sprite.Bitmap.Name,
				W:             sprite.Bitmap.W,
				H:             sprite.Bitmap.H,
				IndexedPixels: sprite.Bitmap.IndexedPixels,
				Palette:       &shiftedPalette,
			}
		}
		return bitmaps
	}

	bitmaps := make([]*core.Bitmap, len(targetPalette))
	for animationIndex := range targetPalette {
		pixels := make([]byte, len(sprite.Bitmap.Pixels))
		copy(pixels, sprite.Bitmap.Pixels)

		for i := 0; i < len(pixels); i += 4 {
			if pixels[i+3] < 128 {
				continue
			}

			gradientIndex := sprite.PaletteSwap.SourcePalette.GradientIndex(pixels[i : i+4])
			if gradientIndex < 0 {
				continue
			}

			targetColor := targetPalette[(gradientIndex+animationIndex)%len(targetPalette)]
			pixels[i] = targetColor.R
			pixels[i+1] = targetColor.G
			pixels[i+2] = targetColor.B
			pixels[i+3] = targetColor.A
		}

		bitmaps[animationIndex] = &core.Bitmap{
			Name:   sprite.Bitmap.Name,
			W:      sprite.Bitmap.W,
			H:      sprite.Bitmap.H,
			Pixels: pixels,
		}
	}

	return bitmaps
}

func buildShiftedPalette(base core.Palette, sourcePalette *core.Palette, targetPalette *core.Palette, animationIndex int) core.Palette {
	shifted := make(core.Palette, len(base))
	copy(shifted, base)

	for i, color := range shifted {
		gradientIndex := sourcePalette.GradientIndex([]byte{color.R, color.G, color.B, color.A})
		if gradientIndex < 0 {
			continue
		}

		targetColor := (*targetPalette)[(gradientIndex+animationIndex)%len(*targetPalette)]
		shifted[i] = targetColor
	}

	return shifted
}
