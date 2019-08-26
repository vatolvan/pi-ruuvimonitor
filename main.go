package main

import (
	"fmt"
	"log"

	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
)

func onStateChanged(device gatt.Device, s gatt.State) {
	switch s {
	case gatt.StatePoweredOn:
		fmt.Println("Scanning for beacon broadcasts...")
		device.Scan([]gatt.UUID{}, true)
		return
	default:
		device.StopScanning()
	}
}

func onPeripheralDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	b, err := NewRuuviTag(p.ID(), a.ManufacturerData)
	if err == nil {
		fmt.Println("id: ", p.ID(), ", temp: ", b.temperature, ", humidity: ", b.humidity, ", pressure: ", b.pressure, ", acc: ", b.acceleration);
		// fmt.Println("Temperature: ", b.temperature, ", humidity: ", b.humidity, ", pressure: ", b.pressure, ", acceleration: ", b.acceleration)
	}
}

func main() {
	device, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		log.Fatalf("Failed to open device, err: %s\n", err)
		return
	}
	device.Handle(gatt.PeripheralDiscovered(onPeripheralDiscovered))
	device.Init(onStateChanged)
	select {}
}