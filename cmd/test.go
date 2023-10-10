// This program reads energy data from a Fronius Symo Gen24 inverter and stores them into a InfluxDB database.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/brutella/fronius"
)

func main() {
	var (
		inverterAddr = flag.String("inverter", fronius.SymoHostClassA, "The ip address of the Fronius inverter.")
	)
	flag.Parse()

	config := fronius.Config{
		Host:    *inverterAddr,
		Timeout: 5 * time.Second,
	}
	config.Host = "192.168.1.65"
	client := fronius.NewClient(config)
	inverter, err := client.Get3PInverterData(1)
	if err != nil {
		log.Fatal(err)
	}

	json, _ := json.MarshalIndent(inverter, "", "\t")
	fmt.Print(string(json))

	meter, err := client.GetMeterSystemRealtimeData()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v", meter)

	// fmt.Println("------- Inverter --------")
	//     fmt.Printf("producing: %v\n", fronius.SystemCurrentPower(sys))
	//     fmt.Printf("today: %v\n", fronius.SystemEnergyToday(sys))
	//     fmt.Printf("this year: %v\n", fronius.SystemEnergyThisYear(sys))
	//     fmt.Printf("total: %v\n", fronius.SystemEnergyTotal(sys))

	fmt.Println("------- SMART METER --------")
	fmt.Printf("imported: %v\n", fronius.SmartMeterGridEnergyImportSum(meter))
	fmt.Printf("exported: %v\n", fronius.SmartMeterGridEnergyExportSum(meter))
}
