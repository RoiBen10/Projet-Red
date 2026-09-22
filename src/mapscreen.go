package main

import (
	"fmt"
	"strings"
)

// tileColors --> couleur d'affichage de chaque type de tuile de worldMap.
var tileColors = map[byte][3]int{
	'.': {140, 190, 90},
	',': {170, 165, 160},
	'T': {30, 90, 60},
	'~': {60, 110, 160},
	'#': {70, 65, 65},
	'H': {170, 110, 70},
	'A': {200, 140, 60},
	'F': {150, 60, 50},
	'M': {180, 160, 220},
	'W': {255, 215, 80},
	'K': {90, 170, 180},
	'X': {120, 90, 50},
	'D': {210, 70, 60},
}

// playerColor --> couleur du repère du joueur (bloc vif, bien distinct du terrain).
var playerColor = [3]int{255, 215, 0}

// tileGlyphWidth --> largeur d'une tuile à l'écran, en colonnes (2 espaces de fond coloré, pour
// un rendu à peu près carré dans un terminal).
const tileGlyphWidth = 2

// spawnPoint --> position de départ du Voyageur : le terrain d'entraînement, au sud du village.
func spawnPoint() (int, int) {
	return 50, len(worldMap) - 3
}

// trainingFenceY --> ligne de la porte nord du terrain d'entraînement (au-delà : le village).
const trainingFenceY = 54

// tileAt --> tuile de worldMap à (x, y), ' ' si hors carte.
func tileAt(x, y int) byte {
	if y < 0 || y >= len(worldMap) || x < 0 || x >= len(worldMap[y]) {
		return ' '
	}
	return worldMap[y][x]
}

// attemptMove --> tente de déplacer le joueur vers (nx, ny) : marcher sur un mannequin ('D')
// déclenche un combat d'entraînement au lieu de bloquer ; sinon, avance si la tuile est praticable.
func attemptMove(c *Character, px, py, nx, ny int) (int, int) {
	if tileAt(nx, ny) == 'D' {
		mummy := initMummy()
		disableLineMode()
		trainingFight(c, &mummy)
		return px, py
	}
	if isWalkable(nx, ny) {
		return nx, ny
	}
	return px, py
}

// viewSize --> dimensions de la caméra (en tuiles), calées sur la taille du terminal pour rester
// proche du joueur plutôt que d'afficher toute la carte d'un coup.
func viewSize() (int, int) {
	cols, rows := terminalSize()
	viewW := (cols - 4) / tileGlyphWidth
	viewH := rows - 5 // place pour le nom + l'aide en bas

	if viewW > 31 {
		viewW = 31
	}
	if viewW < 11 {
		viewW = 11
	}
	if viewW%2 == 0 {
		viewW--
	}
	if viewH > 19 {
		viewH = 19
	}
	if viewH < 9 {
		viewH = 9
	}
	if viewH%2 == 0 {
		viewH--
	}
	return viewW, viewH
}

// mapScreen --> affiche une caméra centrée sur le joueur (pas la carte entière), en blocs de
// couleur pleins ; déplacement aux flèches (bloqué par les obstacles), Entrée pour ouvrir le menu.
func mapScreen(c *Character) {
	px, py := spawnPoint()
	mapH := len(worldMap)
	mapW := len(worldMap[0])

	fmt.Print("\033[H\033[2J")
	fmt.Print(c.Name + ", tu te réveilles au terrain d'entraînement, juste avant Emberhollow.\r\n")
	fmt.Print("Utilise les flèches ↑↓←→ pour te déplacer. Rejoins la porte au nord pour entrer dans le village.\r\n\r\n")
	fmt.Print("Appuie sur une touche pour commencer...\r\n")
	if _, ok := <-arrowChan; !ok {
		return
	}
	fmt.Print("\033[H\033[2J")

	draw := func() {
		cols, rows := terminalSize()
		viewW, viewH := viewSize()

		camX := px - viewW/2
		camY := py - viewH/2
		if camX < 0 {
			camX = 0
		}
		if camY < 0 {
			camY = 0
		}
		if camX > mapW-viewW {
			camX = mapW - viewW
		}
		if camY > mapH-viewH {
			camY = mapH - viewH
		}
		if camX < 0 {
			camX = 0
		}
		if camY < 0 {
			camY = 0
		}

		screen := make([][]screenCell, viewH)
		for ty := 0; ty < viewH; ty++ {
			screen[ty] = make([]screenCell, viewW*tileGlyphWidth)
			for tx := 0; tx < viewW; tx++ {
				mx, my := camX+tx, camY+ty
				color := [3]int{0, 0, 0}
				if my >= 0 && my < mapH && mx >= 0 && mx < mapW {
					color = tileColors[worldMap[my][mx]]
				}
				if mx == px && my == py {
					color = playerColor
				}
				for gi := 0; gi < tileGlyphWidth; gi++ {
					screen[ty][tx*tileGlyphWidth+gi] = screenCell{color: color}
				}
			}
		}

		viewOuterW := viewW * tileGlyphWidth
		leftPad := (cols - viewOuterW) / 2
		if leftPad < 0 {
			leftPad = 0
		}
		topPad := (rows - viewH - 2) / 2
		if topPad < 0 {
			topPad = 0
		}
		left := strings.Repeat(" ", leftPad)

		var b strings.Builder
		b.WriteString("\033[H")
		for i := 0; i < topPad; i++ {
			b.WriteString("\r\n")
		}
		b.WriteString(renderScreenBG(screen, left))
		footer := c.Name + "  --  flèches pour se déplacer, Entrée pour le menu"
		if py > trainingFenceY {
			footer = "Terrain d'entraînement -- avance vers le nord (↑) jusqu'à la porte du village"
		}
		b.WriteString("\r\n" + left + footer + "\r\n")
		fmt.Print(b.String())
	}
	draw()

	midnight := make(chan struct{})
	go runClock(midnight)

	for {
		select {
		case key, ok := <-arrowChan:
			if !ok {
				return
			}
			switch key {
			case "up":
				px, py = attemptMove(c, px, py, px, py-1)
			case "down":
				px, py = attemptMove(c, px, py, px, py+1)
			case "left":
				px, py = attemptMove(c, px, py, px-1, py)
			case "right":
				px, py = attemptMove(c, px, py, px+1, py)
			case "enter":
				enableLineMode()
				quit := runMenu(c, midnight)
				disableLineMode()
				if quit {
					return
				}
				fmt.Print("\033[H\033[2J")
			case "quit":
				return
			}
			draw()
		case <-midnight:
			fmt.Print("\033[H\033[2J")
			fmt.Println("Le feu s'abat sur Emberhollow. " + c.Name + " meurt...")
			return
		}
	}
}
