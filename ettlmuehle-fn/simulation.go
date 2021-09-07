package function

import (
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// SimulateAll triggers simulation of new values
func simulateSensors(connctdClient APIClient) {
	waterLevelSensorOne := os.Getenv("thing-id-water-level-sensor-one")
	waterLevelSensorTwo := os.Getenv("thing-id-water-level-sensor-two")

	currTime := time.Now()

	thingValueOne := simulateWaterSensor(waterLevelSensorOne, currTime)
	thingValueTwo := simulateWaterSensor(waterLevelSensorTwo, currTime)

	connctdClient.UpdateProperty(waterLevelSensorOne, "waterlevel", "value", strconv.FormatInt(thingValueOne, 10))
	connctdClient.UpdateProperty(waterLevelSensorTwo, "waterlevel", "value", strconv.FormatInt(thingValueTwo, 10))
}

func simulateWaterSensor(thingID string, dayTime time.Time) int64 {
	// stretched and shifted sinus wave where max (y=2) is at around x=1 and x=12 and min (y=0) is at around x=7
	// meaning: at around month 7 (july) water level is lowest (dry sommer, wet winter)
	yeartimeMultiplier := math.Sin(float64(dayTime.Month())*0.55+3.14/2) + 1

	// something between 0 (summer) and 1 (winter)
	yeartimeMultiplier = yeartimeMultiplier / 2

	// stretched and shifted sinus wave where max (y=2) is at around x=0 and x=24 and min (y=1) is at around x=12
	// meaning: water level drops until middle of day and then raises again (wet night, dry day)
	daytimeMultiplier := math.Sin(float64(dayTime.Hour())*0.26+3.14/2) + 1

	// something between 0 (middle of day) and 1 (night time)
	daytimeMultiplier = daytimeMultiplier / 2.0

	//minLevel := 0.0
	maxLevel := 1500.0
	perDayWaterLevelVariation := 50.0 // potential variation (night - day)

	yeartimeWaterLevel := yeartimeMultiplier * maxLevel

	simulatedValue := randFloat(yeartimeWaterLevel-(yeartimeWaterLevel/10), yeartimeWaterLevel, 2) // 0 und 0
	simulatedValue += (daytimeMultiplier * perDayWaterLevelVariation) / 10

	return int64(simulatedValue)
}

func randFloat(min float64, max float64, round float64) float64 {
	var mult = math.Pow(10.0, round)
	minInc := int64(min * mult)
	maxInc := int64(max * mult)
	rand.Seed(time.Now().UnixNano())

	if maxInc == minInc {
		maxInc = maxInc + 1
	}

	result := rand.Int63n(maxInc-minInc) + minInc

	return float64(result) / mult
}
