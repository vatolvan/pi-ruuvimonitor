package main

import (
	"fmt"
	"log"
	"time"

	"github.com/influxdata/influxdb1-client/v2"
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
		fmt.Println("id: ", p.ID(), ", temp: ", b.temperature, ", humidity: ", b.humidity, 
			", pressure: ", b.pressure, ", acc: ", [3]float64{b.accelerationX, b.accelerationY, b.accelerationZ},
			", voltage: ", b.voltage, ", txPower: ", b.txPower)
		influxInsertMeasurement(p.ID(), b)
	}
}

func influxInsertMeasurement(macAddress string, ruuvi *RuuviTag) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if err != nil {
		fmt.Println("Error creating InfluxDB Client: ", err.Error())
	}
	defer c.Close()
	
	// Create a new point batch
	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "ruuvi_measurements",
		Precision: "s",
	})

	// Create a point and add to batch
	tags := map[string]string{"mac_address": macAddress}
	fields := map[string]interface{}{
		"temperature":   ruuvi.temperature,
		"humidity": ruuvi.humidity,
		"pressure":   ruuvi.pressure,
		"acceleration_x": ruuvi.accelerationX,
		"acceleration_y": ruuvi.accelerationY,
		"acceleration_z": ruuvi.accelerationZ,
		"voltage": ruuvi.voltage,
		"txPower": ruuvi.txPower,
	}
	pt, err := client.NewPoint("ruuvi_measurement", tags, fields, time.Now())
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
	bp.AddPoint(pt)

	// Write the batch
	c.Write(bp)
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
