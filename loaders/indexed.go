package loaders

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/ukoelguanche/graphicsengine/core"
)

type indexedBitmapMetadata struct {
	Name        string `json:"name"`
	SourcePath  string `json:"sourcePath"`
	Width       int32  `json:"width"`
	Height      int32  `json:"height"`
	IndexedPath string `json:"indexedPath"`
}

type indexedAssetsMetadata struct {
	SpritesFile string                  `json:"spritesFile"`
	PalettePath string                  `json:"palettePath"`
	PaletteType string                  `json:"paletteType"`
	ColorCount  int                     `json:"colorCount"`
	Bitmaps     []indexedBitmapMetadata `json:"bitmaps"`
}

func LoadIndexedSprites(definitionPath string, metadataPath string, sprites *core.Sprites) {
	LoadJson(definitionPath, sprites)

	metadata := loadIndexedMetadata(metadataPath)
	palette := loadRGB24Palette(resolveRelativePath(metadataPath, metadata.PalettePath), metadata.ColorCount)
	bitmaps := loadIndexedBitmaps(metadataPath, metadata, palette)

	assignSpritesAndPalettes(sprites, bitmaps)
}

func loadIndexedMetadata(metadataPath string) indexedAssetsMetadata {
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		log.Fatalf("Error opening indexed metadata %s\n%v", metadataPath, err)
	}

	var metadata indexedAssetsMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		log.Fatalf("Error loading indexed metadata: %s\n%v", metadataPath, err)
	}

	if metadata.PaletteType != "rgb24" {
		log.Fatalf("Unsupported indexed palette type %q in %s", metadata.PaletteType, metadataPath)
	}

	return metadata
}

func loadRGB24Palette(path string, colorCount int) *core.Palette {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error opening palette file %s\n%v", path, err)
	}

	expectedBytes := colorCount * 3
	if len(data) != expectedBytes {
		log.Fatalf("Invalid palette file %s: expected %d bytes, got %d", path, expectedBytes, len(data))
	}

	palette := make(core.Palette, colorCount)
	for i := 0; i < colorCount; i++ {
		offset := i * 3
		alpha := uint8(255)
		if i == 0 {
			alpha = 0
		}
		palette[i] = core.Color{
			R: data[offset],
			G: data[offset+1],
			B: data[offset+2],
			A: alpha,
		}
	}

	return &palette
}

func loadIndexedBitmaps(metadataPath string, metadata indexedAssetsMetadata, palette *core.Palette) core.Bitmaps {
	bitmaps := make(core.Bitmaps, len(metadata.Bitmaps))

	for _, bitmapMetadata := range metadata.Bitmaps {
		indexedPath := resolveRelativePath(metadataPath, bitmapMetadata.IndexedPath)
		log.Printf("Loading indexed bitmap source: %s %s", bitmapMetadata.Name, indexedPath)

		bitmaps[bitmapMetadata.Name] = loadIndexedBitmap(indexedPath, bitmapMetadata, palette)
		bitmaps[bitmapMetadata.Name].Name = bitmapMetadata.Name
	}

	return bitmaps
}

func loadIndexedBitmap(path string, metadata indexedBitmapMetadata, palette *core.Palette) *core.Bitmap {
	indexedPixels, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error opening indexed bitmap %s\n%v", path, err)
	}

	expectedPixelCount := int(metadata.Width * metadata.Height)
	if len(indexedPixels) != expectedPixelCount {
		log.Fatalf("Invalid indexed bitmap %s: expected %d bytes, got %d", path, expectedPixelCount, len(indexedPixels))
	}

	for i, paletteIndex := range indexedPixels {
		if int(paletteIndex) >= len(*palette) {
			log.Fatalf("Invalid palette index %d in %s at pixel %d", paletteIndex, path, i)
		}
	}

	return &core.Bitmap{
		W:             metadata.Width,
		H:             metadata.Height,
		IndexedPixels: indexedPixels,
		Palette:       palette,
	}
}

func resolveRelativePath(baseFile string, targetPath string) string {
	if filepath.IsAbs(targetPath) {
		return targetPath
	}

	baseDir := filepath.Dir(baseFile)
	resolved := filepath.Join(baseDir, targetPath)
	cleanResolved := filepath.Clean(resolved)

	if _, err := os.Stat(cleanResolved); err == nil {
		return cleanResolved
	}

	workspaceResolved := filepath.Clean(targetPath)
	if _, err := os.Stat(workspaceResolved); err == nil {
		return workspaceResolved
	}

	log.Fatalf("Unable to resolve path %q relative to %q", targetPath, baseFile)
	return ""
}
