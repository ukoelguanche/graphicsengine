package loaders

import (
	"log"

	"apodeiktikos.com/fbtest/core"
)

func LoadSprites(definitionPath string, sprites *core.Sprites) {
	LoadJson(definitionPath, sprites)

	// Load all bitmaps
	bitmaps := make(core.Bitmaps)
	for name, path := range sprites.BitmapSources {
		log.Printf("Loading bitmap source: %s %s", name, path)
		bitmaps[name] = LoadBitmap(path)
		bitmaps[name].Name = name
	}

	// Assign bitmap pointers to sprites
	for name, sprite := range sprites.Sprites {
		sprite.Name = name
		sprite.Bitmap = bitmaps[sprite.BitmapSource]
		log.Printf("Bitmap [%s] assigned to sprite: [%s]", sprite.BitmapSource, sprite.Name)
	}

	// Assign palete pointers to sprites
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
		}

	}
}
