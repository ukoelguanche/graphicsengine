package loaders

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/ukoelguanche/graphicsengine/core"
)

func TestLoadIndexedSprites(t *testing.T) {
	tempDir := t.TempDir()

	sprites := core.Sprites{
		BitmapSources: map[string]string{
			"sheet": "./assets/sheet.png",
		},
		Sprites: map[string]*core.Sprite{
			"Hero": {
				BitmapSource: "sheet",
				Frames: []core.Frame{
					{Point: core.Point{X: 0, Y: 0}, Size: core.Size{W: 2, H: 1}},
				},
				Sequences: map[string]core.Sequence{"idle": {0}},
			},
		},
		Palettes: map[string]*core.Palette{},
	}

	definitionPath := filepath.Join(tempDir, "Sprites.json")
	writeJSON(t, definitionPath, sprites)

	palettePath := filepath.Join(tempDir, "palette.pal")
	if err := os.WriteFile(palettePath, []byte{
		0, 0, 0,
		10, 20, 30,
		40, 50, 60,
	}, 0o644); err != nil {
		t.Fatalf("write palette: %v", err)
	}

	indexedPath := filepath.Join(tempDir, "sheet.idx")
	if err := os.WriteFile(indexedPath, []byte{1, 0, 2, 1}, 0o644); err != nil {
		t.Fatalf("write idx: %v", err)
	}

	metadata := indexedAssetsMetadata{
		SpritesFile: "./Sprites.json",
		PalettePath: "./palette.pal",
		PaletteType: "rgb24",
		ColorCount:  3,
		Bitmaps: []indexedBitmapMetadata{
			{
				Name:        "sheet",
				SourcePath:  "./assets/sheet.png",
				Width:       2,
				Height:      2,
				IndexedPath: "./sheet.idx",
			},
		},
	}

	metadataPath := filepath.Join(tempDir, "metadata.json")
	writeJSON(t, metadataPath, metadata)

	var loaded core.Sprites
	LoadIndexedSprites(definitionPath, metadataPath, &loaded)

	sprite := loaded.Sprites["Hero"]
	if sprite == nil {
		t.Fatalf("expected Hero sprite to be loaded")
	}
	if sprite.Bitmap == nil {
		t.Fatalf("expected Hero bitmap to be loaded")
	}
	if sprite.Bitmap.Palette == nil {
		t.Fatalf("expected Hero palette to be loaded")
	}

	wantIndexed := []byte{1, 0, 2, 1}
	if string(sprite.Bitmap.IndexedPixels) != string(wantIndexed) {
		t.Fatalf("unexpected indexed pixels: got %v want %v", sprite.Bitmap.IndexedPixels, wantIndexed)
	}

	if len(*sprite.Bitmap.Palette) != 3 {
		t.Fatalf("unexpected palette size: got %d want 3", len(*sprite.Bitmap.Palette))
	}

	if (*sprite.Bitmap.Palette)[0].A != 0 {
		t.Fatalf("expected palette index 0 to be transparent")
	}

	color := (*sprite.Bitmap.Palette)[2]
	if color.R != 40 || color.G != 50 || color.B != 60 || color.A != 255 {
		t.Fatalf("unexpected palette color at index 2: %+v", color)
	}
}

func writeJSON(t *testing.T, path string, value any) {
	t.Helper()

	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal json: %v", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write json %s: %v", path, err)
	}
}
