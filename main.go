package main

import (
	"goRadio/ic78civCmd"
	"goRadio/serialDataExchange"
	"sync"
	"time"

	"go.bug.st/serial"
)

func main() {
	var port serial.Port
	var serialAccess sync.Mutex
	port = serialDataExchange.OpenSerialPort(19200, 8)
	go ic78civCmd.CivCmdParser(port, &serialAccess)
	for {
		time.Sleep(10 * time.Second)
	}

}
