package main

import (
	"fmt"
	"strings"
)

// centerText / padRightText --> alignement dans une largeur donnée (en runes).
func centerText(s string, w int) string {
	pad := w - len([]rune(s))
	if pad < 0 {
		pad = 0
	}
	left := pad / 2
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", pad-left)
}

func padRightText(s string, w int) string {
	pad := w - len([]rune(s))
	if pad < 0 {
		pad = 0
	}
	return s + strings.Repeat(" ", pad)
}

// Écran de choix de classe : classSelectScreen --> boîte simple avec bordure et curseur en forme
// de cœur qui se déplace, fond noir, façon menu de sélection Undertale (pas d'animation de fond,
// contrairement à l'écran titre).
func classSelectScreen() string {
	const title = "Choisis ta classe"
	const boxWidth = 22 // largeur du contenu, hors bordures
	options := []string{"Humain", "Elfe", "Nain"}
	selected := 0

	fmt.Print("\033[H\033[2J")

	draw := func() {
		cols, rows := terminalSize()
		totalLines := 5 + len(options) // bordure haut + titre + blanc + options + blanc + bordure bas

		boxOuterWidth := boxWidth + 4
		leftPad := (cols - boxOuterWidth) / 2
		if leftPad < 0 {
			leftPad = 0
		}
		topPad := (rows - totalLines) / 2
		if topPad < 0 {
			topPad = 0
		}
		left := strings.Repeat(" ", leftPad)

		var b strings.Builder
		b.WriteString("\033[H\033[37m")
		for i := 0; i < topPad; i++ {
			b.WriteString("\r\n")
		}

		b.WriteString(left + "┌" + strings.Repeat("─", boxWidth+2) + "┐\r\n")
		b.WriteString(left + "│ " + centerText(title, boxWidth) + " │\r\n")
		b.WriteString(left + "│ " + strings.Repeat(" ", boxWidth) + " │\r\n")
		for i, opt := range options {
			cursor := "  "
			if i == selected {
				cursor = "\033[33m🔔\033[37m"
			}
			b.WriteString(left + "│ " + cursor + padRightText(opt, boxWidth-2) + " │\r\n")
		}
		b.WriteString(left + "│ " + strings.Repeat(" ", boxWidth) + " │\r\n")
		b.WriteString(left + "└" + strings.Repeat("─", boxWidth+2) + "┘\r\n")
		b.WriteString("\033[0m")

		fmt.Print(b.String())
	}
	draw()

	for {
		key, ok := <-arrowChan
		if !ok {
			return options[0]
		}
		switch key {
		case "up":
			selected = (selected - 1 + len(options)) % len(options)
			draw()
		case "down":
			selected = (selected + 1) % len(options)
			draw()
		case "enter":
			return options[selected]
		case "quit":
			return options[0]
		}
	}
}
