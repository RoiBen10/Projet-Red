package main

import (
	"fmt"
	"strings"
)

// nameSelectScreen --> saisie du nom via une grille de lettres façon Undertale : flèches pour se
// déplacer, Entrée pour choisir une lettre, "Rub" efface la dernière, "End" valide (nom non vide).
func nameSelectScreen() string {
	grid := [][]string{
		{"A", "B", "C", "D", "E", "F"},
		{"G", "H", "I", "J", "K", "L"},
		{"M", "N", "O", "P", "Q", "R"},
		{"S", "T", "U", "V", "W", "X"},
		{"Y", "Z", "Rub", "End"},
	}
	const maxLen = 12
	const innerWidth = 30
	const cellWidth = 5

	row, col := 0, 0
	var name []rune

	fmt.Print("\033[H\033[2J")

	draw := func() {
		cols, rows := terminalSize()

		var out []string
		out = append(out, "┌"+strings.Repeat("─", innerWidth+2)+"┐")
		out = append(out, "│ "+centerText("Quel est ton nom, Voyageur ?", innerWidth)+" │")
		out = append(out, "└"+strings.Repeat("─", innerWidth+2)+"┘")
		out = append(out, "")

		out = append(out, "┌"+strings.Repeat("─", innerWidth+2)+"┐")
		out = append(out, "│ \033[33m"+centerText(string(name), innerWidth)+"\033[37m │")
		out = append(out, "└"+strings.Repeat("─", innerWidth+2)+"┘")
		out = append(out, "")

		out = append(out, "┌"+strings.Repeat("─", innerWidth+2)+"┐")
		for r, gridRow := range grid {
			var line strings.Builder
			line.WriteString("│ ")
			visibleLen := 0
			for c, cell := range gridRow {
				cellStr := padRightText(cell, cellWidth)
				visibleLen += len(cellStr)
				if r == row && c == col {
					line.WriteString("\033[33m" + cellStr + "\033[37m")
				} else {
					line.WriteString(cellStr)
				}
			}
			pad := innerWidth - visibleLen
			if pad < 0 {
				pad = 0
			}
			line.WriteString(strings.Repeat(" ", pad))
			line.WriteString(" │")
			out = append(out, line.String())
		}
		out = append(out, "└"+strings.Repeat("─", innerWidth+2)+"┘")

		boxOuterWidth := innerWidth + 4
		leftPad := (cols - boxOuterWidth) / 2
		if leftPad < 0 {
			leftPad = 0
		}
		topPad := (rows - len(out)) / 2
		if topPad < 0 {
			topPad = 0
		}
		left := strings.Repeat(" ", leftPad)

		var b strings.Builder
		b.WriteString("\033[H\033[37m")
		for i := 0; i < topPad; i++ {
			b.WriteString("\r\n")
		}
		for _, l := range out {
			b.WriteString(left + l + "\r\n")
		}
		b.WriteString("\033[0m")
		fmt.Print(b.String())
	}
	draw()

	for {
		key, ok := <-arrowChan
		if !ok {
			return "Voyageur"
		}
		switch key {
		case "up":
			if row > 0 {
				row--
				if col >= len(grid[row]) {
					col = len(grid[row]) - 1
				}
			}
		case "down":
			if row < len(grid)-1 {
				row++
				if col >= len(grid[row]) {
					col = len(grid[row]) - 1
				}
			}
		case "left":
			if col > 0 {
				col--
			}
		case "right":
			if col < len(grid[row])-1 {
				col++
			}
		case "enter":
			switch grid[row][col] {
			case "Rub":
				if len(name) > 0 {
					name = name[:len(name)-1]
				}
			case "End":
				if len(name) > 0 {
					return string(name)
				}
			default:
				if len(name) < maxLen {
					name = append(name, []rune(grid[row][col])[0])
				}
			}
		case "quit":
			if len(name) > 0 {
				return string(name)
			}
			return "Voyageur"
		}
		draw()
	}
}
