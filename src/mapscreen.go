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

// defeatedDummies --> mannequins d'entraînement ('D') déjà vaincus : leur tuile redevient de
// l'herbe (praticable, plus de point rouge) au lieu de redéclencher un combat.
var defeatedDummies = map[[2]int]bool{}

// tileAt --> tuile de worldMap à (x, y), ' ' si hors carte. Un mannequin vaincu est rendu comme de
// l'herbe : il disparaît de la carte une fois le combat gagné.
func tileAt(x, y int) byte {
	if y < 0 || y >= len(worldMap) || x < 0 || x >= len(worldMap[y]) {
		return ' '
	}
	tile := worldMap[y][x]
	if tile == 'D' && defeatedDummies[[2]int{x, y}] {
		return '.'
	}
	return tile
}

// attemptMove --> tente de déplacer le joueur vers (nx, ny) : marcher sur un mannequin ('D')
// déclenche un combat d'entraînement au lieu de bloquer ; sinon, avance si la tuile est praticable.
// Modifie directement la position du personnage. Renvoie true si le combat a été perdu (nouvelle
// boucle temporelle déclenchée), pour que l'appelant redémarre l'horloge du jour.
func attemptMove(c *Character, nx, ny int) bool {
	if tileAt(nx, ny) == 'D' {
		mummy := initMummy()
		disableLineMode()
		won := trainingFight(c, &mummy)
		if won {
			defeatedDummies[[2]int{nx, ny}] = true
		}
		// la caméra de la carte ne redessine qu'une petite zone centrée : sans ce nettoyage,
		// le sprite du combat resterait visible en marge après le retour à la carte.
		fmt.Print("\033[H\033[2J")
		return !won
	}
	if isWalkable(nx, ny) {
		c.PosX, c.PosY = nx, ny
	}
	return false
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
	c.PosX, c.PosY = spawnPoint()
	mapH := len(worldMap)
	mapW := len(worldMap[0])

	if !storyIntro(c) {
		return
	}

	draw := func() {
		cols, rows := terminalSize()
		viewW, viewH := viewSize()

		camX := c.PosX - viewW/2
		camY := c.PosY - viewH/2
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
					color = tileColors[tileAt(mx, my)]
				}
				if mx == c.PosX && my == c.PosY {
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
		footer := c.Name + "  --  flèches pour se déplacer, Entrée ou M pour le menu"
		if c.PosY > trainingFenceY {
			footer = "Terrain d'entraînement -- avance vers le nord (↑) jusqu'à la porte du village"
		}
		b.WriteString("\r\n" + left + footer + "\r\n")
		fmt.Print(b.String())
	}
	draw()

	// Une itération de cette boucle = un jour de la boucle temporelle. Une mort (combat perdu,
	// poison, ou minuit) démarre un nouveau jour : le Voyageur revient toujours à la carte, jamais
	// au menu ni à l'écran titre (Tâche 8 / boucle temporelle).
	for {
		midnight := make(chan struct{})
		go runClock(midnight)

		move := func(nx, ny int) bool {
			return attemptMove(c, nx, ny)
		}

		dayEnded := false
		for !dayEnded {
			select {
			case key, ok := <-arrowChan:
				if !ok {
					return
				}
				switch key {
				case "up":
					dayEnded = move(c.PosX, c.PosY-1)
				case "down":
					dayEnded = move(c.PosX, c.PosY+1)
				case "left":
					dayEnded = move(c.PosX-1, c.PosY)
				case "right":
					dayEnded = move(c.PosX+1, c.PosY)
				case "enter", "menu":
					enableLineMode()
					quit, newDay := runMenu(c, midnight)
					disableLineMode()
					if quit {
						return
					}
					fmt.Print("\033[H\033[2J")
					dayEnded = newDay
				case "quit":
					return
				}
				if !dayEnded {
					draw()
				}
			case <-midnight:
				fmt.Print("\033[H\033[2J")
				fmt.Println("Le feu s'abat sur Emberhollow. " + c.Name + " meurt...")
				c.startNewDay()
				fmt.Println("Un nouveau jour commence. " + c.Name + " se réveille au terrain d'entraînement.")
				dayEnded = true
			}
		}
		fmt.Print("\033[H\033[2J")
		draw()
	}
}
