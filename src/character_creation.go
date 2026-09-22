package main

import "strings"

// Tâche 11 : formatName --> majuscule au début, minuscules ensuite (la grille de nom ne fournit
// que des lettres A-Z, la contrainte "uniquement des lettres" est donc garantie à la construction).
func formatName(s string) string {
	s = strings.ToLower(s)
	return strings.ToUpper(s[:1]) + s[1:]
}

// Tâche 11 : classMaxHP --> PV max de départ selon la classe.
func classMaxHP(class string) int {
	switch class {
	case "Elfe":
		return 80
	case "Nain":
		return 120
	default: // Humain
		return 100
	}
}

// Tâche 11 : characterCreation --> construit le personnage (nom et classe déjà choisis, via
// nameSelectScreen et classSelectScreen, avant d'arriver ici).
func characterCreation(name, class string) Character {
	name = formatName(name)
	maxHP := classMaxHP(class)
	currentHP := maxHP / 2
	return initCharacter(name, class, 1, maxHP, currentHP, []string{"Potion de vie", "Potion de vie", "Potion de vie"})
}
