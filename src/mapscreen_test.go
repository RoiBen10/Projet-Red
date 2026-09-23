package main

import "testing"

// tileAt --> renvoie un espace (hors carte) plutôt que de paniquer, quelles que soient les coordonnées.
func TestTileAtOutOfBounds(t *testing.T) {
	if tileAt(-1, -1) != ' ' {
		t.Fatalf("tileAt out of bounds should be space")
	}
	if tileAt(1000000, 1000000) != ' ' {
		t.Fatalf("tileAt far out of bounds should be space")
	}
}

// spawnPoint --> doit toujours pointer sur une case praticable du terrain d'entraînement.
func TestSpawnPointIsWalkable(t *testing.T) {
	x, y := spawnPoint()
	if !isWalkable(x, y) {
		t.Fatalf("spawn point (%d,%d) is not walkable, tile=%q", x, y, tileAt(x, y))
	}
}

// attemptMove --> ne déplace jamais le joueur sur une tuile non praticable (mur, eau, etc.).
func TestAttemptMoveBlockedByWall(t *testing.T) {
	c := &Character{}
	c.PosX, c.PosY = spawnPoint()
	startX, startY := c.PosX, c.PosY
	// première ligne de la carte : uniquement des arbres, non praticable.
	died := attemptMove(c, c.PosX, 0)
	if died {
		t.Fatalf("attemptMove into a wall should never report a death")
	}
	if c.PosX != startX || c.PosY != startY {
		t.Fatalf("position moved into an unwalkable tile: now (%d,%d)", c.PosX, c.PosY)
	}
}

// attemptMove --> avance normalement sur une case praticable du terrain d'entraînement.
func TestAttemptMoveOntoGrass(t *testing.T) {
	c := &Character{}
	c.PosX, c.PosY = spawnPoint()
	_, wantY := spawnPoint()
	wantY--
	died := attemptMove(c, c.PosX, c.PosY-1)
	if died {
		t.Fatalf("moving onto open grass should never report a death")
	}
	if c.PosY != wantY {
		t.Fatalf("player did not move: still at y=%d, want %d", c.PosY, wantY)
	}
}

// tileAt / isWalkable --> un mannequin ('D') vaincu redevient de l'herbe (praticable, plus de
// point rouge à l'écran) au lieu de redéclencher un combat.
func TestDefeatedDummyDisappears(t *testing.T) {
	x, y := 40, 60 // position d'un mannequin, worldmap.go
	if tileAt(x, y) != 'D' {
		t.Fatalf("expected a training dummy at (%d,%d), got %q", x, y, tileAt(x, y))
	}
	if isWalkable(x, y) {
		t.Fatalf("an active dummy should block movement")
	}

	defeatedDummies[[2]int{x, y}] = true
	defer delete(defeatedDummies, [2]int{x, y})

	if tileAt(x, y) != '.' {
		t.Fatalf("defeated dummy tile = %q, want grass '.'", tileAt(x, y))
	}
	if !isWalkable(x, y) {
		t.Fatalf("a defeated dummy's tile should be walkable")
	}
}

// viewSize --> la caméra reste dans les bornes prévues (largeur/hauteur impaires, tailles raisonnables).
func TestViewSizeStaysInBounds(t *testing.T) {
	w, h := viewSize()
	if w < 11 || w > 31 || w%2 == 0 {
		t.Fatalf("viewSize width = %d, want odd in [11,31]", w)
	}
	if h < 9 || h > 19 || h%2 == 0 {
		t.Fatalf("viewSize height = %d, want odd in [9,19]", h)
	}
}
