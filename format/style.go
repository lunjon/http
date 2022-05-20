package format

import (
	"github.com/charmbracelet/lipgloss"
)

type Color uint8

const (
	ColorRed     Color = 31
	ColorGreen   Color = 32
	ColorYellow  Color = 33
	ColorBlue    Color = 34
	ColorMagenta Color = 35
	ColorCyan    Color = 36
)

type Styler struct {
	unit   lipgloss.Style
	whiteB lipgloss.Style
	red    lipgloss.Style
	redB   lipgloss.Style
	green  lipgloss.Style
	greenB lipgloss.Style
	blue   lipgloss.Style
	blueB  lipgloss.Style
	cyan   lipgloss.Style
	cyanB  lipgloss.Style
}

func NewStyler() *Styler {
	return &Styler{
		unit:   lipgloss.NewStyle(),
		whiteB: lipgloss.NewStyle().Bold(true),
		red:    lipgloss.NewStyle().Foreground(lipgloss.Color("31")),
		redB:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("31")),
		green:  lipgloss.NewStyle().Foreground(lipgloss.Color("32")),
		greenB: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("32")),
		blue:   lipgloss.NewStyle().Foreground(lipgloss.Color("34")),
		blueB:  lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("34")),
		cyan:   lipgloss.NewStyle().Foreground(lipgloss.Color("36")),
		cyanB:  lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("36")),
	}
}

func (s *Styler) Render() string {
	return ""
}

func (styler *Styler) WhiteB(s string) string {
	return styler.whiteB.Render(s)
}

func (styler *Styler) Red(s string) string {
	return styler.red.Render(s)
}

func (styler *Styler) RedB(s string) string {
	return styler.redB.Render(s)
}

func (styler *Styler) Green(s string) string {
	return styler.green.Render(s)
}

func (styler *Styler) GreenB(s string) string {
	return styler.greenB.Render(s)
}

func (styler *Styler) Blue(s string) string {
	return styler.blue.Render(s)
}

func (styler *Styler) BlueB(s string) string {
	return styler.blueB.Render(s)
}

func (styler *Styler) Cyan(s string) string {
	return styler.cyan.Render(s)
}

func (styler *Styler) CyanB(s string) string {
	return styler.cyanB.Render(s)
}
