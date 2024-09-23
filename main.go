package main

import (
	"goRadio/menu"
	"sync"
	"time"
)

func main() {
	var serialAccess sync.Mutex
	go menu.Menu(&serialAccess)
	for {
		time.Sleep(10 * time.Second)
	}

}
