Beacon Simulator: A simple beacon simulator (iBeacon/Eddystone) in Golang
==================

[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/kkdai/beacon/master/LICENSE)  [![GoDoc](https://godoc.org/github.com/kkdai/beacon?status.svg)](https://godoc.org/github.com/kkdai/beacon)  [![Build Status](https://travis-ci.org/kkdai/beacon.svg?branch=master)](https://travis-ci.org/kkdai/beacon)


This package is summarized golang beacon simulator with paypal/gatt package. It supports major two beacon as follow:

- iBeacon: Apple [iBeacon](https://developer.apple.com/ibeacon/)
- Eddystone: Google [Eddystone](https://github.com/google/eddystone)

Install library
---------------
`go get github.com/kkdai/beacon`


Install binary
---------------

- Eddystone: `go install github.com/kkdai/beacon/eddystone`
- iBeacon: `go install github.com/kkdai/beacon/ibeacon`

Simulator iBeacon
---------------

```go

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
```


Simulator Eddystone
---------------

```go

package main

import (
	. "github.com/kkdai/beacon"
)

func main() {
	ed := NewEddystoneURLBeacon("http://www.evanlin.com", -20)
	ed.Advertise()
}
```

Inspired by
---------------

- [github.com/go-ble/ble](github.com/go-ble/ble)
- [https://github.com/suapapa/go_eddystone](https://github.com/suapapa/go_eddystone)

Project52
---------------

It is one of my [project 52](https://github.com/kkdai/project52).


License
---------------

This package is licensed under MIT license. See LICENSE for details.

