package main

import (
	. "github.com/kkdai/beacon"
)

func main() {
	ed := NewEddystoneUIDBeacon("00010203040506070809", "aabbccddeeff", -20)
	ed.Advertise()
}
