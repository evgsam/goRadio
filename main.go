package main

import (
	newmenu "goRadio/new_menu"
	"time"
)

func main() {
	//var serialAccess sync.Mutex
	//go menu.Menu(&serialAccess)
	newmenu.NewMenu()
	for {
		time.Sleep(10 * time.Second)
	}

}
