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

// Tâche 8 : isDead --> à 0 PV, le perso meurt puis revient à 50% de ses PV max.
func (c *Character) isDead() bool {
	if c.CurrentHP > 0 {
		return false
	}
	fmt.Println(c.Name + " s'effondre... et se réveille dans sa chambre, comme chaque matin.")
	c.CurrentHP = c.MaxHP / 2
	fmt.Printf("PV : %d/%d\n", c.CurrentHP, c.MaxHP)
	return true
}
