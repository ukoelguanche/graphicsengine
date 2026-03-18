package drivers

/*
#cgo LDFLAGS: -framework ApplicationServices
#include <ApplicationServices/ApplicationServices.h>
*/
import "C"
import (
	"log"

	"github.com/veandco/go-sdl2/sdl"

	"unsafe"
)

type Display struct {
	window       *sdl.Window
	renderer     *sdl.Renderer
	tex          *sdl.Texture
	pixels       []byte
	buffer       []byte
	baseBuffer   []byte
	transformers []PixelTransformer
}

func InitDisplay(title string, vw, vh int) *Display {
	sw, sh := 960, 600
	log.Printf("Detected resolution: %dx%d", sw, sh)

	sdl.Init(sdl.INIT_EVERYTHING)
	w, _ := sdl.CreateWindow(title, 100, 100, int32(sw), int32(sh), sdl.WINDOW_SHOWN)
	r, _ := sdl.CreateRenderer(w, -1, sdl.RENDERER_ACCELERATED)
	t, _ := r.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(vw), int32(vh))

	return &Display{
		window:       w,
		renderer:     r,
		tex:          t,
		pixels:       make([]byte, vw*vh*4),
		buffer:       make([]byte, vw*vh*4),
		transformers: make([]PixelTransformer, 0),
	}
}

func getDisplaySize() (int, int) {
	mainDisplay := C.CGMainDisplayID()
	width := C.CGDisplayPixelsWide(mainDisplay)
	height := C.CGDisplayPixelsHigh(mainDisplay)

	return int(width), int(height)
}

func (d *Display) DrawPixel(x, y int32, c []byte) {
	if x < 0 || x >= VW || y < 0 || y >= VH {
		return
	}

	offset := (y*VW + x) * 4
	copy(d.pixels[offset:offset+4], c)
}

func (d *Display) Clear() {
	for i := range d.pixels {
		d.pixels[i] = 0
	}
}

func (d *Display) Present() {
	copy(d.pixels, d.buffer)

	for _, transformer := range d.transformers {
		transformer.Transform(d.pixels)
	}

	// 3. Renderizamos (usando d.pixels ya transformados)
	d.tex.Update(nil, unsafe.Pointer(&d.pixels[0]), VW*4)
	d.renderer.Copy(d.tex, nil, nil)
	d.renderer.Present()
}

func (d *Display) Close() {
	d.Clear()
	d.window.Destroy()
	sdl.Quit()
}

func (d *Display) SaveBaseBuffer() {
	d.baseBuffer = make([]byte, len(d.buffer))
	copy(d.baseBuffer, d.buffer)
}

func (d *Display) RestoreBaseBuffer() bool {
	if len(d.baseBuffer) == 0 {
		return false
	}
	copy(d.buffer, d.baseBuffer)
	return true
}
