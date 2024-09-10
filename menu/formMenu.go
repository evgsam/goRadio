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

func layoutF(g *gocui.Gui) error {
	/*if v, err := g.SetView("v9", 0, 8, 50, 12); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " IC-78 Set "
	}
	*/
	if v, err := g.SetView("v10", 1, 9, 16, 11); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " freque set "
	}

	if v, err := g.SetView("v11", 17, 9, 27, 11); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " mode"
	}

	if v, err := g.SetView("v12", 28, 9, 38, 11); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " att"
	}

	if v, err := g.SetView("v13", 39, 9, 49, 11); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " preamp"
	}
	/*
		component.NewSelect(g, "mode_set", 17, 9, 0, 6).
			AddOptions("LSB", "USB", "CW", "RTTY", "AM").
			Draw()
		component.NewSelect(g, "att_set", 28, 9, 0, 6).
			AddOptions("ON", "OFF").
			Draw()
	*/
	return nil
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
	//g.SetManagerFunc(layoutF)

	// new form
	signup1 := &signupF{
		component.NewForm(g, " IC-78 Set ", 0, 8, 50, 12),
	}
	signup1.Draw()

	signupFreq := &signupF{
		component.NewForm(g, " Freque Set ", 1, 9, 10, 1),
	}
	//signup.AddInputField_(" ", 11, 9, 11, 18).
	//		AddValidate("required input", requireValidatorF)
	signupFreq.Draw()

	signupMode := &signupF{
		component.NewForm(g, " Mode ", 13, 9, 10, 1),
	}
	//signup2.AddInputField_(" ", 11, 9, 11, 18).
	//	AddValidate("required input", requireValidatorF)
	signupMode.Draw()

	signupAtt := &signupF{
		component.NewForm(g, " ATT ", 25, 9, 10, 1),
	}
	//signup2.AddInputField_(" ", 11, 9, 11, 18).
	//	AddValidate("required input", requireValidatorF)
	signupAtt.Draw()

	signupPreamp := &signupF{
		component.NewForm(g, " Preamp ", 37, 9, 10, 1),
	}
	//signup2.AddInputField_(" ", 11, 9, 11, 18).
	//	AddValidate("required input", requireValidatorF)
	signupPreamp.Draw()

	/*signup.AddInputField("Last Name", 11, 18).
		AddValidate("required input", requireValidatorF)

	signup.AddInputField("Password", 11, 18).
		AddValidate("required input", requireValidatorF).
		SetMask().
		SetMaskKeybinding(gocui.KeyCtrlA)

		// add checkbox
	signup.AddCheckBox("Age 18+", 11)

	// add select
	signup.AddSelect("Language", 11, 10).AddOptions("Japanese", "English", "Chinese")

	// add radios
	signup.AddRadio("Country", 11).
		SetMode(component.VerticalMode).
		AddOptions("Japan", "America", "China")

	// add button
	signup.AddButton("Regist", signup.registF)
	signup.AddButton("Cancel", quitF)
	*/

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

	text += fmt.Sprintf("radio: %s\n", s.GetSelectedRadios())

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
