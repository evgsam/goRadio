package menu

import (
	"errors"
	"fmt"
	datastruct "goRadio/dataStruct"
	"log"

	"github.com/jroimartin/gocui"
	//"github.com/awesome-gocui/gocui"
)

const (
	Status string = "status"
	Mode   string = "mode"
	ATT    string = "ATT"
	Preamp string = "preamp"
	Freque string = "freque"
	AF     string = "AF"
	RF     string = "RF"
	SQL    string = "SQL"
)

func dataToDisplay(g *gocui.Gui, ch chan *datastruct.RadioSettings) {
	for {
		myRadioSettings := <-ch
		g.Update(func(g *gocui.Gui) error {
			v, err := g.View(Status)
			if err != nil {
				return err
			}
			v.Clear()
			fmt.Fprintln(v, myRadioSettings.Status)
			v, err = g.View(Mode)
			if err != nil {
				return err
			}
			v.Clear()
			fmt.Fprintln(v, myRadioSettings.Mode)

			v, err = g.View(ATT)
			if err != nil {
				return err
			}
			v.Clear()
			fmt.Fprintln(v, myRadioSettings.ATT)

			v, err = g.View(Preamp)
			if err != nil {
				return err
			}
			v.Clear()
			fmt.Fprintln(v, myRadioSettings.Preamp)

			v, err = g.View(Freque)
			if err != nil {
				return err
			}
			v.Clear()
			fmt.Fprintln(v, myRadioSettings.Freque)

			v, err = g.View(AF)
			if err != nil {
				return err
			}
			v.Clear()
			fmt.Fprintln(v, myRadioSettings.AF)

			v, err = g.View(RF)
			if err != nil {
				return err
			}
			v.Clear()
			fmt.Fprintln(v, myRadioSettings.RF)

			v, err = g.View(SQL)
			if err != nil {
				return err
			}
			v.Clear()
			fmt.Fprintln(v, myRadioSettings.SQL)

			return nil
		})
	}
}

func Menu(ch chan *datastruct.RadioSettings) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	g.SetManagerFunc(layout)
	go dataToDisplay(g, ch)

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}
	g.Close()

}

func layout(g *gocui.Gui) error {
	if v, err := g.SetView("v0", 0, 0, 50, 7); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " IC-78 Information "
	}

	if v, err := g.SetView(Status, 1, 1, 16, 3); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " status "
		fmt.Fprint(v, "Disconnected")
	}

	if v, err := g.SetView(Mode, 17, 1, 27, 3); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " mode "
	}

	if v, err := g.SetView(ATT, 28, 1, 38, 3); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " ATT "
	}

	if v, err := g.SetView(Preamp, 39, 1, 49, 3); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " preamp "
	}

	if v, err := g.SetView(Freque, 1, 4, 16, 6); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " freque "
	}

	if v, err := g.SetView(AF, 17, 4, 27, 6); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " AF "
	}

	if v, err := g.SetView(RF, 28, 4, 38, 6); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " RF "
	}

	if v, err := g.SetView(SQL, 39, 4, 49, 6); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " SQL "
	}

	/*component.NewSelect(g, "Mode:", 2, 9, 0, 5).
		AddOptions("LSB", "USB", "AM", "RTTY", "CW").
		Draw()

	if v, err := g.SetView("v9", 0, 8, 50, 12); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " IC-78 Set "
	}

	if v, err := g.SetView("v10", 1, 9, 16, 11); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " freque set "
	}
	*/
	/*if v, err := g.SetView("v3", 0, 4, 30, 7, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = "Transiver mode"
	}
	*/
	return nil
}
