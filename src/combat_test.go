package main

import "testing"

// computeDamage --> ne doit jamais tomber à 0 ou en dessous, même avec une défense écrasante.
func TestComputeDamageFloorsAtOne(t *testing.T) {
	for i := 0; i < 500; i++ {
		if d := computeDamage(1, 50); d < 1 {
			t.Fatalf("computeDamage(1,50) = %d, want >= 1", d)
		}
	}
}

// computeDamage --> reste dans une plage raisonnable (variance ±20%, critique x1.5 à 10%).
func TestComputeDamageReasonableRange(t *testing.T) {
	for i := 0; i < 500; i++ {
		d := computeDamage(10, 2)
		if d < 1 || d > 20 {
			t.Fatalf("computeDamage(10,2) = %d, want roughly within [1,20]", d)
		}
	}
}

// grantRewards --> monte de niveau (PV max/Attaque augmentent) quand l'EXP dépasse le seuil.
func TestGrantRewardsLevelsUp(t *testing.T) {
	c := initCharacter("Test", "Humain", 1, 100, 100, nil)
	c.Exp = 0
	c.ExpMax = 50
	m := Monster{ExpReward: 60, GoldReward: 5}
	leveledUp, newLevel := grantRewards(&c, m)
	if !leveledUp {
		t.Fatalf("expected level up with 60 exp vs expMax 50")
	}
	if newLevel != 2 {
		t.Fatalf("newLevel = %d, want 2", newLevel)
	}
	if c.Gold != 105 { // 100 (or de départ) + 5 (récompense)
		t.Fatalf("Gold = %d, want 105", c.Gold)
	}
	if c.Attack != 12 { // 10 (init) + 2
		t.Fatalf("Attack = %d, want 12", c.Attack)
	}
}

// grantRewards --> pas de montée de niveau tant que l'EXP reste sous le seuil.
func TestGrantRewardsNoLevelUp(t *testing.T) {
	c := initCharacter("Test", "Humain", 1, 100, 100, nil)
	c.Exp = 0
	c.ExpMax = 50
	m := Monster{ExpReward: 20, GoldReward: 3}
	leveledUp, _ := grantRewards(&c, m)
	if leveledUp {
		t.Fatalf("did not expect level up with 20 exp vs expMax 50")
	}
	if c.Exp != 20 {
		t.Fatalf("Exp = %d, want 20", c.Exp)
	}
}
