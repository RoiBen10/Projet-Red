package main

import (
	"fmt"
	"strings"
	"time"
)

// portraitPixels --> portrait du Voyageur, converti en pixel-art (grille 10x8, façon icône).
// '.' = transparent, 'D' = cheveux (foncé), 'M' = cheveux (mèche claire), 'K' = mèche latérale,
// 'S' = peau, 'E' = sourcil/œil, 'O' = bouche.
var portraitPixels = []string{
	"..DDDDDD..",
	".DDMMMMDD.",
	".DSSSSSSD.",
	"KDSESSESD.",
	"KDSSSSSSD.",
	".DSESSESD.",
	".DSSOOSSD.",
	"..DDDDDD..",
}

// portraitColors --> couleurs (RGB) associées aux lettres de portraitPixels.
var portraitColors = map[byte][3]int{
	'D': {59, 34, 20},
	'M': {94, 58, 33},
	'K': {168, 104, 58},
	'S': {247, 184, 141},
	'E': {58, 36, 24},
	'O': {201, 122, 69},
}

// portraitCols / portraitDisplayWidth --> dimensions du portrait affiché (2 colonnes par pixel).
const portraitCols = 10
const portraitDisplayWidth = portraitCols * 2

// portraitRow --> une ligne du portrait, en blocs de couleur de fond (transparent = fond par défaut).
func portraitRow(y int) string {
	var b strings.Builder
	lastColor := [3]int{-1, -1, -1}
	transparent := true
	for x := 0; x < portraitCols; x++ {
		ch := byte('.')
		if y < len(portraitPixels) && x < len(portraitPixels[y]) {
			ch = portraitPixels[y][x]
		}
		if ch == '.' {
			if !transparent {
				b.WriteString("\033[49m")
				transparent = true
			}
			b.WriteString("  ")
			continue
		}
		color := portraitColors[ch]
		if transparent || color != lastColor {
			fmt.Fprintf(&b, "\033[48;2;%d;%d;%dm", color[0], color[1], color[2])
			lastColor, transparent = color, false
		}
		b.WriteString("  ")
	}
	b.WriteString("\033[0m")
	return b.String()
}

// villageColor --> teinte unique (clair de lune) pour tout le fond ASCII du village : "sans
// couleur" comme demandé, la densité des caractères fait le contraste, pas une palette.
var villageColor = [3]int{150, 160, 175}

// buildMonoScreen --> grille targetW x targetH à partir de caractères ASCII seuls, tous rendus
// dans une teinte unique color (contrairement à buildScreen, qui utilise une couleur par pixel).
func buildMonoScreen(chars []string, color [3]int, targetW, targetH int) [][]screenCell {
	screen := make([][]screenCell, targetH)
	for y := 0; y < targetH; y++ {
		screen[y] = make([]screenCell, targetW)
		for x := 0; x < targetW; x++ {
			cell := screenCell{ch: ' ', color: color}
			if y < len(chars) && x < len(chars[y]) {
				cell.ch = rune(chars[y][x])
			}
			screen[y][x] = cell
		}
	}
	return screen
}

// villageFrameIndex --> avance/recule dans villageFrames façon va-et-vient (comme l'écran titre),
// pour une boucle fluide sans saut brusque entre la dernière et la première image.
func villageFrameIndex(step int) int {
	n := len(villageFrames)
	if n == 0 {
		return 0
	}
	cycle := 2 * (n - 1)
	if cycle <= 0 {
		return 0
	}
	pos := step % cycle
	if pos < 0 {
		pos += cycle
	}
	if pos < n {
		return pos
	}
	return cycle - pos
}

// drawIntroFrame --> le fond animé (Emberhollow, de nuit) et la boîte de dialogue par-dessus, dans
// le même rendu : un seul repositionnement de curseur, pas d'effacement, donc pas de scintillement.
func drawIntroFrame(step int, lines []string, showPrompt bool, showPortrait bool) {
	if len(villageFrames) > 0 {
		cols, rows := terminalSize()
		chars := resampleGrid(villageFrames[villageFrameIndex(step)], cols, rows)
		screen := buildMonoScreen(chars, villageColor, cols, rows)
		dimRegion(screen, 0, 0, cols, rows, 0.5) // fond discret, pour ne pas gêner la lecture
		fmt.Print("\033[H" + renderScreen(screen))
	}
	drawDialogueBox(lines, showPrompt, showPortrait)
}

// dialogueTextWidth / dialogueBoxHeight --> dimensions du texte et hauteur totale de la boîte
// (au moins aussi haute que le portrait, pour que les deux cadres s'alignent).
func dialogueTextWidth() int {
	cols, _ := terminalSize()
	w := 56
	if maxW := cols - portraitDisplayWidth - 10; w > maxW {
		w = maxW
	}
	if w < 20 {
		w = 20
	}
	return w
}

