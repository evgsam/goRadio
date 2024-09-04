package menu

import (
	component "goRadio/gocui-component"

	"github.com/jroimartin/gocui"
)

type signup struct {
	*component.Form
	ch chan []string
}

func SerialPortSelectMenu(ch chan []string) {
	portList := <-ch

	gui, err := gocui.NewGui(gocui.OutputNormal)

	if err != nil {
		panic(err)
	}
	defer gui.Close()

	if err := gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		panic(err)
	}

	// new form
	//maxX, maxY := gui.Size()
	signup := &signup{
		component.NewForm(gui, "Select Port", 0, 0, 0, 0), ch,
	}
	signup.AddRadio(" ", 0).
		SetMode(component.VerticalMode).
		AddOptions(portList...)

	// add button
	signup.AddButton("Ok", signup.regist)
	signup.AddButton("Cancel", quit)
	signup.Draw()

	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}
}

func (s *signup) regist(g *gocui.Gui, v *gocui.View) error {
	if !s.Validate() {
		return nil
	}
	var text string

	for _, val := range s.GetSelectedRadios() {
		text = val
	}
	s.ch <- []string{text}
	s.Close(g, v)
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func requireValidator(value string) bool {
	if value == "" {
		return false
	}
	return true
}
