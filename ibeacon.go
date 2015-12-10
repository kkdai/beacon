package beacon

import (
	"fmt"
	"log"

	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
	"github.com/paypal/gatt/examples/service"
)

type IBeacon struct {
	serviceList []*gatt.Service
	beaconDev   *gatt.Device

	//Device Information
	DevUUID         string
	DevName         string
	DevMajorVersion uint16
	DevMinorVersion uint16
	PowerLevel      int8
}

func NewIBeacon(devUUID string, name string, powerLevel int8) *IBeacon {
	ib := new(IBeacon)
	ib.DevUUID = devUUID
	ib.DevName = name
	ib.DevMajorVersion = 1 //default
	ib.DevMinorVersion = 1 //default
	ib.PowerLevel = powerLevel
	return ib
}

func (ib *IBeacon) SetiBeaconVersion(major, minor uint16) {
	ib.DevMajorVersion = major
	ib.DevMinorVersion = minor
}

func (ib *IBeacon) AddBatteryService() {
	sev := service.NewBatteryService()
	ib.serviceList = append(ib.serviceList, sev)
}

func (ib *IBeacon) AddCountService() {
	sev := service.NewCountService()
	ib.serviceList = append(ib.serviceList, sev)
}

func (ib *IBeacon) Advertise() error {

	d, err := gatt.NewDevice(option.DefaultServerOptions...)
	if err != nil {
		log.Fatalf("Failed to open device, err: %s", err)
		return err
	}

	// Register optional handlers.
	d.Handle(
		gatt.CentralConnected(func(c gatt.Central) { fmt.Println("Connect: ", c.ID()) }),
		gatt.CentralDisconnected(func(c gatt.Central) { fmt.Println("Disconnect: ", c.ID()) }),
	)

	// A mandatory handler for monitoring device state.
	onStateChanged := func(d gatt.Device, s gatt.State) {
		fmt.Printf("State: %s\n", s)
		switch s {
		case gatt.StatePoweredOn:
			// Setup GAP and GATT services for Linux implementation.
			// OS X doesn't export the access of these services.
			d.AddService(service.NewGapService(ib.DevName)) // no effect on OS X
			d.AddService(service.NewGattService())          // no effect on OS X

			//Append services
			serviceUUIDs := []gatt.UUID{}
			for _, v := range ib.serviceList {
				d.AddService(v)
				serviceUUIDs = append(serviceUUIDs, v.UUID())
			}

			// Advertise device name and service's UUIDs.
			d.AdvertiseNameAndServices(ib.DevName, serviceUUIDs)

			// Advertise as an OpenBeacon iBeacon
			d.AdvertiseIBeacon(gatt.MustParseUUID(ib.DevUUID), ib.DevMajorVersion, ib.DevMinorVersion, ib.PowerLevel)

		default:
		}
	}

	d.Init(onStateChanged)
	select {}
}
