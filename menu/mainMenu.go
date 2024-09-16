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
	mainViews
)

var (
	//displayAccessMutex sync.Mutex
	viewsNames = make(map[byte]string)
	viewArray  = make([]viewsStruct, 0)
	mainGui    *gocui.Gui
	setGui     *gocui.Gui
	err        error
	newGui     bool
)

type viewsStruct struct {
	cmd            commandName
	name           string
	x0, y0, x1, y1 int
}

func viewsValueUpdate(g *gocui.Gui, name, value string) error {
	v, err := g.View(name)
	if err != nil {
		return err
	}
	v.Clear()
	fmt.Fprintln(v, value)
	return nil
}

func dataToDisplay(g *gocui.Gui, ch chan map[byte]string) {
	for {
		m := <-ch
		g.Update(func(g *gocui.Gui) error {
			for key, value := range m {
				switch key {
				case byte(status):
					_ = viewsValueUpdate(g, viewsNames[byte(status)], value)
				case byte(mode):
					_ = viewsValueUpdate(g, viewsNames[byte(mode)], value)
				case byte(att):
					_ = viewsValueUpdate(g, viewsNames[byte(att)], value)
				case byte(preamp):
					_ = viewsValueUpdate(g, viewsNames[byte(preamp)], value)
				case byte(freqRead):
					_ = viewsValueUpdate(g, viewsNames[byte(freqRead)], value)
				case byte(af):
					_ = viewsValueUpdate(g, viewsNames[byte(af)], value)
				case byte(rf):
					_ = viewsValueUpdate(g, viewsNames[byte(rf)], value)
				case byte(sql):
					_ = viewsValueUpdate(g, viewsNames[byte(sql)], value)
				}
			}
			return nil
		})
	}
}

func delNewView_(g *gocui.Gui) error {
	newGui = false
	setGui.Close()

	mainGui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer mainGui.Close()

	mainGui.SetManagerFunc(layoutMainMenu)

	if err := initKeybindings_(mainGui); err != nil {
		log.Panicln(err)
	}
	if err := mainGui.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}
	return nil
}

func newView_(gui *gocui.Gui) error {
	if !newGui {
		mainGui.Close()
		setGui, err = gocui.NewGui(gocui.OutputNormal)
		if err != nil {
			log.Panicln(err)
		}
		defer setGui.Close()

		setGui.SetManagerFunc(layoutNewMenu)

		if err := initKeybindings_(setGui); err != nil {
			log.Panicln(err)
		}

		newGui = true
		if err := setGui.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
			log.Panicln(err)
		}
	}

	return nil
}

func Menu(ch chan map[byte]string) {
	viewArray = []viewsStruct{
		{cmd: mainViews, name: "IC-78Information", x0: 0, y0: 0, x1: 50, y1: 7},
		{cmd: status, name: "status", x0: 1, y0: 1, x1: 16, y1: 3},
		{cmd: mode, name: "mode", x0: 17, y0: 1, x1: 27, y1: 3},
		{cmd: att, name: "ATT", x0: 28, y0: 1, x1: 38, y1: 3},
		{cmd: preamp, name: "preamp", x0: 39, y0: 1, x1: 49, y1: 3},
		{cmd: freqRead, name: "freque", x0: 1, y0: 4, x1: 16, y1: 6},
		{cmd: af, name: "AF", x0: 17, y0: 4, x1: 27, y1: 6},
		{cmd: rf, name: "RF", x0: 28, y0: 4, x1: 38, y1: 6},
		{cmd: sql, name: "SQL", x0: 39, y0: 4, x1: 49, y1: 6},
	}
	for _, v := range viewArray {
		viewsNames[byte(v.cmd)] = v.name
	}
	mainGui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer mainGui.Close()

	mainGui.SetManagerFunc(layoutMainMenu)

	if err := initKeybindings_(mainGui); err != nil {
		log.Panicln(err)
	}
	go dataToDisplay(mainGui, ch)
	if err := mainGui.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}

}

func setView(g *gocui.Gui, name string, x0, y0, x1, y1 int) error {
	if v, err := g.SetView(name, x0, y0, x1, y1); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = name
	}
	return nil
}

func layoutMainMenu(g *gocui.Gui) error {
	for _, v := range viewArray {
		_ = setView(g, v.name, v.x0, v.y0, v.x1, v.y1)
	}
	return nil
}

func layoutNewMenu(g *gocui.Gui) error {
	_ = setView(g, "set", 0, 0, 50, 7)
	return nil
}

func initKeybindings_(g *gocui.Gui) error {
	if err := g.SetKeybinding("", 'n', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return newView_(g)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'b', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return delNewView_(g)
		}); err != nil {
		return err
	}

	return nil
}
