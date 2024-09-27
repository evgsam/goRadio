package newmenu

import (
	component "goRadio/gocui-component"
	"goRadio/serialDataExchange"

	"github.com/jroimartin/gocui"
)

type signup struct {
	*component.Form
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
	_ = serialDataExchange.OpenSerialPort(19200, 8, text)

	return nil
}

func requireValidator(value string) bool {
	if value == "" {
		return false
	}
	return true
}

func spSelectMenu(g *gocui.Gui) error {
	signup := &signup{
		component.NewForm(g, "Select Port", 0, 0, 0, 0),
	}

	signup.AddRadio(" ", 0).
		SetMode(component.VerticalMode).
		AddOptions(serialDataExchange.GetSerialPortList()...)
	signup.AddButton("Ok", signup.regist)
	signup.AddButton("Cancel", signup.Close)
	signup.Draw()
	return nil
}
