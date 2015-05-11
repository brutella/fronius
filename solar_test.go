package solar

import (
	"encoding/json"
	"testing"
)

func TestRealtimeData(t *testing.T) {
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
	rsp := InverterResponse{}
	if err := json.Unmarshal(b, &rsp); err != nil {
		t.Fatal(err)
	}

	if x := rsp.Head.RequestArguments.Scope; x != "Device" {
		t.Error(x)
	}

	if x := rsp.Head.RequestArguments.DeviceId; x != "0" {
		t.Error(x)
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
