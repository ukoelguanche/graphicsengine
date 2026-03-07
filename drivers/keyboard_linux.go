package drivers

import (
	"encoding/binary"
	"log"
	"os"
	"strings"
	"syscall"
)

type Keyboard struct {
	keyboardFile *os.File
}

func InitKeyboard() *Keyboard {
	kbdPath := findKeyboardDevice()
	keyboardFile, err := os.OpenFile(kbdPath, os.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		log.Fatalf("Failed to open keyboard file: %s", err)
	}
	log.Printf("Keyboard file is: %s", kbdPath)

	return &Keyboard{keyboardFile: keyboardFile}
}

func (k *Keyboard) GetInput() KeyboardInput {
	if k.keyboardFile == nil {
		return KBD_NONE
	}

	buffer := make([]byte, 256)
	n, err := syscall.Read(int(k.keyboardFile.Fd()), buffer)

	if err != nil || n < 24 {
		return KBD_NONE
	}

	for i := 0; i+24 <= n; i += 24 {
		chunk := buffer[i : i+24]

		typ := binary.LittleEndian.Uint16(chunk[16:18])
		code := binary.LittleEndian.Uint16(chunk[18:20])
		val := binary.LittleEndian.Uint32(chunk[20:24])

		if typ == 1 {
			if val == 1 || val == 2 { // Pulsado o mantenido
				switch code {
				case 1:
					return KBD_ESCAPE
				case 28:
					return KBD_RETURN
				case 75:
					return KBD_SPACE
				case 103:
					return KBD_UP
				case 108:
					return KBD_DOWN
				case 105:
					return KBD_LEFT
				case 106:
					return KBD_RIGHT
				default:
					return KBD_NONE
				}
			} else if val == 0 {
				return KBD_NONE
			}
		}
	}

	return KBD_NONE
}

func findKeyboardDevice() string {
	data, err := os.ReadFile("/proc/bus/input/devices")
	if err != nil {
		log.Printf("error reading /proc/bus/input/devices: %v. Keyboard not found", err)
		return ""
	}

	sections := strings.Split(string(data), "\n\n")
	var keyboardSection string

	for _, section := range sections {
		if isKeyboardSection(section) {
			keyboardSection = section
			break
		}
	}

	lines := strings.Split(keyboardSection, "\n")
	for _, line := range lines {
		if !strings.Contains(line, "Handlers=") {
			continue
		}
		parts := strings.Fields(line)
		for _, event := range parts {
			if strings.HasPrefix(event, "event") {
				return "/dev/input/" + event
			}
		}
	}

	return ""
}

func isKeyboardSection(section string) bool {
	if !strings.Contains(section, "H: Handlers=sysrq kbd event") {
		return false
	}
	if !(strings.Contains(section, "B: EV=120013") || strings.Contains(section, "B: EV=120011")) {
		return false
	}
	if !strings.Contains(section, "P: Phys=usb-") {
		return false
	}

	return true
}
