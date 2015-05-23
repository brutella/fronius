package solar

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	solarAPIFormat          = "solar_api/v1"
	getInverterRealtimeData = "GetInverterRealtimeData.cgi"
	sysScope                = "Scope=System"
	deviceScopeFormat       = "Scope=Device&DeviceId=%s"
)

func GetRealtimeDataRequestURL(ip string) string {
	base := toHTTPHost(ip)
	base = path(base, solarAPIFormat, getInverterRealtimeData)
	arg := args(sysScope)

	return append(base, arg)
}

func GetRealtimeData(ip string) (*http.Response, error) {
	url := GetRealtimeDataRequestURL(ip)

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
