package style

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	bold = lipgloss.NewStyle().Bold(true)

	None    = lipgloss.NewStyle().Render
	Bold    = bold.Render
	Red     = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render
	RedB    = bold.Copy().Foreground(lipgloss.Color("1")).Render
	Green   = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render
	GreenB  = bold.Copy().Foreground(lipgloss.Color("2")).Render
	Yellow  = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Render
	YellowB = bold.Copy().Foreground(lipgloss.Color("3")).Render
	Blue    = lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Render
	BlueB   = bold.Copy().Foreground(lipgloss.Color("4")).Render
	Cyan    = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render
	CyanB   = bold.Copy().Foreground(lipgloss.Color("6")).Render
	Grey    = bold.Copy().Foreground(lipgloss.Color("245")).Render
)
