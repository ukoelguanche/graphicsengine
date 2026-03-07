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
	window   *sdl.Window
	renderer *sdl.Renderer
	tex      *sdl.Texture
	pixels   []byte
}

func InitDisplay(vw, vh int) *Display {
	//sw, sh := getDisplaySize()
	sw, sh := 640, 480
	log.Printf("Detected resolution: %dx%d", sw, sh)

	sdl.Init(sdl.INIT_EVERYTHING)
	w, _ := sdl.CreateWindow("framebuffer", 100, 100, int32(sw), int32(sh), sdl.WINDOW_SHOWN)
	r, _ := sdl.CreateRenderer(w, -1, sdl.RENDERER_ACCELERATED)
	t, _ := r.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(vw), int32(vh))
	return &Display{w, r, t, make([]byte, vw*vh*4)}
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
	d.tex.Update(nil, unsafe.Pointer(&d.pixels[0]), VW*4)
	d.renderer.Copy(d.tex, nil, nil)
	d.renderer.Present()
}

func (d *Display) Close() {
	d.Clear()
	d.window.Destroy()
	sdl.Quit()
}
