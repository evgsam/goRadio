package menu

import (
	component "goRadio/gocui-component"

	"github.com/jroimartin/gocui"
)

type signup struct {
	*component.Form
	ch chan []string
}

func serialPortSelectMenu(g *gocui.Gui, ports []string, ch chan []string) {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		panic(err)
	}
	signup := &signup{
		component.NewForm(g, "Select Port", 0, 0, 0, 0), ch,
	}
	signup.AddRadio(" ", 0).
		SetMode(component.VerticalMode).
		AddOptions(ports...)

	// add button
	signup.AddButton("Ok", signup.regist)
	signup.AddButton("Cancel", quit)
	signup.Draw()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
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
	s.Close(g, v)
	quit(g, v)
	s.ch <- []string{text}
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
