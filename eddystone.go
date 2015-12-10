package beacon

import (
	"fmt"
	"log"

	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
	"github.com/suapapa/go_eddystone"
)

const (
	advTypeAllUUID16     = 0x03 // Complete List of 16-bit Service Class UUIDs
	advTypeServiceData16 = 0x16 // Service Data - 16-bit UUID
)

const (
	advFlagGeneralDiscoverable = 0x02
	advFlagLEOnly              = 0x04
)

type EddyStoneBeacon struct {
	//Output frame data for eddystone beacon
	OutFrame *eddystone.Frame

	//Tx power is the received power at 0 meters, in dBm, and the value ranges from -100 dBm to +20 dBm to a resolution of 1 dBm.
	//
	//Note to developers:
	// The best way to determine the precise value to put into this field is to measure the actual output of
	// your beacon from 1 meter away and then add 41dBm to that. 41dBm is the signal loss that occurs over 1 meter.
	PowerLevel int
}

//Create a Eddystone beacon with URL frame (physical web)
// url:  It need include with "http://" or "https://"
// PowerLevel:  The power level adjustment, refer PowerLevel description
// Ex: NewEddystoneURLBeacon("http://google.com", -20)
func NewEddystoneURLBeacon(url string, powerLevel int) *EddyStoneBeacon {
	eb := new(EddyStoneBeacon)
	eb.PowerLevel = powerLevel
	eb.addURL(url)
	return eb
}

//Create a Eddystone beacon with UID frame (10 bytes for namesapce, 6 bytes for instance)
// Namespace:  10 bytes for UID namespace
// Instance:   6  bytes for instance
// PowerLevel:  The power level adjustment, refer PowerLevel description
// Ex: NewEddystoneUIDBeacon("00010203040506070809", "aabbccddeeff", -20)
func NewEddystoneUIDBeacon(namespace, instance string, powerLevel int) *EddyStoneBeacon {
	eb := new(EddyStoneBeacon)
	eb.PowerLevel = powerLevel
	eb.addUID(namespace, instance)
	return eb
}

//Create a Eddystone beacon with TLM frame (for tansfer telemetry)
// BATT: The battery voltage
// TEMP: Beacon temperature is the temperature in degrees Celsius
// AdvCnt: Count of advertisement frames of all types emitted by the beacon since power-up or reboot
// SecCnt: Counter that represents time since beacon power-up or reboot.
// PowerLevel:  The power level adjustment, refer PowerLevel description
// Ex: NewEddystoneUIDBeacon(500, 22.0, 100, 1000, -20)
func NewEddystoneTLMBeacon(batt uint16, temp float32, advCnt, secCnt uint32, powerLevel int) *EddyStoneBeacon {
	eb := new(EddyStoneBeacon)
	eb.PowerLevel = powerLevel
	eb.addTLM(batt, temp, advCnt, secCnt)
	return eb
}

func (eb *EddyStoneBeacon) addTLM(batt uint16, temp float32, advCnt, secCnt uint32) {
	f, err := eddystone.MakeTLMFrame(batt, temp, advCnt, secCnt)
	if err != nil {
		panic(err)
	}
	eb.OutFrame = &f
}

func (eb *EddyStoneBeacon) addUID(ns, inst string) {
	f, err := eddystone.MakeUIDFrame(ns, inst, eb.PowerLevel)
	if err != nil {
		panic(err)
	}
	eb.OutFrame = &f
}

func (eb *EddyStoneBeacon) addURL(url string) {
	f, err := eddystone.MakeURLFrame(url, eb.PowerLevel)
	if err != nil {
		panic(err)
	}
	eb.OutFrame = &f
}

// Start to Advertise your beacon, it is block API.
func (eb *EddyStoneBeacon) Advertise() {
	a := &gatt.AdvPacket{}
	a.AppendFlags(advFlagGeneralDiscoverable | advFlagLEOnly)
	a.AppendField(advTypeAllUUID16, eddystone.SvcUUIDBytes)
	a.AppendField(advTypeServiceData16, append(eddystone.SvcUUIDBytes, *eb.OutFrame...))

	fmt.Println(a.Len(), a)

	d, err := gatt.NewDevice(option.DefaultServerOptions...)
	if err != nil {
		log.Fatalf("Failed to open device, err: %s", err)
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
			d.Advertise(a)
		default:
			log.Println(s)
		}
	}

	d.Init(onStateChanged)
	select {}
}
