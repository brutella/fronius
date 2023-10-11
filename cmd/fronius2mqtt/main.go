// This program reads energy data from a Fronius Symo Gen24 inverter and stores them into a InfluxDB database.
package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/brutella/fronius"
	"github.com/pascaldekloe/mqtt"
)

var (
	inverterAddr = flag.String("inverter", fronius.SymoHostClassA, "The ip address of the Fronius inverter.")
	mqttAddr     = flag.String("mqtt_addr", "", "URL of the MQTT broker")
	mqttUser     = flag.String("mqtt_user", "", "MQTT user name")
	mqttPassword = flag.String("mqtt_pw", "", "MQTT password")
	tlsFlag      = flag.Bool("mqtt_tls", false, "Use TLS to connect to the broker")
	delay        = flag.Duration("delay", 0*time.Second, "Delay before sending a request")
)

func config() *mqtt.Config {
	var TLS *tls.Config
	if *tlsFlag {
		TLS = new(tls.Config)
	}

	config := mqtt.Config{
		PauseTimeout: 5 * time.Second,
		UserName:     *mqttUser,
		Password:     []byte(*mqttPassword),
	}

	if TLS != nil {
		config.Dialer = mqtt.NewTLSDialer("tcp", *mqttAddr, TLS)
	} else {
		config.Dialer = mqtt.NewDialer("tcp", *mqttAddr)
	}
	return &config
}

func main() {
	flag.Parse()

	if *inverterAddr == "" || *mqttAddr == "" || *mqttUser == "" {
		flag.Usage()
		return
	}

	// gokrazy expects a non-zero status code
	defer os.Exit(1)

	if delay != nil {
		time.Sleep(*delay)
	}

	mqttClient, err := mqtt.VolatileSession("fronius2mqtt", config())
	if err != nil {
		log.Fatal(err)
	}

	config := fronius.Config{
		Host:    *inverterAddr,
		Timeout: 5 * time.Second,
	}
	client := fronius.NewClient(config)
	sys, err := client.GetInverterSystemRealtimeData()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v", sys)

	go func() {
		for {
			message, topic, err := mqttClient.ReadSlices()
			log.Println(message, topic, err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = mqttClient.Publish(ctx.Done(), []byte(fmt.Sprintf("%v", fronius.SystemCurrentPower(sys))), "pv/inverter/current")

	log.Println(err)

	mqttClient.Disconnect(nil)
	mqttClient.Close()
}
