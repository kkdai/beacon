package main

import (
	. "github.com/kkdai/beacon"
)

func main() {
	ed := NewEddystoneURLBeacon("http://www.evanlin.com", -20)
	ed.Advertise()
}
