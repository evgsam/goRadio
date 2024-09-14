package menu

import (
	"errors"
	"fmt"
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

type commandName int

const (
	freqRead commandName = iota
	taddr
	mode
	att
	af
	rf
	sql
	preamp
	status
)

var (
	views_ = []string{}
)

func dataToDisplay(g *gocui.Gui, ch chan map[byte]string) {
	for {
		m := <-ch
		g.Update(func(g *gocui.Gui) error {
			for key, value := range m {
				switch key {
				case byte(status):
					v, err := g.View(Status)
					if err != nil {
						return err
					}
					v.Clear()
					fmt.Fprintln(v, value)

				case byte(mode):
					v, err := g.View(Mode)
					if err != nil {
						return err
					}
					v.Clear()
					fmt.Fprintln(v, value)

				case byte(att):
					v, err := g.View(ATT)
					if err != nil {
						return err
					}
					v.Clear()
					fmt.Fprintln(v, value)

				case byte(preamp):
					v, err := g.View(Preamp)
					if err != nil {
						return err
					}
					v.Clear()
					fmt.Fprintln(v, value)

				case byte(freqRead):
					v, err := g.View(Freque)
					if err != nil {
						return err
					}
					v.Clear()
					fmt.Fprintln(v, m[byte(freqRead)])

				case byte(af):
					v, err := g.View(AF)
					if err != nil {
						return err
					}
					v.Clear()
					fmt.Fprintln(v, value)

				case byte(rf):
					v, err := g.View(RF)
					if err != nil {
						return err
					}
					v.Clear()
					fmt.Fprintln(v, value)

				case byte(sql):
					v, err := g.View(SQL)
					if err != nil {
						return err
					}
					v.Clear()
					fmt.Fprintln(v, value)
				}
			}
			return nil
		})
	}
}

func newView_(g *gocui.Gui) error {
	_, err := g.SetView("set", 0, 0, 15, 15)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	if _, err := g.SetCurrentView("set"); err != nil {
		return err
	}
	return nil
}

func Menu(ch chan map[byte]string) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := initKeybindings_(g); err != nil {
		log.Panicln(err)
	}
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
	views_ = append(views_, "v0")

	if v, err := g.SetView(Status, 1, 1, 16, 3); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " status "
		fmt.Fprint(v, "Disconnected")
	}
	views_ = append(views_, Status)

	if v, err := g.SetView(Mode, 17, 1, 27, 3); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " mode "
	}
	views_ = append(views_, Mode)

	if v, err := g.SetView(ATT, 28, 1, 38, 3); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " ATT "
	}
	views_ = append(views_, ATT)

	if v, err := g.SetView(Preamp, 39, 1, 49, 3); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " preamp "
	}
	views_ = append(views_, Preamp)

	if v, err := g.SetView(Freque, 1, 4, 16, 6); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " freque "
	}
	views_ = append(views_, Freque)

	if v, err := g.SetView(AF, 17, 4, 27, 6); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " AF "
	}
	views_ = append(views_, AF)

	if v, err := g.SetView(RF, 28, 4, 38, 6); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " RF "
	}
	views_ = append(views_, RF)

	if v, err := g.SetView(SQL, 39, 4, 49, 6); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " SQL "
	}
	views_ = append(views_, SQL)

	return nil
}

func initKeybindings_(g *gocui.Gui) error {
	if err := g.SetKeybinding("", 'n', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return newView_(g)
		}); err != nil {
		return err
	}
	return nil
}
