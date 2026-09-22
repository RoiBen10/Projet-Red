package main

import "fmt"

// Tâche 14 : buyItem --> achète un item si le joueur a assez d'or et de place dans l'inventaire.
func (c *Character) buyItem(item string, cost int) {
	if c.Gold < cost {
		fmt.Printf("Pas assez d'or pour %s (%d pièce(s) requises, vous en avez %d).\n", item, cost, c.Gold)
		return
	}
	if !c.addInventory(item) {
		return
	}
	c.Gold -= cost
	fmt.Printf("Achat : %s (-%d pièce(s) d'or)\n", item, cost)
}

// Tâche 7 / 9 / 10 / 14 : merchantMenu --> interface du marchand.
func merchantMenu(c *Character) {
	for {
		fmt.Println("\n=== Marchand ===")
		fmt.Println("1. Potion de vie (3 or)")
		fmt.Println("2. Potion de poison (6 or)")
		fmt.Println("3. Livre de Sort : Boule de Feu (25 or)")
		fmt.Println("4. Fourrure de Loup (4 or)")
		fmt.Println("5. Peau de Troll (7 or)")
		fmt.Println("6. Cuir de Sanglier (3 or)")
		fmt.Println("7. Plume de Corbeau (1 or)")
		fmt.Println("8. Augmentation d'inventaire (30 or)")
		fmt.Println("9. Retour")
		fmt.Print("> ")

		switch readChoice() {
		case "1":
			c.buyItem("Potion de vie", 3)
		case "2":
			c.buyItem("Potion de poison", 6)
		case "3":
			c.buyItem("Livre de Sort : Boule de Feu", 25)
		case "4":
			c.buyItem("Fourrure de Loup", 4)
		case "5":
			c.buyItem("Peau de Troll", 7)
		case "6":
			c.buyItem("Cuir de Sanglier", 3)
		case "7":
			c.buyItem("Plume de Corbeau", 1)
		case "8":
			c.buyInventoryUpgrade()
		case "9":
			return
		default:
			fmt.Println("Choix invalide.")
		}
	}
}
