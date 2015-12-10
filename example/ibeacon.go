package main

import (
	. "github.com/kkdai/beacon"
)

func main() {
	ib := NewIBeacon("AA6062F098CA42118EC4193EB73CCEB6", "Gopher", -59)
	ib.SetiBeaconVersion(1, 2)
	ib.AddCountService()
	ib.AddBatteryService()
	ib.Advertise()
}
