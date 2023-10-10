package fronius

import (
	"encoding/json"
	"testing"
)

var systemRealtimeData string = `{
	"Head" : {
		"RequestArguments" : {
			"DataCollection" : "",
			"Scope" : "System"
		},
		"Status" : {
			"Code" : 0,
			"Reason" : "",
			"UserMessage" : ""
		},
		"Timestamp" : "2015-05-23T10:42:29+02:00"
	},
	"Body" : {
		"Data" : {
			"PAC" : {
				"Unit" : "W",
				"Values" : {
					"1" : 766
				}
			},
			"DAY_ENERGY" : {
				"Unit" : "Wh",
				"Values" : {
					"1" : 1622
				}
			},
			"YEAR_ENERGY" : {
				"Unit" : "Wh",
				"Values" : {
					"1" : 46146
				}
			},
			"TOTAL_ENERGY" : {
				"Unit" : "Wh",
				"Values" : {
					"1" : 46146
				}
			}
		}
	}
}`

func TestSystemRealtimeData(t *testing.T) {
	b := []byte(systemRealtimeData)
	rsp := InverterSystemResponse{}
	if err := json.Unmarshal(b, &rsp); err != nil {
		t.Fatal(err)
	}

	if is, want := rsp.Head.RequestArguments.Scope, "System"; is != want {
		t.Fatalf("is=%v want=%v", is, want)
	}

	if x := rsp.Head.Timestamp; x != "2015-05-23T10:42:29+02:00" {
		t.Error(x)
	}

	if is, want := rsp.Body.Data.Power.Values["1"], float64(766); is != want {
		t.Fatalf("is=%v want=%v", is, want)
	}

	if x := rsp.Body.Data.Power.Unit; x != "W" {
		t.Error(x)
	}

	if is, want := rsp.Body.Data.EnergyToday.Values["1"], float64(1622); is != want {
		t.Fatalf("is=%v want=%v", is, want)
	}

	if x := rsp.Body.Data.EnergyToday.Unit; x != "Wh" {
		t.Error(x)
	}
}

func TestSystemRealtimeDataGetter(t *testing.T) {
	b := []byte(systemRealtimeData)
	rsp := &InverterSystemResponse{}
	if err := json.Unmarshal(b, rsp); err != nil {
		t.Fatal(err)
	}

	p := SystemCurrentPower(rsp)
	if is, want := p.Value, float64(766); is != want {
		t.Fatalf("is=%v want=%v", is, want)
	}

	if is, want := p.Unit, "W"; is != want {
		t.Fatalf("is=%v want=%v", is, want)
	}
}

func TestDeviceRealtimeData(t *testing.T) {
	str := `{
    "Head": {
        "RequestArguments": {
            "Scope": "Device",
            "DeviceId": "0",
            "DataCollection": "CommonInverterData"
        },
        "Status": {
            "Code": 0,
            "Reason": "",
            "UserMessage": ""
        },
        "Timestamp": "2011-10-20T10:09:14+02:00"
    },
    "Body" : {
        "Data": {
            "DAY_ENERGY": {
                "Value": 8000,
                "Unit": "Wh"
            },
            "PAC": {
                "Value": 3373,
                "Unit": "W"
            },
            "TOTAL_ENERGY": {
                "Value": 45000,
                "Unit": "Wh"
            },
            "YEAR_ENERGY": {
                "Value": 44000,
                "Unit": "Wh"
            },
            "DeviceStatus": {
                "DeviceState": 7,
                "MgmtTimerRemainingTime": -1,
                "ErrorCode": 0,
                "LEDColor": 2,
                "LEDState": 0,
                "StateToReset": false
            }
        }
    }
}`
	b := []byte(str)
	rsp := InverterDeviceResponse{}
	if err := json.Unmarshal(b, &rsp); err != nil {
		t.Fatal(err)
	}

	if is, want := rsp.Head.RequestArguments.Scope, "Device"; is != want {
		t.Fatalf("is=%v want=%v", is, want)
	}

	if is, want := rsp.Head.RequestArguments.DeviceId, "0"; is != want {
		t.Fatalf("is=%v want=%v", is, want)
	}

	if x := rsp.Head.RequestArguments.DataCollection; x != "CommonInverterData" {
		t.Error(x)
	}

	if x := rsp.Head.Timestamp; x != "2011-10-20T10:09:14+02:00" {
		t.Error(x)
	}

	if x := rsp.Body.Data.Power.Value; x != 3373 {
		t.Error(x)
	}

	if x := rsp.Body.Data.Power.Unit; x != "W" {
		t.Error(x)
	}

	if x := rsp.Body.Data.EnergyToday.Value; x != 8000 {
		t.Error(x)
	}

	if x := rsp.Body.Data.EnergyToday.Unit; x != "Wh" {
		t.Error(x)
	}
}
