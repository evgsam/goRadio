package menu

import (
	"errors"
	"fmt"
	"log"

	"github.com/jroimartin/gocui"

	component "goRadio/gocui-component"
	//component "github.com/skanehira/gocui-component"
)

type signupF struct {
	*component.Form
}

func InputMenuForm() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quitF); err != nil {
		log.Panicln(err)
	}

	// new form
	signup := &signupF{
		component.NewForm(g, " IC-78 Set ", 0, 8, 20, 10),
	}
	signup.AddInputField_("Freque", 0, 9, 6, 10).
		AddValidate("required input", requireValidatorF)

	// add select
	signup.AddSelect_("Mode", 18, 9, 4, 5).AddOptions("LSB", "USB", "CW", "RTTY", "AM")
	signup.AddSelect_("ATT", 29, 9, 3, 4).AddOptions("ON", "OFF")
	signup.AddSelect_("Preamp", 37, 9, 6, 6).AddOptions("P.AMP", "OFF")

	signup.AddInputField_("AF", 7, 11, 4, 5).
		AddValidate("required input", requireValidatorF)

	signup.AddInputField_("RF", 17, 11, 4, 5).
		AddValidate("required input", requireValidatorF)

	signup.AddInputField_("SQL", 27, 11, 4, 5).
		AddValidate("required input", requireValidatorF)
	// add button
	signup.AddButton("SEND Settins", signup.registF)
	signup.AddButton("Cancel", quitF)

	signup.Draw()
	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}
	//g.Close()

}

func (s *signupF) registF(g *gocui.Gui, v *gocui.View) error {
	if !s.Validate() {
		return nil
	}

	var text string

	for label, ftext := range s.GetFieldTexts() {
		text += fmt.Sprintf("%s: %s\n", label, ftext)
	}

	for label, state := range s.GetCheckBoxStates() {
		text += fmt.Sprintf("%s: %t\n", label, state)
	}

	for label, opt := range s.GetSelectedOpts() {
		text += fmt.Sprintf("%s: %s\n", label, opt)
	}

	modal := component.NewModal(g, 0, 0, 30).SetText(text)
	modal.AddButton("OK", gocui.KeyEnter, func(g *gocui.Gui, v *gocui.View) error {
		modal.Close()
		s.SetCurrentItem(s.GetCurrentItem())
		return nil
	})

	modal.Draw()

	return nil
}

func quitF(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func requireValidatorF(value string) bool {
	if value == "" {
		return false
	}
	return true
}
