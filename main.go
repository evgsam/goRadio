package main

import (
	"goRadio/ic78civCmd"
	newmenu "goRadio/new_menu"
	"sync"
	"time"

	"go.bug.st/serial"
)

func main() {
	chRadioSettings := make(chan map[byte]string, 30)
	var serialAccess sync.Mutex
	portCh := make(chan serial.Port)

	go ic78civCmd.CivCmdParser(portCh, &serialAccess, chRadioSettings)
	newmenu.NewMenu(portCh, chRadioSettings)
	for {
		time.Sleep(10 * time.Second)
	}

}
