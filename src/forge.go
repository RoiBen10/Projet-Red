package main

import "fmt"

// Tâche 15 : craftItem --> fabrique un équipement si le joueur a l'or, les ressources et la place.
func (c *Character) craftItem(name string, cost int, materials map[string]int) {
	if c.Gold < cost {
		fmt.Printf("Pas assez d'or pour fabriquer %s (%d pièce(s) requises, vous en avez %d).\n", name, cost, c.Gold)
		return
	}
	for material, qty := range materials {
		if c.countItem(material) < qty {
			fmt.Printf("Il vous manque des ressources pour fabriquer %s (%d %s requis).\n", name, qty, material)
			return
		}
	}

	qtyTotal := 0
	for _, qty := range materials {
		qtyTotal += qty
	}
	if len(c.Inventory)-qtyTotal+1 > c.InventoryCapacity {
		fmt.Println("Pas assez de place dans l'inventaire pour " + name + ".")
		return
	}

	for material, qty := range materials {
		for i := 0; i < qty; i++ {
			c.removeInventory(material)
		}
	}
	c.addInventory(name)
	c.Gold -= cost
	fmt.Printf("Fabriqué : %s (-%d pièce(s) d'or)\n", name, cost)
}

// Tâche 15 : forgeMenu --> interface du forgeron.
func forgeMenu(c *Character) {
	for {
		fmt.Println("\n=== Forgeron ===")
		fmt.Println("1. Chapeau de l'aventurier (5 or)")
		fmt.Println("2. Tunique de l'aventurier (5 or)")
		fmt.Println("3. Bottes de l'aventurier (5 or)")
		fmt.Println("4. Retour")
		fmt.Print("> ")

		switch readChoice() {
		case "1":
			c.craftItem("Chapeau de l'aventurier", 5, map[string]int{
				"Plume de Corbeau": 1,
				"Cuir de Sanglier": 1,
			})
		case "2":
			c.craftItem("Tunique de l'aventurier", 5, map[string]int{
				"Fourrure de Loup": 2,
				"Peau de Troll":    1,
			})
		case "3":
			c.craftItem("Bottes de l'aventurier", 5, map[string]int{
				"Fourrure de Loup": 1,
				"Cuir de Sanglier": 1,
			})
		case "4":
			return
		default:
			fmt.Println("Choix invalide.")
		}
	}
}
