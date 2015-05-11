package solar

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

type InverterResponse struct {
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

type Property struct {
	Value float64
	Unit  string
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
