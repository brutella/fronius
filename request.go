package fronius

import (
	"fmt"
	"strings"
)

const (
	SymoHostClassA          = "169.254.0.180:80"
	solarAPIFormat          = "solar_api/v1"
	getInverterRealtimeData = "GetInverterRealtimeData.cgi"
	sysScope                = "Scope=System"
)

// SystemRealtimeDataRequestURL returns an url to request realtime inverter 
// data from host.
func SystemRealtimeDataRequestURL(host string) string {
	base := path(toHTTPHost(host), solarAPIFormat, getInverterRealtimeData)
	arg := args(sysScope)

	return append(base, arg)
}

func toHTTPHost(ip string) string {
	return fmt.Sprintf("http://%s", ip)
}

func args(a ...string) string {
	return strings.Join(a, "&")
}

func append(a ...string) string {
	return strings.Join(a, "?")
}

func path(a ...string) string {
	return strings.Join(a, "/")
}
