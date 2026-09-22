//go:build !windows

package main

import (
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/term"
)

// enableRawInput --> passe le terminal en mode "cbreak" (ICANON/ECHO désactivés) tout en réactivant
// OPOST ensuite, pour garder la conversion \n -> \r\n automatique en sortie (macOS/Linux).
func enableRawInput() {
	fd := int(os.Stdin.Fd())
	old, err := term.MakeRaw(fd)
	if err != nil {
		return
	}
	savedTermState = old

	var t syscall.Termios
	syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.TIOCGETA, uintptr(unsafe.Pointer(&t)))
	t.Oflag |= syscall.OPOST
	syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.TIOCSETA, uintptr(unsafe.Pointer(&t)))
}
