package main

import (
	"fmt"
	"time"
)

var hourDuration = 12*time.Second + 500*time.Millisecond

// Système d'heures : runClock --> fait défiler les heures , prévient sur midnight quand minuit sonne.
func runClock(midnight chan<- struct{}) {
	for heure := 0; heure < 24; heure++ {
		fmt.Printf("\n[Horloge] %02d:00\n> ", heure)
		time.Sleep(hourDuration)
	}
	fmt.Println("\n[Horloge] La cloche du temple sonne... minuit.")
	midnight <- struct{}{}
}
