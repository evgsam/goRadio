package main

import (
	"fmt"
	"log"

	"go.bug.st/serial"
)

func OpenSerialPort() serial.Port {
	ports, err := serial.GetPortsList()
	var portsnum int
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}
	fmt.Print("Ports list: \n")
	for _, port := range ports {
		fmt.Printf("Port #%d: %v\n", portsnum, port)
		portsnum++
	}

	if len(ports) > 1 {
		fmt.Print("Please, select port:")
		fmt.Scan(&portsnum)
	} else {
		portsnum = 0
	}
	mode := &serial.Mode{
		BaudRate: 19200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	port, err := serial.Open(ports[portsnum], mode)
	if err != nil {
		log.Fatal(err)
	}
	return port
}
