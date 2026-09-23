package main

import (
	"math/rand"
)

// combatAction --> action choisie par le joueur pendant son tour.
type combatAction int

const (
	actionAttack combatAction = iota
	actionDefend
)

// Tâche 20 : computeDamage --> dégâts d'un coup (attaque - défense, ±20% aléatoire, 10% de critique x1.5).
func computeDamage(attack, defense int) int {
	base := attack - defense
	if base < 1 {
		base = 1
	}
	variance := 0.8 + rand.Float64()*0.4
	dmg := float64(base) * variance
	if rand.Float64() < 0.10 {
		dmg *= 1.5
	}
	result := int(dmg)
	if result < 1 {
		result = 1
	}
	return result
}

// grantRewards --> ajoute l'expérience/l'or gagnés, gère la montée de niveau (Mission 2).
func grantRewards(c *Character, m Monster) (leveledUp bool, newLevel int) {
	c.Gold += m.GoldReward
	c.Exp += m.ExpReward
	for c.Exp >= c.ExpMax {
		c.Exp -= c.ExpMax
		c.Level++
		c.MaxHP += 10
		c.CurrentHP += 10
		if c.CurrentHP > c.MaxHP {
			c.CurrentHP = c.MaxHP
		}
		c.Attack += 2
		c.ExpMax = c.ExpMax * 3 / 2
		leveledUp = true
	}
	return leveledUp, c.Level
}

// Tâche 22.1 / 22.2 : trainingFight --> combat d'entraînement tour par tour entre le joueur et un
// monstre. Reprend Tâche 20 (dégâts doublés tous les 3 tours, ici annoncés un tour à l'avance) et
// Tâche 21 (tour du joueur : attaquer ou utiliser l'inventaire). Renvoie true si le joueur gagne.
func trainingFight(c *Character, m *Monster) bool {
	turn := 1
	guarding := false

	combatIntro(c, m)

	for {
		telegraph := turn%3 == 0 // toutes les 3 tours : le monstre prépare une attaque puissante

		action := combatPlayerChoice(c, m, turn, telegraph)
		guarding = false

		switch action {
		case actionAttack:
			dmg := computeDamage(c.Attack, m.Defense)
			m.CurrentHP -= dmg
			if m.CurrentHP < 0 {
				m.CurrentHP = 0
			}
			combatShowPlayerAttack(c, m, dmg)
		case actionDefend:
			guarding = true
			heal := 2
			c.CurrentHP += heal
			if c.CurrentHP > c.MaxHP {
				c.CurrentHP = c.MaxHP
			}
			combatShowPlayerDefend(c, m)
		}

		if m.CurrentHP <= 0 {
			leveledUp, newLevel := grantRewards(c, *m)
			combatVictory(c, *m, leveledUp, newLevel)
			return true
		}

		dmg := m.Attack
		if telegraph {
			dmg = m.Attack * 2
		}
		if guarding {
			dmg /= 2
		}
		c.CurrentHP -= dmg
		if c.CurrentHP < 0 {
			c.CurrentHP = 0
		}
		combatShowEnemyAttack(c, m, dmg, telegraph)

		if c.CurrentHP <= 0 {
			combatDefeat(c, *m)
			c.checkDeath() // Tâche 8 : le Voyageur s'effondre, nouvelle boucle
			return false
		}

		turn++
	}
}
