// Example to demonstrate how to use the inverter simulator
package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/brutella/fronius"
)

func main() {
	s := fronius.NewSymoSimulator()

	defer s.Stop()

	url, _ := url.Parse(s.URL())
	resp, err := http.Get(fronius.InverterSystemRealtimeDataRequestURL(url.Host))

	if err != nil {
		log.Fatal(err)
	}

	inv, err := fronius.NewInverterSystemResponse(resp)

	fmt.Printf("current power: %v\n", fronius.SystemCurrentPower(inv))
	fmt.Printf("today: %v\n", fronius.SystemEnergyToday(inv))
	fmt.Printf("this year: %v\n", fronius.SystemEnergyThisYear(inv))
	fmt.Printf("total: %v\n", fronius.SystemEnergyTotal(inv))
}
