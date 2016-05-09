package fronius

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	SymoHostClassA          = "http://169.254.0.180"
	solarAPIFormat          = "solar_api/v1"
	getInverterRealtimeData = "GetInverterRealtimeData.cgi"
	sysScope                = "Scope=System"
)

func SystemRealtimeDataRequestURL(host string) string {
	base := path(toHTTPHost(host), solarAPIFormat, getInverterRealtimeData)
	arg := args(sysScope)

	return append(base, arg)
}

func GetSystemRealtimeData(base string) (*http.Response, error) {
	url := SystemRealtimeDataRequestURL(base)

	return http.Get(url)
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
