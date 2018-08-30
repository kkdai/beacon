package beacon

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib"
	"github.com/go-ble/ble/examples/lib/dev"
)

type IBeacon struct {
	serviceList []*ble.Service
	BeaconDev   ble.Device

	//Device Information
	DevUUID         string
	DevName         string
	DevMajorVersion uint16
	DevMinorVersion uint16
	PowerLevel      int8
}

func NewIBeacon(devUUID string, name string, powerLevel int8) *IBeacon {
	d, err := dev.NewDevice("default")
	if err != nil {
		log.Fatalf("can't new device : %s", err)
	}

	return &IBeacon{
		BeaconDev:       d,
		DevUUID:         devUUID,
		DevName:         name,
		DevMajorVersion: 1,
		DevMinorVersion: 1,
		PowerLevel:      powerLevel,
	}
	return &IBeacon{}
}

func (ib *IBeacon) SetiBeaconVersion(major, minor uint16) {
	ib.DevMajorVersion = major
	ib.DevMinorVersion = minor
}

func (ib *IBeacon) AddBatteryService() {
	sev := ble.NewService(ble.BatteryUUID)
	ib.serviceList = append(ib.serviceList, sev)
}

func (ib *IBeacon) AddCountService() {
	testSvc := ble.NewService(lib.TestSvcUUID)
	testSvc.AddCharacteristic(lib.NewCountChar())
	testSvc.AddCharacteristic(lib.NewEchoChar())
	ib.serviceList = append(ib.serviceList, testSvc)
}

func (ib *IBeacon) Advertise() error {
	ble.SetDefaultDevice(ib.BeaconDev)

	for _, service := range ib.serviceList {
		if err := ble.AddService(service); err != nil {
			log.Fatalf("can't add service: %s", err)
		}
	}

	// Advertise for specified durantion, or until interrupted by user.
	fmt.Printf("Advertising for %s...\n", 30*time.Second)
	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), 30*time.Second))
	return ble.AdvertiseNameAndServices(ctx, "Gopher")

	// // Register optional handlers.
	// d.Handle(
	// 	gatt.CentralConnected(func(c gatt.Central) { fmt.Println("Connect: ", c.ID()) }),
	// 	gatt.CentralDisconnected(func(c gatt.Central) { fmt.Println("Disconnect: ", c.ID()) }),
	// )

	// // A mandatory handler for monitoring device state.
	// onStateChanged := func(d gatt.Device, s gatt.State) {
	// 	fmt.Printf("State: %s\n", s)
	// 	switch s {
	// 	case gatt.StatePoweredOn:
	// 		// Setup GAP and GATT services for Linux implementation.
	// 		// OS X doesn't export the access of these services.
	// 		d.AddService(service.NewGapService(ib.DevName)) // no effect on OS X
	// 		d.AddService(service.NewGattService())          // no effect on OS X

	// 		//Append services
	// 		serviceUUIDs := []gatt.UUID{}
	// 		for _, v := range ib.serviceList {
	// 			d.AddService(v)
	// 			serviceUUIDs = append(serviceUUIDs, v.UUID())
	// 		}

	// 		// Advertise device name and service's UUIDs.
	// 		d.AdvertiseNameAndServices(ib.DevName, serviceUUIDs)

	// 		// Advertise as an OpenBeacon iBeacon
	// 		d.AdvertiseIBeacon(gatt.MustParseUUID(ib.DevUUID), ib.DevMajorVersion, ib.DevMinorVersion, ib.PowerLevel)

	// 	default:
	// 	}
	// }

	// d.Init(onStateChanged)
	// select {}
}
