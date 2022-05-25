package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func main() {
	nbrRows := 25
	nbrCols := 10

	rows := [][]string{}

	for row := 1; row <= nbrRows; row++ {
		columns := []string{}
		for col := 0; col < nbrCols; col++ {
			color := fmt.Sprint(row + nbrRows*col)
			style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
			columns = append(columns, style.Render(color))
		}
		rows = append(rows, columns)
	}

	for _, row := range rows {
		s := strings.Join(row, "\t")
		fmt.Println(s)
	}
}
