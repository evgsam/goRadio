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
	/*var text string

	for _, val := range s.GetSelectedRadios() {
		text = val
	}

	modal := component.NewModal(g, 0, 0, 30).SetText(text)
	modal.AddButton("OK", gocui.KeyEnter, func(g *gocui.Gui, v *gocui.View) error {
		modal.Close()
		s.SetCurrentItem(s.GetCurrentItem())
		return nil
	})
	modal.Draw()
	*/
	s.Close(g, v)

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

	// add button

	signup.AddButton("Ok", signup.regist)
	signup.AddButton("Cancel", signup.Close)
	signup.Draw()
	return nil
}
