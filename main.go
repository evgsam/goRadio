package main

import (
	"goRadio/ic78civCmd"
	"goRadio/serialDataExchange"
	"sync"
	"time"
)

func main() {
	//menu.Exemple()
	port := serialDataExchange.OpenSerialPort(19200, 8)
	//ch := make(chan map[byte]string, 3)
	//go ic78civCmd.IC78civCmdSet(port, ch)
	//menu.InputMenuForm(ch)

	var serialAccess sync.Mutex

	go ic78civCmd.CivCmdParser(port, &serialAccess)

	for {
		time.Sleep(10 * time.Second)
	}

}
