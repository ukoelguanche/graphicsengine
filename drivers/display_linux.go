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
	file       *os.File
	pixels     []byte
	buffer     []byte
	LineLength int // <--- Añade esto
	VW, VH     int // <--- Y esto
}

func InitDisplay(vw, vh int) *Display {

	exec.Command("stty", "-F", "/dev/tty", "-echo", "-icanon").Run()
	fmt.Print("\033[?25l")

	sw, sh = getDisplaySize()
	f, err := os.OpenFile("/dev/fb0", os.O_RDWR, 0)
	if err != nil {
		panic(err)
	}

	// OBTENER LINE LENGTH REAL
	// Esto evita que la imagen se vea "hacia un lado" o "torcida"
	var fixInfo struct {
		id                            [16]byte
		smem_start                    uintptr
		smem_len                      uint32
		type_                         uint32
		type_aux                      uint32
		visual                        uint32
		xpanstep, ypanstep, ywrapstep uint16
		line_length                   uint32 // Este es el que nos importa
	}
	// FBIOGET_FSCREENINFO = 0x4602
	syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), 0x4602, uintptr(unsafe.Pointer(&fixInfo)))

	lineLen := int(fixInfo.line_length)
	size := lineLen * sh // Tamaño real de la memoria de video

	data, _ := syscall.Mmap(int(f.Fd()), 0, size, syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED)
	backBuffer := make([]byte, size)

	fd := int(os.Stdin.Fd())

	oldState, err = term.MakeRaw(fd)
	if err != nil {
		panic(err)
	}

	return &Display{
		file:       f,
		pixels:     data,
		buffer:     backBuffer,
		LineLength: lineLen,
		VW:         vw,
		VH:         vh,
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
	// Usamos las constantes o variables VW y VH que tengas definidas
	if vx < 0 || vx >= int32(VW) || vy < 0 || vy >= int32(VH) {
		return
	}

	// Proyección dinámica: calculamos el área real que ocupa el píxel virtual
	// Esto reparte los píxeles sobrantes automáticamente
	xStart := int(float64(vx) * float64(sw) / float64(VW))
	xEnd := int(float64(vx+1) * float64(sw) / float64(VW))

	yStart := int(float64(vy) * float64(sh) / float64(VH))
	yEnd := int(float64(vy+1) * float64(sh) / float64(VH))

	r, g, b, a := c[0], c[1], c[2], c[3]

	// Dibujamos el bloque estirado
	for py := yStart; py < yEnd; py++ {
		for px := xStart; px < xEnd; px++ {
			// Importante: No olvides que si Alpine usa LineLength,
			// deberías usar d.LineLength en lugar de sw aquí.
			offset := (py*sw + px) * 4

			if offset+3 < len(d.buffer) {
				d.buffer[offset] = b
				d.buffer[offset+1] = g
				d.buffer[offset+2] = r
				d.buffer[offset+3] = a
			}
		}
	}
}

func (d *Display) Clear() {
	for i := range d.buffer {
		d.buffer[i] = 0
	}
}

func (d *Display) ClearFinal() {
	log.Printf("CLEARING display")
	for i := range d.buffer {
		d.buffer[i] = 0
	}
}

func (d *Display) Present() {
	copy(d.pixels, d.buffer)
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
