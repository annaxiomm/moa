package screen

import (
	"github.com/annaxiomm/moa/cpu"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	Width       = 128
	Height      = 128
	Scale       = 4
	Framebuffer = uint16(0x6600)
)

type Screen struct {
	Cpu    *cpu.Cpu
	Image  *ebiten.Image
	Pixels []byte
}

func New(c *cpu.Cpu) *Screen {
	return &Screen{
		Cpu:    c,
		Image:  ebiten.NewImage(Width, Height),
		Pixels: make([]byte, Width*Height*4),
	}
}

func (s *Screen) Update() error {
	// Execute up to 1000 instructions per frame.
	for i := 0; i < 1000 && !s.Cpu.Halted; i++ {
		if err := s.Cpu.Step(); err != nil {
			return err
		}
	}

	return nil
}

func (s *Screen) Draw(dst *ebiten.Image) {
	// Read the framebuffer and convert each pixel to RGBA.
	for i := range Width * Height {
		value := s.Cpu.Memory.Read(Framebuffer + uint16(i))
		offset := i * 4

		s.Pixels[offset] = value   // Red
		s.Pixels[offset+1] = value // Green
		s.Pixels[offset+2] = value // Blue
		s.Pixels[offset+3] = 255   // Alpha
	}

	s.Image.WritePixels(s.Pixels)
	dst.DrawImage(s.Image, nil)
}

func (s *Screen) Layout(outsideWidth, outsideHeight int) (int, int) {
	return Width, Height
}

func (s *Screen) Run() error {
	ebiten.SetWindowSize(Width*Scale, Height*Scale)
	ebiten.SetWindowTitle("Moa v0.1")

	return ebiten.RunGame(s)
}
