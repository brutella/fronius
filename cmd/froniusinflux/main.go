// Example to demonstrate how to get realitme inverter data
package main

import (
	"flag"
	"fmt"
	"github.com/brutella/fronius"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
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
	var (
		host     = flag.String("host", fronius.SymoHostClassA, "Host")
		simulate = flag.Bool("simulate", false, "Simulate Fronius symo")
	)

	client := NewTimeoutClient(5*time.Second, 5*time.Second)

	var resp *http.Response
	var err error
	if *simulate {
		s := fronius.NewSymoSimulator()
		defer s.Stop()
		url, _ := url.Parse(s.URL())
		resp, err = client.Get(fronius.SystemRealtimeDataRequestURL(url.Host))
	} else {
		resp, err = client.Get(fronius.SystemRealtimeDataRequestURL(*host))
	}

	if err != nil {
		log.Fatal(err)
	}
	inv, err := fronius.NewInverterSystemResponse(resp)

	pairs := []string{
		fmt.Sprintf("power=%d", fronius.SystemCurrentPower(inv).Value),
		fmt.Sprintf("energy_today=%d", fronius.SystemEnergyToday(inv).Value),
		fmt.Sprintf("energy_year=%d", fronius.SystemEnergyThisYear(inv).Value),
		fmt.Sprintf("energy_total=%d", fronius.SystemEnergyTotal(inv).Value),
	}
	fmt.Printf("fronius,host=%s %s", *host, strings.Join(pairs, ","))
}
