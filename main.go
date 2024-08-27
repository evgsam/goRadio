package main

import (
	"fmt"
	"goRadio/ic78civCmd"
	"goRadio/serialDataExchange"
	"sync"
	"time"

	"go.bug.st/serial"
)

func main() {

	/*if err := Oops(); err != nil {
		fmt.Println(err)
	}
	*/
	var port serial.Port
	var serialAccess sync.Mutex
	port = serialDataExchange.OpenSerialPort(19200, 8)
	//go ic78civCmd.DataPollingGorutine(port, &serialAccess)
	err := ic78civCmd.IC78connect(port, &serialAccess)
	if err != nil {
		fmt.Println("error:", err)
	}
	for {
		time.Sleep(10 * time.Second)
	}
}
