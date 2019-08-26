package main

import (
	// "fmt"
	// "encoding/binary"
	"errors"	
)

var allowedMacAddresses = []string {"C9:53:A2:99:0E:51"}

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
	temperature float32
	humidity float32
	pressure uint16
	acceleration [3]float32	
	data []byte
}

// NewRuuviTag Parse manufacturer data and output ruuvitag information
func NewRuuviTag(macAddress string, manufacturerData []byte) (*RuuviTag, error) {
	if !contains(allowedMacAddresses, macAddress) {
		return nil, errors.New("Not monitored RuuviTag")
	}
	ruuvitag := new(RuuviTag)		
	ruuvitag.temperature = float32(int16(manufacturerData[3]) << 8 | int16(manufacturerData[4]) & 0xFF) * 0.005;
	ruuvitag.humidity = float32(uint16(manufacturerData[5]) << 8 | uint16(manufacturerData[6]) & 0xFF) * 0.0025;
	ruuvitag.pressure = (uint16(manufacturerData[7]) << 8 | uint16(manufacturerData[8]) & 0xFF)
	ruuvitag.acceleration = [3]float32{
		float32(int16(manufacturerData[9]) << 8 | int16(manufacturerData[10]) & 0xFF) * 0.001,
		float32(int16(manufacturerData[11]) << 8 | int16(manufacturerData[12]) & 0xFF) * 0.001,
		float32(int16(manufacturerData[13]) << 8 | int16(manufacturerData[14]) & 0xFF) * 0.001,
	}
	ruuvitag.data = manufacturerData
	return ruuvitag, nil
}
