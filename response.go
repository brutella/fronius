package fronius

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Response is the common response header (CRH)
type Response struct {
	Head struct {
		RequestArguments struct {
			Scope          string
			DeviceId       string
			DataCollection string
		}
		Status struct {
			Code        int16
			Reason      string
			UserMessage string
		}
		Timestamp string
	}
}

// InverterDeviceResponse is the 
type InverterDeviceResponse struct {
	Response
	Body struct {
		Data struct {
			Power          Property `json:"PAC"`
			EnergyToday    Property `json:"DAY_ENERGY"`
			EnergyThisYear Property `json:"YEAR_ENERGY"`
			EnergyTotal    Property `json:"TOTAL_ENERGY"`
			Status         struct {
				Status                 int16 `json:"StatusCode"`
				Error                  int16 `json:"ErrorCode"`
				LEDColor               int16 // ?
				LEDState               int16 // ?
				Reset                  bool  // ?
				MgmtTimerRemainingTime int16 // ?
			} `json:"DeviceStatus"`
		}
	}
}

type InverterSystemResponse struct {
	Response
	Body struct {
		Data struct {
			Power          SystemProperty `json:"PAC"`
			EnergyToday    SystemProperty `json:"DAY_ENERGY"`
			EnergyThisYear SystemProperty `json:"YEAR_ENERGY"`
			EnergyTotal    SystemProperty `json:"TOTAL_ENERGY"`
			Status         struct {
				Status                 int16 `json:"StatusCode"`
				Error                  int16 `json:"ErrorCode"`
				LEDColor               int16 // ?
				LEDState               int16 // ?
				Reset                  bool  // ?
				MgmtTimerRemainingTime int16 // ?
			} `json:"DeviceStatus"`
		}
	}
}

func NewInverterSystemResponse(res *http.Response) (inv InverterSystemResponse, err error) {
	b, err := ioutil.ReadAll(res.Body)

	if err == nil {
		err = json.Unmarshal(b, &inv)
	}

	return inv, err
}

type Property struct {
	Value int64
	Unit  string
}

func (p Property) String() string {
	return fmt.Sprintf("%v %v", p.Value, p.Unit)
}

type SystemProperty struct {
	Values map[string]int64
	Unit   string
}

func Sum(p SystemProperty) int64 {
	var sum int64
	for _, value := range p.Values {
		sum = sum + value
	}

	return sum
}

const (
	// 0-6 == Startup
	StatusRunning     int16 = 7
	StatusStandby     int16 = 8
	StatusBootLoading int16 = 9
	StatusError       int16 = 10
)

// FIXME(brutella) Define more error codes
const (
	ErrorNo int16 = 0
)

func SystemCurrentPower(r InverterSystemResponse) Property {
	value := Sum(r.Body.Data.Power)
	unit := r.Body.Data.Power.Unit

	return Property{value, unit}
}

func SystemEnergyToday(r InverterSystemResponse) Property {
	value := Sum(r.Body.Data.EnergyToday)
	unit := r.Body.Data.EnergyToday.Unit

	return Property{value, unit}
}

func SystemEnergyThisYear(r InverterSystemResponse) Property {
	value := Sum(r.Body.Data.EnergyThisYear)
	unit := r.Body.Data.EnergyThisYear.Unit

	return Property{value, unit}
}

func SystemEnergyTotal(r InverterSystemResponse) Property {
	value := Sum(r.Body.Data.EnergyTotal)
	unit := r.Body.Data.EnergyTotal.Unit

	return Property{value, unit}
}
