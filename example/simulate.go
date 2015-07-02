package main

import (
	"fmt"
	"github.com/brutella/fronius"
	"log"
)

func main() {
	s := fronius.NewSymoSimulator()

	defer s.Stop()

	resp, err := fronius.GetSystemRealtimeData(s.URL())

	if err != nil {
		log.Fatal(err)
	}

	inv, err := fronius.NewInverterSystemResponse(resp)

	fmt.Printf("current power: %v\n", fronius.SystemCurrentPower(inv))
	fmt.Printf("today: %v\n", fronius.SystemEnergyToday(inv))
	fmt.Printf("this year: %v\n", fronius.SystemEnergyThisYear(inv))
	fmt.Printf("total: %v\n", fronius.SystemEnergyTotal(inv))
}
