// Example to demonstrate how to get realitme inverter data
package main

import (
	"fmt"
	"github.com/brutella/fronius"
	"log"
	"net"
	"net/http"
	"time"
)

// from http://stackoverflow.com/a/16930649/424814
func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, cTimeout)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(rwTimeout))
		return conn, nil
	}
}

func NewTimeoutClient(connectTimeout time.Duration, readWriteTimeout time.Duration) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial: TimeoutDialer(connectTimeout, readWriteTimeout),
		},
	}
}

func main() {
	client := NewTimeoutClient(5*time.Second, 5*time.Second)
	resp, err := client.Get(fronius.SystemRealtimeDataRequestURL(fronius.SymoHostClassA))

	if err != nil {
		log.Fatal(err)
	}

	inv, err := fronius.NewInverterSystemResponse(resp)

	fmt.Printf("current power: %v\n", fronius.SystemCurrentPower(inv))
	fmt.Printf("today: %v\n", fronius.SystemEnergyToday(inv))
	fmt.Printf("this year: %v\n", fronius.SystemEnergyThisYear(inv))
	fmt.Printf("total: %v\n", fronius.SystemEnergyTotal(inv))
}
