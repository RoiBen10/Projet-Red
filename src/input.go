package main

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"golang.org/x/term"
)

var (
	inputChan       chan string // lignes complètes (menus numérotés, saisie de texte)
	arrowChan       chan string // "up" / "down" / "enter" / "quit" (écrans animés aux flèches)
	lineModeFlag    int32       // 0 = mode menu (flèches), 1 = mode ligne (texte)
	inputReaderOnce sync.Once
	savedTermState  *term.State
)

// enableLineMode --> bascule le lecteur clavier en mode "ligne" : Entrée valide un texte tapé
// (menus numérotés, saisie de nom) plutôt qu'un simple appui pour un écran animé aux flèches.
func enableLineMode() {
	atomic.StoreInt32(&lineModeFlag, 1)
}

// disableLineMode --> repasse en mode "menu" (flèches + Entrée envoyées sur arrowChan), le temps
// d'un écran aux flèches (ex. le combat) appelé depuis un contexte en mode ligne (ex. le menu).
func disableLineMode() {
	atomic.StoreInt32(&lineModeFlag, 0)
}

// startInputReader --> démarre l'unique lecteur clavier du jeu, une seule fois pour toute sa durée.
// Mode "cbreak" : touches lues immédiatement (pas besoin d'attendre Entrée pour les flèches), écho
// géré à la main, mais la conversion \n -> \r\n reste automatique en sortie (OPOST conservé), donc
// tout le reste du jeu (fmt.Println, etc.) continue de s'afficher normalement.
func startInputReader() {
	inputReaderOnce.Do(func() {
		inputChan = make(chan string, 4)
		arrowChan = make(chan string, 4)

		enableRawInput()

		go func() {
			defer close(inputChan)
			defer close(arrowChan)
			var line []rune
			buf := make([]byte, 1)
			for {
				n, err := os.Stdin.Read(buf)
				if err != nil || n == 0 {
					return
				}
				b := buf[0]
				switch {
				case b == 27: // ESC --> potentielle séquence flèche (ESC [ A/B)
					var seq [2]byte
					if n1, _ := os.Stdin.Read(seq[:1]); n1 > 0 && seq[0] == '[' {
						if n2, _ := os.Stdin.Read(seq[1:2]); n2 > 0 {
							switch seq[1] {
							case 'A':
								arrowChan <- "up"
							case 'B':
								arrowChan <- "down"
							case 'C':
								arrowChan <- "right"
							case 'D':
								arrowChan <- "left"
							}
						}
					}
				case b == '\r' || b == '\n':
					if atomic.LoadInt32(&lineModeFlag) == 1 {
						fmt.Print("\n")
						inputChan <- string(line)
						line = nil
					} else {
						arrowChan <- "enter"
					}
				case b == 127 || b == 8: // retour arrière
					if atomic.LoadInt32(&lineModeFlag) == 1 && len(line) > 0 {
						line = line[:len(line)-1]
						fmt.Print("\b \b")
					}
				case b == 3: // Ctrl+C
					arrowChan <- "quit"
				default:
					if atomic.LoadInt32(&lineModeFlag) == 1 {
						r := rune(b)
						line = append(line, r)
						fmt.Print(string(r))
					}
				}
			}
		}()
	})
}

// Tâche 6 : readChoice --> attend la prochaine ligne tapée par le joueur.
func readChoice() string {
	choice, ok := <-inputChan
	if !ok {
		return "4"
	}
	return choice
}

// restoreInput --> remet le terminal dans son état d'origine (à appeler à la sortie du jeu).
func restoreInput() {
	if savedTermState != nil {
		term.Restore(int(os.Stdin.Fd()), savedTermState)
	}
}
