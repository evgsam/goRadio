package newmenu

import (
	component "goRadio/gocui-component"

	"github.com/jroimartin/gocui"
)

type signupFreqSet struct {
	*component.Form
}

func (s *signupFreqSet) inputRegist(g *gocui.Gui, v *gocui.View) error {
	if !s.Validate() {
		return nil
	}

	/*var text string

	for _, val := range s.GetFieldTexts() {
		text = val
	}

	*/
	s.Close(g, v)
	return nil
}

func freqSetMenu(g *gocui.Gui) error {
	signupFreq := &signupFreqSet{
		component.NewForm(g, "Freque set", 0, 0, 0, 0),
	}
	signupFreq.AddInputField("Freque, hz", 11, 9).
		AddValidate("required input", requireValidator)
	signupFreq.AddButton("Ok", signupFreq.inputRegist)

	signupFreq.Draw()
	return nil
}
