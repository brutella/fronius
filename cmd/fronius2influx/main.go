// This program reads energy data from a Fronius Symo Gen24 inverter and stores them into a InfluxDB database.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/brutella/fronius"
	influxdb "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

func main() {
	var (
		inverterAddr = flag.String("inverter", fronius.SymoHostClassA, "The ip address of the Fronius inverter.")
		dbAddr       = flag.String("influxdb", "http://localhost:8086", "The ip address of InfluxDB.")
		bucket       = flag.String("bucket", "pv", "Name of the InfluxDB bucket.")
		org          = flag.String("org", "", "Name of the InfluxDB organisation.")
		token        = flag.String("token", "", "InfluxDB API token")
		exitOnErr    = flag.Bool("exit-on-err", false, "Flag if app should exit on error.")
	)
	flag.Parse()

	config := fronius.Config{
		Host:    *inverterAddr,
		Timeout: 5 * time.Second,
	}
	// config.Host = "192.168.1.65"
	client := fronius.NewClient(config)
	sys, err := client.GetInverterSystemRealtimeData()
	if err != nil {
		log.Fatal(err)
	}

	phases, err := client.Get3PInverterData(1)
	if err != nil {
		log.Fatal(err)
	}

	meter, err := client.GetMeterSystemRealtimeData()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("------- Inverter --------")
	fmt.Printf("producing: %v\n", fronius.SystemCurrentPower(sys))
	fmt.Printf("today: %v\n", fronius.SystemEnergyToday(sys))
	fmt.Printf("this year: %v\n", fronius.SystemEnergyThisYear(sys))
	fmt.Printf("total: %v\n", fronius.SystemEnergyTotal(sys))

	fmt.Println("------- Phases --------")
	fmt.Printf("1: %v\t%v\n", phases.Body.Data.Current1, phases.Body.Data.Voltage1)
	fmt.Printf("2: %v\t%v\n", phases.Body.Data.Current2, phases.Body.Data.Voltage2)
	fmt.Printf("3: %v\t%v\n", phases.Body.Data.Current3, phases.Body.Data.Voltage3)

	fmt.Println("------- SMART METER --------")
	fmt.Printf("imported: %v\n", fronius.SmartMeterGridEnergyImportSum(meter))
	fmt.Printf("exported: %v\n", fronius.SmartMeterGridEnergyExportSum(meter))

	{
		client := influxdb.NewClient(*dbAddr, *token)
		writeAPI := client.WriteAPIBlocking(*org, *bucket)
		pw := fronius.SmartMeterGridPowerSum(meter).Value
		var in, out float64
		if pw < 0 {
			out = pw * -1
		} else {
			in = pw
		}

		points := []*write.Point{}

		inverter := influxdb.NewPoint(
			"inverter",          // measurement
			map[string]string{}, // tags
			map[string]interface{}{ // fields
				"production":       fronius.SystemCurrentPower(sys).Value,
				"production-total": fronius.SystemEnergyTotal(sys).Value,
			},
			time.Now())
		points = append(points, inverter)

		current := influxdb.NewPoint(
			"current",           // measurement
			map[string]string{}, // tags
			map[string]interface{}{ // fields
				"phase1": phases.Body.Data.Current1.Value,
				"phase2": phases.Body.Data.Current2.Value,
				"phase3": phases.Body.Data.Current3.Value,
			},
			time.Now())

		points = append(points, current)

		voltage := influxdb.NewPoint(
			"voltage",           // measurement
			map[string]string{}, // tags
			map[string]interface{}{ // fields
				"phase1": phases.Body.Data.Voltage1.Value,
				"phase2": phases.Body.Data.Voltage2.Value,
				"phase3": phases.Body.Data.Voltage3.Value,
			},
			time.Now())

		points = append(points, voltage)

		sm := influxdb.NewPoint(
			"meter",             // measurement
			map[string]string{}, // tags
			map[string]interface{}{ // fields
				"import":   in,
				"export":   out,
				"usage":    in + fronius.SystemCurrentPower(sys).Value - out,
				"imported": fronius.SmartMeterGridEnergyImportSum(meter).Value,
				"exported": fronius.SmartMeterGridEnergyExportSum(meter).Value,
			},
			time.Now())
		points = append(points, sm)

		for _, p := range points {
			err = writeAPI.WritePoint(context.Background(), p)
			if err != nil {
				log.Println("error:", err)
				if *exitOnErr {
					return
				}
			}
		}
	}
}
