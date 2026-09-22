//go:build windows

package main

import (
	"os"

	"golang.org/x/sys/windows"
	"golang.org/x/term"
)

// enableRawInput --> passe le terminal en mode "cbreak" (ICANON/ECHO désactivés) et active le
// traitement des séquences ANSI (couleurs, positionnement du curseur) sur la sortie, pour que le
// jeu s'affiche correctement dans un terminal Windows moderne (Windows Terminal, PowerShell, cmd).
func enableRawInput() {
	fd := int(os.Stdin.Fd())
	old, err := term.MakeRaw(fd)
	if err != nil {
		return
	}
	savedTermState = old

	out := windows.Handle(os.Stdout.Fd())
	var mode uint32
	if err := windows.GetConsoleMode(out, &mode); err == nil {
		windows.SetConsoleMode(out, mode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING|windows.ENABLE_PROCESSED_OUTPUT)
	}
}
