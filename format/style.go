package format

import (
	"fmt"
	"strings"
)

func init() {
	WhiteB = newBuilder().bold(true).build()
	Red = newBuilder().fg(ColorRed).build()
	RedB = newBuilder().fg(ColorRed).bold(true).build()
	Green = newBuilder().fg(ColorGreen).build()
	GreenB = newBuilder().fg(ColorGreen).bold(true).build()
	Blue = newBuilder().fg(ColorBlue).build()
	BlueB = newBuilder().fg(ColorBlue).bold(true).build()
	Cyan = newBuilder().fg(ColorCyan).build()
	CyanB = newBuilder().fg(ColorCyan).bold(true).build()
}

type StyleFunc = func(string) string

type Color uint8

const (
	ColorRed     Color = 31
	ColorGreen   Color = 32
	ColorYellow  Color = 33
	ColorBlue    Color = 34
	ColorMagenta Color = 35
	ColorCyan    Color = 36
)

var (
	unitFunc = func(s string) string { return s }
	WhiteB   StyleFunc
	Red      StyleFunc
	RedB     StyleFunc
	Green    StyleFunc
	GreenB   StyleFunc
	Blue     StyleFunc
	BlueB    StyleFunc
	Cyan     StyleFunc
	CyanB    StyleFunc
)

func DisableColors() {
	WhiteB = unitFunc
	Red = unitFunc
	RedB = unitFunc
	Green = unitFunc
	GreenB = unitFunc
	Blue = unitFunc
	BlueB = unitFunc
}

type builder struct {
	fgColor Color
	isBold  bool
}

func newBuilder() *builder {
	return &builder{}
}

func (b *builder) fg(c Color) *builder {
	b.fgColor = c
	return b
}

func (b *builder) bold(v bool) *builder {
	b.isBold = v
	return b
}

func (b *builder) build() StyleFunc {
	codes := []string{}
	if b.fgColor > 0 {
		codes = append(codes, fmt.Sprint(b.fgColor))
	}

	if b.isBold {
		codes = append(codes, "1")
	}

	format := "\x1b[" + strings.Join(codes, ";") + "m%s\x1b[0m"
	return func(s string) string {
		return fmt.Sprintf(format, s)
	}
}