// dialogueBoxHeight --> hauteur intérieure de la boîte : la hauteur du portrait quand il est
// affiché, sinon une petite boîte de narration classique (pas d'espace vide inutile).
func dialogueBoxHeight(showPortrait bool) int {
	if showPortrait && len(portraitPixels) > combatBoxLines {
		return len(portraitPixels)
	}
	return combatBoxLines
}

// drawDialogueBox --> boîte de dialogue ancrée en bas de l'écran, façon Undertale : portrait à
// gauche (uniquement si showPortrait, c'est-à-dire quand c'est notre personnage qui parle) et
// texte à droite. showPrompt affiche un repère "▼" en bas à droite une fois le texte affiché.
func drawDialogueBox(lines []string, showPrompt bool, showPortrait bool) {
	cols, rows := terminalSize()
	textWidth := dialogueTextWidth()
	boxHeight := dialogueBoxHeight(showPortrait)

	portraitBlockWidth := 0
	if showPortrait {
		portraitBlockWidth = portraitDisplayWidth + 3 // "portrait │ "
	}
	innerWidth := portraitBlockWidth + textWidth
	leftPad := (cols - (innerWidth + 4)) / 2
	if leftPad < 0 {
		leftPad = 0
	}
	left := strings.Repeat(" ", leftPad)

	topPad := rows - (boxHeight + 4)
	if topPad < 1 {
		topPad = 1
	}

	var b strings.Builder
	b.WriteString("\033[H\033[49m\033[37m") // \033[49m : annule le fond hérité du décor animé
	for i := 0; i < topPad; i++ {
		b.WriteString("\r\n")
	}
	b.WriteString(left + "┌" + strings.Repeat("─", innerWidth+2) + "┐\r\n")

	textRow := boxHeight - 1 // ligne du repère "▼", en bas de la boîte
	for i := 0; i < boxHeight; i++ {
		var row strings.Builder
		row.WriteString(left + "│ ")
		if showPortrait {
			row.WriteString(portraitRow(i) + "\033[49m\033[37m │ ")
		}
		text := ""
		if i < len(lines) {
			text = lines[i]
		}
		if i == textRow && showPrompt {
			text = padRightText(text, textWidth-2) + " ▼"
		} else {
			text = padRightText(text, textWidth)
		}
		row.WriteString(text + " │\r\n")
		b.WriteString(row.String())
	}
	b.WriteString(left + "└" + strings.Repeat("─", innerWidth+2) + "┘\r\n")
	b.WriteString("\033[0m")
	fmt.Print(b.String())
}

// showDialogue --> affiche un texte lettre par lettre dans une boîte de dialogue façon Undertale.
// Une touche affiche le texte en entier d'un coup, une touche suivante fait avancer. speaker=true
// affiche le portrait du Voyageur à gauche (quand c'est notre personnage qui parle), false pour une
// simple narration (pas de portrait). Renvoie false si le flux d'entrée s'est fermé.
func showDialogue(text string, speaker bool, startFrame int) bool {
	fullLines := wrapNarration(text, dialogueTextWidth())
	joined := strings.Join(fullLines, "\n")
	runes := []rune(joined)

	for i := range runes {
		partial := strings.Split(string(runes[:i+1]), "\n")
		drawIntroFrame(startFrame+i/5, partial, false, speaker)
		select {
		case _, ok := <-arrowChan:
			if !ok {
				return false
			}
			drawIntroFrame(startFrame+i/5, fullLines, true, speaker)
			_, ok = <-arrowChan
			return ok
		case <-time.After(20 * time.Millisecond):
		}
	}
	drawIntroFrame(startFrame+len(runes)/5, fullLines, true, speaker)
	_, ok := <-arrowChan
	return ok
}

// storyIntro --> l'intrigue du jeu, présentée façon Undertale (boîtes de dialogue successives),
// juste avant l'arrivée sur la carte. Renvoie false si le flux d'entrée s'est fermé (jeu à quitter).
func storyIntro(c *Character) bool {
	type beat struct {
		text    string
		speaker bool
	}
	beats := []beat{
		{"Un voyageur arrive un matin à Emberhollow, un village tranquille perdu dans les collines.", false},
		{"Chaque nuit, à minuit, la cloche du temple sonne son dernier coup. Et le feu s'abat sur le village.", false},
		{"Pourtant, chaque matin, tout recommence... comme si rien ne s'était passé.", false},
		{"Où... où suis-je ? Ce village... j'ai l'impression de l'avoir déjà vu.", true},
		{"Il faut que je comprenne ce qui se passe ici, avant que la nuit ne tombe à nouveau.", true},
		{c.Name + ", tu te réveilles au terrain d'entraînement, juste avant Emberhollow.", false},
		{"Utilise les flèches ↑↓←→ pour te déplacer. Entrée ou M ouvrent le menu à tout moment.", false},
	}

	for i, beat := range beats {
		if !showDialogue(beat.text, beat.speaker, i*25) {
			return false
		}
	}
	fmt.Print("\033[H\033[2J")
	return true
}
