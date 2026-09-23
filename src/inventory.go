package main

import (
	"fmt"
	"time"
)

// Tâche 4 : accessInventory --> liste les objets de l'inventaire.
func (c Character) accessInventory() {
	fmt.Println("=== Inventaire de " + c.Name + " ===")
	if len(c.Inventory) == 0 {
		fmt.Println("(vide)")
		return
	}
	for i, item := range c.Inventory {
		fmt.Printf("%d. %s\n", i+1, item)
	}
}

// Tâche 7 (astuce) / Tâche 12 : addInventory --> ajoute un item si l'inventaire n'est pas plein.
func (c *Character) addInventory(item string) bool {
	if len(c.Inventory) >= c.InventoryCapacity {
		fmt.Printf("Inventaire plein (%d/%d), impossible d'ajouter %s.\n", len(c.Inventory), c.InventoryCapacity, item)
		return false
	}
	c.Inventory = append(c.Inventory, item)
	return true
}

// Tâche 18 : upgradeInventorySlot --> +10 capacité, utilisable 3 fois maximum.
func (c *Character) upgradeInventorySlot() bool {
	if c.InventoryUpgrades >= 3 {
		fmt.Println("Vous avez atteint la limite maximale d'amélioration (3/3).")
		return false
	}
	c.InventoryCapacity += 10
	c.InventoryUpgrades++
	fmt.Printf("Inventaire agrandi ! Nouvelle capacité : %d (Amélioration %d/3)\n", c.InventoryCapacity, c.InventoryUpgrades)
	return true
}

// Tâche 18 : buyInventoryUpgrade --> achète une augmentation d'inventaire chez le marchand (30 or).
func (c *Character) buyInventoryUpgrade() {
	const cost = 30
	if c.InventoryUpgrades >= 3 {
		fmt.Println("Vous ne pouvez plus acheter cette amélioration (limite atteinte).")
		return
	}
	if c.Gold < cost {
		fmt.Printf("Pas assez d'or ! Il vous faut %d pièces d'or.\n", cost)
		return
	}
	c.Gold -= cost
	c.upgradeInventorySlot()
}

// Tâche 15 : countItem --> compte les occurrences d'un item dans l'inventaire (utilisé par le forgeron).
func (c Character) countItem(item string) int {
	count := 0
	for _, it := range c.Inventory {
		if it == item {
			count++
		}
	}
	return count
}

// Tâche 7 (astuce) : removeInventory --> retire la 1ère occurrence de l'item, false si absent.
func (c *Character) removeInventory(item string) bool {
	for i, it := range c.Inventory {
		if it == item {
			c.Inventory = append(c.Inventory[:i], c.Inventory[i+1:]...)
			return true
		}
	}
	return false
}

// Tâche 5 : takePot --> consomme 1 potion de vie, rend 50 PV.
func (c *Character) takePot() {
	if !c.removeInventory("Potion de vie") {
		fmt.Println("Aucune potion de vie dans l'inventaire.")
		return
	}
	c.CurrentHP += 50
	if c.CurrentHP > c.MaxHP {
		c.CurrentHP = c.MaxHP
	}
	fmt.Printf("Potion de vie utilisée. PV : %d/%d\n", c.CurrentHP, c.MaxHP)
}

// Tâche 9 : poisonPot --> inflige 10 dégâts/seconde pendant 3s. Renvoie true si le poison a tué le
// Voyageur (nouvelle boucle déclenchée), pour que l'appelant referme le menu et revienne à la carte.
func (c *Character) poisonPot() bool {
	if !c.removeInventory("Potion de poison") {
		fmt.Println("Aucune potion de poison dans l'inventaire.")
		return false
	}
	fmt.Println("Vous buvez la potion de poison...")
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second)
		c.CurrentHP -= 10
		if c.CurrentHP < 0 {
			c.CurrentHP = 0
		}
		fmt.Printf("PV : %d/%d\n", c.CurrentHP, c.MaxHP)
		if c.checkDeath() {
			return true
		}
	}
	return false
}
