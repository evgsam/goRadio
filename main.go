package main

import (
	"goRadio/ic78civCmd"
	"goRadio/menu"
	"sync"
	"time"

	"go.bug.st/serial"
)

func main() {
	chRadioSettings := make(chan map[byte]string, 30)
	chSetData := make(chan map[byte]string, 30)
	var serialAccess sync.Mutex
	portCh := make(chan serial.Port)

	go ic78civCmd.CivCmdParser(portCh, &serialAccess, chRadioSettings, chSetData)

	menu.MainMenu(portCh, chRadioSettings, chSetData)
	for {
		time.Sleep(10 * time.Second)
	}

}
