package menu

import (
	"errors"
	"log"

	"github.com/jroimartin/gocui"

	component "goRadio/gocui-component"
	//component "github.com/skanehira/gocui-component"
)

type signupF struct {
	*component.Form
	ch chan map[byte]string
}

func updateMenu(g *gocui.Gui, signup *signupF) {
	for {
		g.Update(func(g *gocui.Gui) error {
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
			return nil
		})
	}
}

func InputMenuForm(ch chan map[byte]string) {
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
		component.NewForm(g, " IC-78 Set ", 0, 8, 20, 10), ch,
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
	m := make(map[byte]string)
	if !s.Validate() {
		return nil
	}

	for label, ftext := range s.GetFieldTexts() {
		switch label {
		case "Freque":
			m[byte(freqRead)] = ftext
		case "AF":
			m[byte(af)] = ftext
		case "RF":
			m[byte(rf)] = ftext
		case "SQL":
			m[byte(sql)] = ftext

		}
	}

	for label, opt := range s.GetSelectedOpts() {
		switch label {
		case "Mode":
			m[byte(mode)] = opt
		case "ATT":
			m[byte(att)] = opt
		case "Preamp":
			m[byte(preamp)] = opt
		}
	}
	s.ch <- m
	//updateMenu(g, s)
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
