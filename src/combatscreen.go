package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

const combatBoxWidth = 60 // largeur intérieure du cadre de narration
const combatBoxLines = 3  // nombre de lignes affichées dans le cadre

var ansiCodeRe = regexp.MustCompile("\x1b\\[[0-9;]*m")

// visibleWidth --> longueur affichée d'une chaîne, sans compter les codes ANSI.
func visibleWidth(s string) int {
	return len([]rune(ansiCodeRe.ReplaceAllString(s, "")))
}

// hpBarColored --> barre de PV : jaune pour la vie restante, rouge sombre pour la vie perdue.
func hpBarColored(current, max, width int) string {
	if max <= 0 {
		max = 1
	}
	if current < 0 {
		current = 0
	}
	filled := current * width / max
	if filled > width {
		filled = width
	}
	return "\033[38;2;235;200;60m" + strings.Repeat("█", filled) +
		"\033[38;2;150;40;40m" + strings.Repeat("░", width-filled) + "\033[0m"
}

// wrapNarration --> découpe un texte en lignes d'au plus `width` caractères visibles (mot par mot).
func wrapNarration(text string, width int) []string {
	words := strings.Fields(text)
	var lines []string
	current := ""
	for _, w := range words {
		candidate := w
		if current != "" {
			candidate = current + " " + w
		}
		if len([]rune(candidate)) > width && current != "" {
			lines = append(lines, current)
			current = w
		} else {
			current = candidate
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

// combatCenterLines --> centre chaque ligne dans la largeur du terminal (ignore les codes ANSI).
func combatCenterLines(lines []string, cols int) []string {
	out := make([]string, len(lines))
	for i, l := range lines {
		pad := (cols - visibleWidth(l)) / 2
		if pad < 0 {
			pad = 0
		}
		out[i] = strings.Repeat(" ", pad) + l
	}
	return out
}

// padVisible --> complète une chaîne avec des espaces jusqu'à `width` colonnes visibles.
func padVisible(s string, width int) string {
	pad := width - visibleWidth(s)
	if pad < 0 {
		pad = 0
	}
	return s + strings.Repeat(" ", pad)
}

// drawCombatFrame --> dessine tout l'écran de combat (sprite, PV ennemi, cadre de narration déjà
// découpé en lignes, PV joueur, boutons), en place (pas de scintillement, pas d'effacement).
func drawCombatFrame(c *Character, m *Monster, narrLines []string, buttons []string, selected int, showButtons bool) {
	cols, _ := terminalSize()

	var lines []string
	lines = append(lines, "", "\033[1m"+m.Name+"\033[0m", "")
	if m.Sprite != "" {
		for _, l := range strings.Split(strings.Trim(m.Sprite, "\n"), "\n") {
			lines = append(lines, l)
		}
	}
	lines = append(lines, "")
	if m.CurrentHP < m.MaxHP {
		lines = append(lines, fmt.Sprintf("PV %s %d/%d", hpBarColored(m.CurrentHP, m.MaxHP, 20), m.CurrentHP, m.MaxHP))
	}
	lines = append(lines, "")

	lines = append(lines, "┌"+strings.Repeat("─", combatBoxWidth)+"┐")
	for i := 0; i < combatBoxLines; i++ {
		text := ""
		if i < len(narrLines) {
			text = narrLines[i]
		}
		lines = append(lines, "│ "+padVisible(text, combatBoxWidth-2)+" │")
	}
	lines = append(lines, "└"+strings.Repeat("─", combatBoxWidth)+"┘")
	lines = append(lines, "")

	status := fmt.Sprintf("%-10s NV %-2d  PV %s %d/%d",
		strings.ToUpper(c.Name), c.Level, hpBarColored(c.CurrentHP, c.MaxHP, 20), c.CurrentHP, c.MaxHP)
	lines = append(lines, status, "")

	if showButtons {
		var parts []string
		for i, btn := range buttons {
			if i == selected {
				parts = append(parts, "\033[1;38;2;255;215;60m> [ "+btn+" ]\033[0m")
			} else {
				parts = append(parts, "\033[38;2;220;140;60m  [ "+btn+" ]\033[0m")
			}
		}
		lines = append(lines, strings.Join(parts, "    "))
	} else {
		lines = append(lines, "")
	}

	centered := combatCenterLines(lines, cols)

	var out strings.Builder
	out.WriteString("\033[H")
	for _, l := range centered {
		out.WriteString(l + "\033[K\r\n")
	}
	fmt.Print(out.String())
}

// showNarration --> affiche le texte du cadre lettre par lettre (~25ms/caractère), en gardant le
// reste de l'écran (sprite, PV, boutons) inchangé. Une touche saute l'animation.
func showNarration(c *Character, m *Monster, text string, buttons []string, selected int, showButtons bool) {
	fullLines := wrapNarration(text, combatBoxWidth-2)
	joined := strings.Join(fullLines, "\n")
	runes := []rune(joined)

	for i := range runes {
		partial := strings.Split(string(runes[:i+1]), "\n")
		drawCombatFrame(c, m, partial, buttons, selected, showButtons)
		select {
		case <-arrowChan:
			drawCombatFrame(c, m, fullLines, buttons, selected, showButtons)
			return
		case <-time.After(25 * time.Millisecond):
		}
	}
	drawCombatFrame(c, m, fullLines, buttons, selected, showButtons)
}

// combatIntro --> phrase d'apparition du monstre, au début du combat.
func combatIntro(c *Character, m *Monster) {
	fmt.Print("\033[H\033[2J")
	showNarration(c, m, "* "+m.Name+" surgit devant toi !", nil, 0, false)
	time.Sleep(600 * time.Millisecond)
}

// combatPlayerChoice --> affiche le choix ATTAQUER/DÉFENDRE et attend la sélection du joueur.
// Tâche 21 : characterTurn --> le tour du joueur, ici en ATTAQUER/DÉFENDRE façon Undertale.
func combatPlayerChoice(c *Character, m *Monster, turn int, telegraph bool) combatAction {
	buttons := []string{"⚔ ATTAQUER", "🛡 DÉFENDRE"}
	selected := 0

	text := "* Que vas-tu faire ?"
	if telegraph {
		text = "* " + m.Name + " prend son élan pour une attaque puissante !"
	}
	showNarration(c, m, text, buttons, selected, true)

	for {
		key, ok := <-arrowChan
		if !ok {
			return actionAttack
		}
		switch key {
		case "left":
			selected = (selected - 1 + len(buttons)) % len(buttons)
			drawCombatFrame(c, m, wrapNarration(text, combatBoxWidth-2), buttons, selected, true)
		case "right":
			selected = (selected + 1) % len(buttons)
			drawCombatFrame(c, m, wrapNarration(text, combatBoxWidth-2), buttons, selected, true)
		case "enter":
			if selected == 0 {
				return actionAttack
			}
			return actionDefend
		case "quit":
			return actionDefend
		}
	}
}

// combatShowPlayerAttack --> affiche le résultat d'une attaque du joueur.
func combatShowPlayerAttack(c *Character, m *Monster, dmg int) {
	text := fmt.Sprintf("* Tu attaques %s : %d dégâts !", m.Name, dmg)
	showNarration(c, m, text, nil, 0, false)
	time.Sleep(400 * time.Millisecond)
}

// combatShowPlayerDefend --> affiche le message de garde.
func combatShowPlayerDefend(c *Character, m *Monster) {
	showNarration(c, m, "* Tu te mets en garde.", nil, 0, false)
	time.Sleep(400 * time.Millisecond)
}

// Tâche 20 : combatShowEnemyAttack --> affiche l'attaque du monstre (reprend goblinPattern : dégâts
// doublés tous les 3 tours).
func combatShowEnemyAttack(c *Character, m *Monster, dmg int, telegraph bool) {
	verb := "attaque"
	if telegraph {
		verb = "frappe de toutes ses forces"
	}
	text := fmt.Sprintf("* %s %s : %d dégâts !", m.Name, verb, dmg)
	showNarration(c, m, text, nil, 0, false)
	time.Sleep(400 * time.Millisecond)
}

// combatVictory --> écran de victoire (EXP, or, montée de niveau éventuelle).
func combatVictory(c *Character, m Monster, leveledUp bool, newLevel int) {
	text := fmt.Sprintf("* Tu as gagné ! Tu obtiens %d EXP et %d pièces d'or.", m.ExpReward, m.GoldReward)
	showNarration(c, &m, text, nil, 0, false)
	time.Sleep(700 * time.Millisecond)
	if leveledUp {
		showNarration(c, &m, fmt.Sprintf("* %s passe au niveau %d !", c.Name, newLevel), nil, 0, false)
		time.Sleep(700 * time.Millisecond)
	}
}

// combatDefeat --> écran de défaite (avant le réveil du Voyageur, Tâche 8).
func combatDefeat(c *Character, m Monster) {
	showNarration(c, &m, "* Tu t'effondres...", nil, 0, false)
	time.Sleep(900 * time.Millisecond)
}
