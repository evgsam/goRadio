package main

import (
	"fmt"
	"goRadio/ic78civCmd"
	"goRadio/serialDataExchange"

	"go.bug.st/serial"
)

func printByte(data []byte) {
	for _, value := range data {
		fmt.Printf("%#x ", value)
	}
	fmt.Println()
}

func main() {
	var port serial.Port
	port = serialDataExchange.OpenSerialPort(19200, 8)
	ic78civCmd.IC78connect(port)
}
