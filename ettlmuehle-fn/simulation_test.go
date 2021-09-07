package function

import (
	"fmt"
	"testing"
	"time"
)

func TestSimulation(t *testing.T) {
	for fakeMonth := 0; fakeMonth < 12; fakeMonth++ {
		fmt.Printf("Month: %d", fakeMonth)
		for fakeHour := 0; fakeHour < 23; fakeHour++ {
			fakeTime := time.Date(2021, time.Month(fakeMonth), 1, fakeHour, 0, 0, 0, time.Local)
			value := simulateWaterSensor("", fakeTime)
			fmt.Printf(" %v ", value)
		}
		fmt.Println("")
	}
}
