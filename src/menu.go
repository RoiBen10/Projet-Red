package main

import "fmt"

// Tâche 6 / 7 / 15 : runMenu --> affiche le menu (appelé depuis la carte, Entrée pour l'ouvrir).
// L'horloge (midnight) tourne en continu côté carte ; runMenu la surveille aussi pendant qu'il est
// ouvert. Renvoie true si la partie doit se terminer (Quitter ou minuit), false pour revenir à la carte.
func runMenu(c *Character, midnight chan struct{}) bool {
	for {
		fmt.Println("\n=== MENU ===")
		fmt.Println("1. Afficher les informations du personnage")
		fmt.Println("2. Accéder au contenu de l'inventaire")
		fmt.Println("3. Marchand")
		fmt.Println("4. Forgeron")
		fmt.Println("5. Entrainement")
		fmt.Println("6. Retour à la carte")
		fmt.Println("7. Quitter")
		fmt.Print("> ")

		select {
		case choice, ok := <-inputChan:
			if !ok {
				return true
			}
			switch choice {
			case "1":
				c.displayInfo()
			case "2":
				inventoryMenu(c)
			case "3":
				merchantMenu(c)
			case "4":
				forgeMenu(c)
			case "5":
				goblin := initGoblin()
				disableLineMode()
				trainingFight(c, &goblin)
				enableLineMode()
			case "6":
				return false
			case "7":
				fmt.Println("À demain, Voyageur...")
				return true
			default:
				fmt.Println("Choix invalide.")
			}
		case <-midnight:
			fmt.Println("Le feu s'abat sur Emberhollow. " + c.Name + " meurt...")
			return true
		}
	}
}

// Tâche 6 / 5 / 9 / 10 : inventoryMenu --> affiche l'inventaire.
func inventoryMenu(c *Character) {
	for {
		c.accessInventory()
		fmt.Println("1. Utiliser une potion de vie")
		fmt.Println("2. Utiliser une potion de poison")
		fmt.Println("3. Utiliser Livre de Sort : Boule de Feu")
		fmt.Println("4. Équiper Chapeau de l'aventurier")
		fmt.Println("5. Équiper Tunique de l'aventurier")
		fmt.Println("6. Équiper Bottes de l'aventurier")
		fmt.Println("7. Retour")
		fmt.Print("> ")

		switch readChoice() {
		case "1":
			c.takePot()
		case "2":
			c.poisonPot()
		case "3":
			c.useSpellBook()
		case "4":
			c.equipItem("Chapeau de l'aventurier")
		case "5":
			c.equipItem("Tunique de l'aventurier")
		case "6":
			c.equipItem("Bottes de l'aventurier")
		case "7":
			return
		default:
			fmt.Println("Choix invalide.")
		}
	}
}
