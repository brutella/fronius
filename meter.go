package fronius

type MeterSystemResponse struct {
	Response
	Body struct {
		Data map[string]SmartMeterReponseData
	}
}

type SmartMeterReponseData struct {
	Details struct {
		Manufacturer string
		Model        string
		Serial       string
	}
	Current_AC_Phase_1                    float64
	Current_AC_Phase_2                    float64
	Current_AC_Phase_3                    float64
	Current_AC_Sum                        float64
	Enable                                byte
	EnergyReactive_VArAC_Phase_1_Consume  float64
	EnergyReactive_VArAC_Phase_1_Produced float64
	EnergyReactive_VArAC_Sum_Consumed     float64
	EnergyReactive_VArAC_Sum_Produced     float64
	EnergyReal_WAC_Minus_Absolute         float64
	EnergyReal_WAC_Phase_1_Consumed       float64
	EnergyReal_WAC_Phase_1_Produced       float64
	EnergyReal_WAC_Phase_2_Consumed       float64
	EnergyReal_WAC_Phase_2_Produced       float64
	EnergyReal_WAC_Phase_3_Consumed       float64
	EnergyReal_WAC_Phase_3_Produced       float64
	EnergyReal_WAC_Plus_Absolute          float64
	EnergyReal_WAC_Sum_Consumed           float64
	EnergyReal_WAC_Sum_Produced           float64
	Frequency_Phase_Average               float64
	Meter_Location_Current                float32 // 0-511 (uint8 doesn't work because the api reports float values)
	PowerApparent_S_Phase_1               float64
	PowerApparent_S_Phase_2               float64
	PowerApparent_S_Phase_3               float64
	PowerApparent_S_Sum                   float64
	PowerFactor_Phase_1                   float64
	PowerFactor_Phase_2                   float64
	PowerFactor_Phase_3                   float64
	PowerFactor_Sum                       float64
	PowerReactive_Q_Phase_1               float64
	PowerReactive_Q_Phase_2               float64
	PowerReactive_Q_Phase_3               float64
	PowerReactive_Q_Sum                   float64
	PowerReal_P_Phase_1                   float64
	PowerReal_P_Phase_2                   float64
	PowerReal_P_Phase_3                   float64
	PowerReal_P_Sum                       float64
	TimeStamp                             float64
	Visible                               byte
	Voltage_AC_PhaseToPhase_12            float64
	Voltage_AC_PhaseToPhase_23            float64
	Voltage_AC_PhaseToPhase_31            float64
	Voltage_AC_Phase_1                    float64
	Voltage_AC_Phase_2                    float64
	Voltage_AC_Phase_3                    float64
	Voltage_AC_Phase_Average              float64
}

// SmartMeterGridPower returns the current power from the grid.
// A positive value means that power is consumed from the grid.
// A negative value means that power is feed in to the grid.
func SmartMeterGridPowerSum(res *MeterSystemResponse) Property {
	var sum float64
	for _, data := range res.Body.Data {
		if data.Meter_Location_Current == 0 {
			sum += data.PowerReal_P_Sum
		} else {
			sum -= data.PowerReal_P_Sum
		}
	}

	return Property{sum, "W"}
}

// SmartMeterGridPower returns the current power from the grid.
// A positive value means that power is consumed from the grid.
// A negative value means that power is feed in to the grid.
func SmartMeterGridPower(data SmartMeterReponseData) float64 {
	if data.Meter_Location_Current == 0 {
		return data.PowerReal_P_Sum
	}

	return -data.PowerReal_P_Sum
}

func SmartMeterGridEnergyImport(data SmartMeterReponseData) float64 {
	if data.Meter_Location_Current == 0 {
		return data.EnergyReal_WAC_Plus_Absolute
	}

	return -data.EnergyReal_WAC_Minus_Absolute
}

func SmartMeterGridEnergyImportSum(res *MeterSystemResponse) Property {
	var sum float64
	for _, data := range res.Body.Data {
		if data.Meter_Location_Current == 0 {
			sum += data.EnergyReal_WAC_Plus_Absolute
		} else {
			sum -= data.EnergyReal_WAC_Minus_Absolute
		}
	}

	return Property{sum, "Wh"}
}

func SmartMeterGridEnergyExport(data SmartMeterReponseData) float64 {
	if data.Meter_Location_Current == 0 {
		return data.EnergyReal_WAC_Minus_Absolute
	}

	return -data.EnergyReal_WAC_Plus_Absolute
}

func SmartMeterGridEnergyExportSum(res *MeterSystemResponse) Property {
	var sum float64
	for _, data := range res.Body.Data {
		if data.Meter_Location_Current == 0 {
			sum += data.EnergyReal_WAC_Minus_Absolute
		} else {
			sum -= data.EnergyReal_WAC_Plus_Absolute
		}
	}

	return Property{sum, "Wh"}
}
