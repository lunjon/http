package style

import (
	"fmt"
	"strings"
)

type StyleFunc = func(string) string

type Color uint8

const (
	Red     Color = 31
	Green   Color = 32
	Yellow  Color = 33
	Blue    Color = 34
	Magenta Color = 35
	Cyan    Color = 36
)

type Styler struct {
	format string
}

func (styler *Styler) Style(s string) string {
	return fmt.Sprintf(styler.format, s)
}

type Builder struct {
	fg   Color
	bold bool
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Fg(c Color) *Builder {
	b.fg = c
	return b
}

func (b *Builder) Bold(v bool) *Builder {
	b.bold = v
	return b
}

func (b *Builder) Build() StyleFunc {
	codes := []string{}
	if b.fg > 0 {
		codes = append(codes, fmt.Sprint(b.fg))
	}

	if b.bold {
		codes = append(codes, "1")
	}

	format := "\x1b[" + strings.Join(codes, ";") + "m%s\x1b[0m"
	return func(s string) string {
		return fmt.Sprintf(format, s)
	}
}
