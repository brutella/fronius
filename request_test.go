package solar

import (
	"testing"
)

func TestGetRequest(t *testing.T) {
	if is, want := GetRealtimeDataRequestURL("127.0.0.1"), "http://127.0.0.1/solar_api/v1/GetInverterRealtimeData.cgi?Scope=System"; is != want {
		t.Fatalf("is=%v want=%v", is, want)
	}
}
