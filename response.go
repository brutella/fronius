package fronius

import (
	"fmt"
)

const (
	// 0-6 == Startup
	InverterStatusRunning     byte = 7
	InverterStatusStandby          = 8
	InverterStatusBootLoading      = 9
	InverterStatusError            = 10
)

const (
	StatusOk                 byte = 0   // Request successfully finished, Data are valid
	StatusNotImplemented          = 1   // The request or a part of the request is not implemented yet
	StatusUninitialized           = 2   // Instance of APIRequest created, but not yet configured
	StatusInitialized             = 3   // Request is configured and ready to be sent
	StatusRunning                 = 4   // Request is currently being processed (waiting for response)
	StatusTimeout                 = 5   // Response was not received within desired time
	StatusArgumentError           = 6   // Invalid arguments/combination of arguments or missing arguments
	StatusLNRequestError          = 7   // Something went wrong during sending/receiving of LN-message
	StatusLNRequestTimeout        = 8   // LN-request timed out
	StatusLNParseError            = 9   // Something went wrong during parsing of successfully received LN-message
	StatusConfigIOError           = 10  // Something went wrong while reading settings from local config
	StatusNotSupported            = 11  // The operation/feature or whatever is not supported
	StatusDeviceNotAvailable      = 12  // The device is not available
	StatusUnknownError            = 255 // undefined runtime error
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

type Inverter3PhasesDeviceResponse struct {
	Response
	Body struct {
		Data struct {
			Current1      Property `json:"IAC_L1"`
			Current2      Property `json:"IAC_L2"`
			Current3      Property `json:"IAC_L3"`
			Voltage1      Property `json:"UAC_L1"`
			Voltage2      Property `json:"UAC_L2"`
			Voltage3      Property `json:"UAC_L3"`
			FanFrontLeft  Property `json:"ROTATION_SPEED_FAN_FL"`
			FanFrontRight Property `json:"ROTATION_SPEED_FAN_FR"`
			FanBackLeft   Property `json:"ROTATION_SPEED_FAN_BL"`
			FanBackRight  Property `json:"ROTATION_SPEED_FAN_BR"`
		}
	}
}

type InverterCommonDeviceResponse struct {
	Response
	Body struct {
		Data struct {
			Power     Property `json:"PAC"`
			Current   Property `json:"IAC"`
			Voltage   Property `json:"UAC"`
			Frequence Property `json:"FAC"`
			CurrentDC Property `json:"IDC"`
			VoltageDC Property `json:"UDC"`
		}
	}
}

type InverterSystemResponse struct {
	Response
	Body struct {
		Data struct {
			Power          SystemProperty `json:"PAC"`
			Current        SystemProperty `json:"IAC"`
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

type InverterSystemReponseData struct {
}

type Property struct {
	Value float64
	Unit  string
}

type NamedProperty struct {
	Value float64
	Unit  string
}

func (p Property) String() string {
	return fmt.Sprintf("%v %v", p.Value, p.Unit)
}

type SystemProperty struct {
	Values map[string]float64
	Unit   string
}

func Sum(p SystemProperty) float64 {
	var sum float64
	for _, value := range p.Values {
		sum = sum + value
	}

	return sum
}

func SystemCurrentPower(r *InverterSystemResponse) Property {
	value := Sum(r.Body.Data.Power)
	unit := r.Body.Data.Power.Unit

	return Property{value, unit}
}

func SystemEnergyToday(r *InverterSystemResponse) Property {
	value := Sum(r.Body.Data.EnergyToday)
	unit := r.Body.Data.EnergyToday.Unit

	return Property{value, unit}
}

func SystemEnergyThisYear(r *InverterSystemResponse) Property {
	value := Sum(r.Body.Data.EnergyThisYear)
	unit := r.Body.Data.EnergyThisYear.Unit

	return Property{value, unit}
}

func SystemEnergyTotal(r *InverterSystemResponse) Property {
	value := Sum(r.Body.Data.EnergyTotal)
	unit := r.Body.Data.EnergyTotal.Unit

	return Property{value, unit}
}
