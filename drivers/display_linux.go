package drivers

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/term"
)

type inputEvent struct {
	Type  uint16
	Code  uint16
	Value int32
}

const (
	KeyEnter = 28
	KeyUp    = 103
	KeyDown  = 108
)

var sw, sh int
var oldState *term.State

type Display struct {
	file         *os.File
	pixels       []byte
	buffer       []byte
	baseBuffer   []byte
	scaledRow    []byte
	LineLength   int
	VW, VH       int
	xStarts      []int
	xEnds        []int
	yStarts      []int
	yEnds        []int
	transformers []PixelTransformer
}

func InitDisplay(title string, vw, vh int) *Display {

	exec.Command("stty", "-F", "/dev/tty", "-echo", "-icanon").Run()
	fmt.Print("\033[?25l")

	sw, sh = getDisplaySize()
	f, err := os.OpenFile("/dev/fb0", os.O_RDWR, 0)
	if err != nil {
		panic(err)
	}

	// Get actual line length to avoid tilted or off-centered images
	var fixInfo struct {
		id                            [16]byte
		smem_start                    uintptr
		smem_len                      uint32
		type_                         uint32
		type_aux                      uint32
		visual                        uint32
		xpanstep, ypanstep, ywrapstep uint16
		line_length                   uint32 // the important one
	}
	// FBIOGET_FSCREENINFO = 0x4602
	syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), 0x4602, uintptr(unsafe.Pointer(&fixInfo)))

	lineLen := int(fixInfo.line_length)
	size := lineLen * sh // Real video memory size

	data, _ := syscall.Mmap(int(f.Fd()), 0, size, syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED)

	fd := int(os.Stdin.Fd())

	oldState, err = term.MakeRaw(fd)
	if err != nil {
		panic(err)
	}

	return &Display{
		file:         f,
		pixels:       data,
		buffer:       make([]byte, vw*vh*4),
		scaledRow:    make([]byte, sw*4),
		LineLength:   lineLen,
		VW:           vw,
		VH:           vh,
		xStarts:      buildScaleStarts(vw, sw),
		xEnds:        buildScaleEnds(vw, sw),
		yStarts:      buildScaleStarts(vh, sh),
		yEnds:        buildScaleEnds(vh, sh),
		transformers: make([]PixelTransformer, 0),
	}
}
func getDisplaySize() (int, int) {
	vsBytes, err := os.ReadFile("/sys/class/graphics/fb0/virtual_size")
	if err != nil {
		panic(err)
	}

	parts := strings.Split(strings.TrimSpace(string(vsBytes)), ",")
	realWidth, _ := strconv.Atoi(parts[0])
	realHeight, _ := strconv.Atoi(parts[1])

	return realWidth, realHeight
}

func (d *Display) DrawPixel(vx, vy int32, c []byte) {
	if vx < 0 || vx >= int32(VW) || vy < 0 || vy >= int32(VH) {
		return
	}

	offset := (int(vy)*d.VW + int(vx)) * 4
	copy(d.buffer[offset:offset+4], c)
}

func (d *Display) Clear() {
	for i := range d.buffer {
		d.buffer[i] = 0
	}
}

func (d *Display) ClearFinal() {
	log.Printf("CLEARING display")
	for i := range d.pixels {
		d.pixels[i] = 0
	}
}

func (d *Display) Present() {
	for _, transformer := range d.transformers {
		transformer.Transform(d.buffer)
	}

	d.scaleVirtualBuffer()
}

func (d *Display) Close() {
	log.Printf("CLOSING display")
	d.Clear()
	for i := range d.pixels {
		d.pixels[i] = 0
	}

	syscall.Munmap(d.pixels)

	if d.file != nil {
		d.file.Close()
	}

	if oldState != nil {
		term.Restore(int(os.Stdin.Fd()), oldState)
	}

	syscall.Munmap(d.pixels)
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

func (d *Display) scaleVirtualBuffer() {
	for vy := 0; vy < d.VH; vy++ {
		yStart := d.yStarts[vy]
		yEnd := d.yEnds[vy]

		d.buildScaledRow(vy)

		for py := yStart; py < yEnd; py++ {
			rowOffset := py * d.LineLength
			copy(d.pixels[rowOffset:rowOffset+len(d.scaledRow)], d.scaledRow)
		}
	}
}

func (d *Display) buildScaledRow(vy int) {
	srcRowOffset := vy * d.VW * 4

	for vx := 0; vx < d.VW; vx++ {
		xStart := d.xStarts[vx]
		xEnd := d.xEnds[vx]

		srcOffset := srcRowOffset + vx*4
		r := d.buffer[srcOffset]
		g := d.buffer[srcOffset+1]
		b := d.buffer[srcOffset+2]
		a := d.buffer[srcOffset+3]

		for px := xStart; px < xEnd; px++ {
			dstOffset := px * 4
			d.scaledRow[dstOffset] = b
			d.scaledRow[dstOffset+1] = g
			d.scaledRow[dstOffset+2] = r
			d.scaledRow[dstOffset+3] = a
		}
	}
}

func buildScaleStarts(virtualSize, realSize int) []int {
	starts := make([]int, virtualSize)
	for i := 0; i < virtualSize; i++ {
		starts[i] = i * realSize / virtualSize
	}
	return starts
}

func buildScaleEnds(virtualSize, realSize int) []int {
	ends := make([]int, virtualSize)
	for i := 0; i < virtualSize; i++ {
		ends[i] = (i + 1) * realSize / virtualSize
	}
	return ends
}
