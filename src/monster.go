package main

// Tâche 19 : Monster --> struct pour les adversaires de combat.
type Monster struct {
	Name       string
	MaxHP      int
	CurrentHP  int
	Attack     int
	Defense    int
	ExpReward  int
	GoldReward int
	Sprite     string
}

// goblinSprite --> ASCII, gobelin d'entraînement.
const goblinSprite = `
      /\_/\
     ( o.o )
    > ^   ^ <
   /|  \_/  |\
  ' |  | |  | `

// Tâche 19 : initGoblin --> initialise un Gobelin d'entraînement.
func initGoblin() Monster {
	return Monster{
		Name:       "Gobelin d'entrainement",
		MaxHP:      40,
		CurrentHP:  40,
		Attack:     5,
		Defense:    1,
		ExpReward:  20,
		GoldReward: 5,
		Sprite:     goblinSprite,
	}
}

// mummySprite --> ASCII, mannequin d'entraînement du terrain (momie bandée, yeux rouges).
const mummySprite = `
     .:=======:.
   ,'  ::::::::  ` + "`" + `.
  /   (@)    (@)   \
 |                   |
 |      ~~~~~~~      |
  \  :::::::::::::  /
   ` + "`" + `._____________.'
      |=========|
      |=========|----o
     /|=========|\
    ' |=========| ` + "`" + `
      |=========|
     /|    |    |\
    ' |    |    | ` + "`" + `
      '----'----'
`

// initMummy --> initialise une Momie d'entraînement (mannequin du terrain, Tâche 19-like).
func initMummy() Monster {
	return Monster{
		Name:       "Momie d'entrainement",
		MaxHP:      30,
		CurrentHP:  30,
		Attack:     4,
		Defense:    0,
		ExpReward:  15,
		GoldReward: 3,
		Sprite:     mummySprite,
	}
}
