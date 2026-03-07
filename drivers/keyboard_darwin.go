package drivers

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Keyboard struct {
}

func InitKeyboard() *Keyboard {
	return &Keyboard{}
}

func (k *Keyboard) GetInput() KeyboardInput {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		if t, ok := event.(*sdl.KeyboardEvent); ok && t.Type == sdl.KEYDOWN {
			switch t.Keysym.Sym {
			case sdl.K_UP:
				return KBD_UP
			case sdl.K_SPACE:
				return KBD_SPACE
			case sdl.K_DOWN:
				return KBD_DOWN
			case sdl.K_LEFT:
				return KBD_LEFT
			case sdl.K_RIGHT:
				return KBD_RIGHT
			case sdl.K_ESCAPE:
				return KBD_ESCAPE
			case sdl.K_RETURN:
				return KBD_RETURN
			}
		}
		if _, ok := event.(*sdl.QuitEvent); ok {
			return KBD_ESCAPE
		}
	}
	return KBD_NONE
}
