package menu

import (
	component "goRadio/gocui-component"

	"github.com/jroimartin/gocui"
)

type signupFreqSet struct {
	*component.Form
	chDataSet chan map[byte]string
}

func (s *signupFreqSet) inputRegist(g *gocui.Gui, v *gocui.View) error {
	freqMenuActive = false
	if !s.Validate() {
		return nil
	}
	var text string
	for _, val := range s.GetFieldTexts() {
		text = val
	}
	s.Close(g, v)
	s.chDataSet <- map[byte]string{
		byte(freqRead): text,
	}
	return nil
}

func freqSetMenu(g *gocui.Gui, chDataSet chan map[byte]string) error {
	freqMenuActive = true
	signupFreq := &signupFreqSet{
		component.NewForm(g, "Freque set", 0, 0, 0, 0), chDataSet,
	}
	signupFreq.AddInputField("Freque, hz", 11, 9).
		AddValidate("required input", requireValidator)
	signupFreq.AddButton("Ok", signupFreq.inputRegist)

	signupFreq.Draw()
	return nil
}
