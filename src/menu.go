package main

import "fmt"

// Tâche 6 / 7 / 15 : runMenu --> affiche le menu (appelé depuis la carte, Entrée ou M pour l'ouvrir).
// L'horloge (midnight) tourne en continu côté carte ; runMenu la surveille aussi pendant qu'il est
// ouvert. Renvoie (quit, newDay) : quit=true si la partie doit se terminer (flux fermé) ; newDay=true
// si une nouvelle boucle temporelle a démarré pendant que le menu était ouvert (mort, minuit), pour
// que la carte redémarre son horloge. Après un combat (gagné ou perdu), on revient toujours
// directement à la carte, jamais au menu.
func runMenu(c *Character, midnight chan struct{}) (quit bool, newDay bool) {
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
				return true, false
			}
			switch choice {
			case "1":
				c.displayInfo()
			case "2":
				if inventoryMenu(c) {
					return false, true
				}
			case "3":
				merchantMenu(c)
			case "4":
				forgeMenu(c)
			case "5":
				goblin := initGoblin()
				disableLineMode()
				won := trainingFight(c, &goblin)
				enableLineMode()
				return false, !won
			case "6":
				return false, false
			case "7":
				fmt.Println("À demain, Voyageur...")
				return true, false
			default:
				fmt.Println("Choix invalide.")
			}
		case <-midnight:
			fmt.Println("Le feu s'abat sur Emberhollow. " + c.Name + " meurt...")
			c.startNewDay()
			return false, true
		}
	}
}

// Tâche 6 / 5 / 9 / 10 : inventoryMenu --> affiche l'inventaire. Renvoie true si une potion de
// poison a tué le Voyageur (nouvelle boucle déclenchée), pour que runMenu referme tout et revienne
// directement à la carte.
func inventoryMenu(c *Character) bool {
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
			if c.poisonPot() {
				return true
			}
		case "3":
			c.useSpellBook()
		case "4":
			c.equipItem("Chapeau de l'aventurier")
		case "5":
			c.equipItem("Tunique de l'aventurier")
		case "6":
			c.equipItem("Bottes de l'aventurier")
		case "7":
			return false
		default:
			fmt.Println("Choix invalide.")
		}
	}
}
