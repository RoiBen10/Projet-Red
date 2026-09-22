package main

import "fmt"

// Tâche 10 : spellBook --> apprend "Boule de Feu" (un sort ne peut être appris qu'une fois).
func (c *Character) spellBook() {
	for _, s := range c.Skills {
		if s == "Boule de Feu" {
			fmt.Println("Vous connaissez déjà Boule de Feu.")
			return
		}
	}
	c.Skills = append(c.Skills, "Boule de Feu")
	fmt.Println("Nouveau sort appris : Boule de Feu")
}

// Tâche 10 : useSpellBook --> consomme le livre de sort et apprend Boule de Feu.
func (c *Character) useSpellBook() {
	if !c.removeInventory("Livre de Sort : Boule de Feu") {
		fmt.Println("Oups, vous n'avez pas ce livre de sort.")
		return
	}
	c.spellBook()
}
