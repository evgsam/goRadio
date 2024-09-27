package newmenu

import (
	component "goRadio/gocui-component"
	"goRadio/serialDataExchange"

	"github.com/jroimartin/gocui"
	"go.bug.st/serial"
)

type signup struct {
	*component.Form
	portCh chan serial.Port
}

func (s *signup) regist(g *gocui.Gui, v *gocui.View) error {
	if !s.Validate() {
		return nil
	}

	var text string

	for _, val := range s.GetSelectedRadios() {
		text = val
	}
	s.Close(g, v)
	s.portCh <- serialDataExchange.OpenSerialPort(19200, 8, text)
	return nil
}

func requireValidator(value string) bool {
	if value == "" {
		return false
	}
	return true
}

func spSelectMenu(g *gocui.Gui, portCh chan serial.Port) error {
	signup := &signup{
		component.NewForm(g, "Select Port", 0, 0, 0, 0), portCh,
	}

	signup.AddRadio(" ", 0).
		SetMode(component.VerticalMode).
		AddOptions(serialDataExchange.GetSerialPortList()...)
	signup.AddButton("Ok", signup.regist)
	signup.AddButton("Cancel", signup.Close)
	signup.Draw()
	return nil
}
