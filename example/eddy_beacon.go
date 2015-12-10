package main

import (
	. "github.com/kkdai/beacon"
)

func main() {
	ed := NewEddystoneBeacon(-20)
	ed.AddURL("http://www.evanlin.com")
	ed.Advertise()
}
