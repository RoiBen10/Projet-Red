package main

import "testing"

// checkDeath --> à 0 PV, déclenche une nouvelle boucle : PV pleins, inventaire et or perdus,
// retour au point de départ, mais le niveau (ce que le Voyageur a appris) est conservé.
func TestCheckDeathResetsLoop(t *testing.T) {
	c := initCharacter("Test", "Humain", 3, 100, 0, []string{"Potion de vie"})
	c.Gold = 42
	c.PosX, c.PosY = 12, 34
	died := c.checkDeath()
	if !died {
		t.Fatalf("checkDeath() = false at 0 HP, want true")
	}
	if c.CurrentHP != c.MaxHP {
		t.Fatalf("CurrentHP = %d, want full heal to %d", c.CurrentHP, c.MaxHP)
	}
	if len(c.Inventory) != 0 {
		t.Fatalf("Inventory not cleared: %v", c.Inventory)
	}
	if c.Gold != 0 {
		t.Fatalf("Gold = %d, want 0 after death", c.Gold)
	}
	wantX, wantY := spawnPoint()
	if c.PosX != wantX || c.PosY != wantY {
		t.Fatalf("position = (%d,%d), want spawn (%d,%d)", c.PosX, c.PosY, wantX, wantY)
	}
	if c.Level != 3 {
		t.Fatalf("Level = %d, want kept at 3 (only inventory/gold reset)", c.Level)
	}
}

// checkDeath --> ne doit rien changer tant que le personnage est encore en vie.
func TestCheckDeathNoopWhenAlive(t *testing.T) {
	c := initCharacter("Test", "Humain", 1, 100, 50, []string{"Potion de vie"})
	if c.checkDeath() {
		t.Fatalf("checkDeath() = true at 50 HP, want false")
	}
	if c.CurrentHP != 50 || len(c.Inventory) != 1 {
		t.Fatalf("state mutated despite being alive: HP=%d inv=%v", c.CurrentHP, c.Inventory)
	}
}
