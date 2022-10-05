package style

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	bold    = lipgloss.NewStyle().Bold(true)
	None    = lipgloss.NewStyle()
	Bold    = bold
	Red     = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	RedB    = bold.Copy().Foreground(lipgloss.Color("1"))
	Green   = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	GreenB  = bold.Copy().Foreground(lipgloss.Color("2"))
	Yellow  = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	YellowB = bold.Copy().Foreground(lipgloss.Color("3"))
	Blue    = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	BlueB   = bold.Copy().Foreground(lipgloss.Color("4"))
	Cyan    = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	CyanB   = bold.Copy().Foreground(lipgloss.Color("6"))
	Grey    = bold.Copy().Foreground(lipgloss.Color("245"))
)
