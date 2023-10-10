package fronius

import (
	"fmt"
	"strings"
)

const (
	SymoHostClassA       = "169.254.0.180:80"
	solarAPIFormat       = "solar_api/v1"
	ScopeSystem          = "System"
	ScopeDevice          = "Device"
	CollectionCumulation = "CumulationInverterData"
	CollectionCommon     = "CommonInverterData"
	Collection3Phases    = "3PInverterData"
	CollectionMinMax     = "MinMaxInverterData"
)

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
