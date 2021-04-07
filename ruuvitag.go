package main

import (
	// "fmt"
	// "encoding/binary"
	"errors"
	"math"
)

var allowedMacAddresses = []string{"C9:53:A2:99:0E:51"}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

// RuuviTag data container for RuuviTag data
type RuuviTag struct {
	temperature   float64
	humidity      float64
	pressure      float64
	accelerationX float64
	accelerationY float64
	accelerationZ float64
	voltage float64
	txPower int16
}

// NewRuuviTag Parse manufacturer data and output ruuvitag information
func NewRuuviTag(macAddress string, manufacturerData []byte) (*RuuviTag, error) {
	if !contains(allowedMacAddresses, macAddress) {
		return nil, errors.New("Not monitored RuuviTag")
	}
	ruuvitag := new(RuuviTag)
	ruuvitag.temperature = float64(int16(manufacturerData[3])<<8|int16(manufacturerData[4])&0xFF) * 0.005
	ruuvitag.humidity = float64(uint16(manufacturerData[5])<<8|uint16(manufacturerData[6])&0xFF) * 0.0025
	ruuvitag.pressure = (float64(uint16(manufacturerData[7])<<8|uint16(manufacturerData[8])&0xFF) + 50000) * 0.01
	ruuvitag.accelerationX = float64(int16(manufacturerData[9])<<8|int16(manufacturerData[10])&0xFF) * 0.001
	ruuvitag.accelerationY = float64(int16(manufacturerData[11])<<8|int16(manufacturerData[12])&0xFF) * 0.001
	ruuvitag.accelerationZ = float64(int16(manufacturerData[13])<<8|int16(manufacturerData[14])&0xFF) * 0.001

	powerInfo := (uint16(manufacturerData[15]) & 0xFF) << 8 | uint16(manufacturerData[16]) & 0xFF
	if powerInfo >> 5 != 0x7FF {
		ruuvitag.voltage = float64((uint16(powerInfo) >> 5) / 1000) + 1.6;
	}
	if powerInfo & 0x1F != 0x1F {
		ruuvitag.txPower = int16(powerInfo & 0x1F) * 2 - 40;
	}

	ruuvitag.temperature = roundTo(ruuvitag.temperature, 2)
	ruuvitag.humidity = roundTo(ruuvitag.humidity, 2)
	ruuvitag.pressure = roundTo(ruuvitag.pressure, 2)
	ruuvitag.accelerationX = roundTo(ruuvitag.accelerationX, 4)
	ruuvitag.accelerationY = roundTo(ruuvitag.accelerationY, 4)
	ruuvitag.accelerationZ = roundTo(ruuvitag.accelerationZ, 4)
	ruuvitag.voltage = roundTo(ruuvitag.voltage, 4)
	

	return ruuvitag, nil
}

func roundTo(value float64, precision int) float64 {
	return float64(math.Floor(value * math.Pow10(precision)) / math.Pow10(precision))
}