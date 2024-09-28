package newmenu

import (
	component "goRadio/gocui-component"
	"goRadio/serialDataExchange"

	"github.com/jroimartin/gocui"
	"go.bug.st/serial"
)

type signupSpSelect struct {
	*component.Form
	portCh chan serial.Port
}

func (s *signupSpSelect) close(g *gocui.Gui, v *gocui.View) error {
	spMenuActive = false
	s.Close(g, v)
	return nil
}

func (s *signupSpSelect) radioRegist(g *gocui.Gui, v *gocui.View) error {
	spMenuActive = false
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

func spSelectMenu(g *gocui.Gui, portCh chan serial.Port) error {
	spMenuActive = true
	signup := &signupSpSelect{
		component.NewForm(g, "Select Port", 0, 0, 0, 0), portCh,
	}

	signup.AddRadio(" ", 0).
		SetMode(component.VerticalMode).
		AddOptions(serialDataExchange.GetSerialPortList()...)
	signup.AddButton("Ok", signup.radioRegist)
	signup.AddButton("Cancel", signup.close)
	signup.Draw()
	return nil
}
