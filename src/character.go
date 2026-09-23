package main

import "fmt"

// Tâche 1 : Character = le Voyageur.
type Character struct {
	Name              string
	Class             string
	Level             int
	MaxHP             int
	CurrentHP         int
	Inventory         []string
	InventoryCapacity int
	InventoryUpgrades int
	Skills            []string
	Gold              int
	Equip             Equipment
	Attack            int
	Defense           int
	Exp               int
	ExpMax            int
	PosX              int
	PosY              int
}

// Tâche 2 : initCharacter --> construit un perso, avec "Coup de poing" comme sort de base et 100 pièces d'or.
func initCharacter(name, class string, level, maxHP, currentHP int, inventory []string) Character {
	return Character{
		Name:              name,
		Class:             class,
		Level:             level,
		MaxHP:             maxHP,
		CurrentHP:         currentHP,
		Inventory:         inventory,
		InventoryCapacity: 10,
		Skills:            []string{"Coup de poing"},
		Gold:              100,
		Attack:            10,
		Defense:           2,
		ExpMax:            50,
	}
}

// Tâche 3 : displayInfo --> affiche la fiche du perso.
func (c Character) displayInfo() {
	fmt.Println("=== " + c.Name + " ===")
	fmt.Printf("Classe     : %s\n", c.Class)
	fmt.Printf("Niveau     : %d\n", c.Level)
	fmt.Printf("PV         : %d/%d\n", c.CurrentHP, c.MaxHP)
	fmt.Printf("Or         : %d pièce(s)\n", c.Gold)
	fmt.Printf("Inventaire : %d objet(s)\n", len(c.Inventory))
}

// startNewDay --> la boucle temporelle recommence (minuit, ou le Voyageur s'effondre) : PV pleins,
// inventaire et or perdus, retour au terrain d'entraînement. Il garde ce qu'il a appris (niveau, sorts).
func (c *Character) startNewDay() {
	c.CurrentHP = c.MaxHP
	c.Inventory = nil
	c.Gold = 0
	c.PosX, c.PosY = spawnPoint()
}

// Tâche 8 : checkDeath --> à 0 PV (combat perdu, poison...), déclenche une nouvelle boucle et le
// signale à l'appelant, pour qu'il referme tout écran ouvert et revienne directement à la carte.
func (c *Character) checkDeath() bool {
	if c.CurrentHP > 0 {
		return false
	}
	fmt.Println(c.Name + " s'effondre... Le feu s'abat sur Emberhollow, et tout recommence au matin.")
	c.startNewDay()
	fmt.Printf("PV restaurés : %d/%d. Inventaire et or perdus, mais tu gardes ce que tu as appris.\n", c.CurrentHP, c.MaxHP)
	return true
}
