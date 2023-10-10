package fronius

import (
	"fmt"
)

// InverterSystemRealtimeDataRequestURL returns an url to request realtime inverter
// data from host.
func InverterSystemRealtimeDataRequestURL(host string) string {
	base := path(toHTTPHost(host), solarAPIFormat, "GetInverterRealtimeData.cgi")
	arg := args(fmt.Sprintf("Scope=%s", ScopeSystem))

	return append(base, arg)
}

func InverterDeviceRealtimeDataRequestURL(host string, deviceId int, collection string) string {
	base := path(toHTTPHost(host), solarAPIFormat, "GetInverterRealtimeData.cgi")
	arg := args(fmt.Sprintf("Scope=%s", ScopeDevice), ScopeSystem, fmt.Sprintf("DeviceId=%d", deviceId), fmt.Sprintf("DataCollection=%s", collection))

	return append(base, arg)
}
