package main

import "fmt"

func main() {
	startInputReader()
	defer restoreInput()

	if !startScreen() {
		fmt.Println("À bientôt, Voyageur...")
		return
	}

	class := classSelectScreen()
	name := nameSelectScreen()

	c1 := characterCreation(name, class)
	mapScreen(&c1)
}
