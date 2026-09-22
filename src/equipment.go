package main

import "fmt"

// Tâche 16 : Equipment --> les emplacements équipables du personnage.
type Equipment struct {
	Head  string
	Chest string
	Feet  string
}

// Tâche 17 : equipItem --> équipe un objet fabriqué, remplace l'ancien et ajuste les PV max.
func (c *Character) equipItem(item string) {
	var slot *string
	var bonus int
	switch item {
	case "Chapeau de l'aventurier":
		slot, bonus = &c.Equip.Head, 10
	case "Tunique de l'aventurier":
		slot, bonus = &c.Equip.Chest, 25
	case "Bottes de l'aventurier":
		slot, bonus = &c.Equip.Feet, 15
	default:
		fmt.Println(item + " ne peut pas être équipé.")
		return
	}

	if !c.removeInventory(item) {
		fmt.Println("Vous n'avez pas " + item + " dans votre inventaire.")
		return
	}

	if *slot != "" {
		c.addInventory(*slot)
		c.MaxHP -= bonus
	}
	*slot = item
	c.MaxHP += bonus
	if c.CurrentHP > c.MaxHP {
		c.CurrentHP = c.MaxHP
	}
	fmt.Printf("%s équipé. PV max : %d\n", item, c.MaxHP)
}
