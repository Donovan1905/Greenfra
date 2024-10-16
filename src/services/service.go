package services

import "fmt"

type ResourceChange struct {
	Address string                 `json:"address"`
	Mode    string                 `json:"mode"`
	Type    string                 `json:"type"`
	Name    string                 `json:"name"`
	Change  map[string]interface{} `json:"change"`
}

type Service interface {
	Analyze(changes []ResourceChange) error
}

const powerConsumptionPerVCPU = 2.10           // Wh per vCPU
const powerConsumptionPerMBofMemory = 0.000384 // Wh per MB of memory
const hoursInMonth = 720

var carbonCosts = map[string]float64{ // gCO2eq per KWh
	"eu-west-3":    20,
	"eu-central-1": 200,
}

func calculateMonthlyPowerConsumption(vCPUs float64, memory int, durationInHours float64) float64 {
	return (vCPUs * powerConsumptionPerVCPU * durationInHours) + (float64(memory) * powerConsumptionPerMBofMemory * durationInHours)
}

func calculateCarbonFootprint(power float64, region string) float64 {
	regionalCarbonCost, ok := carbonCosts[region]
	if !ok {
		fmt.Println("Unknown region. Using default carbon cost.")
	}

	return power * regionalCarbonCost
}
