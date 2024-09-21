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
	viewsNames     = make(map[byte]string)
	viewArray      = make([]viewsStruct, 0)
	err            error
	newGui         bool
	dataUpdateFlag bool
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

func dataToDisplay(ch chan map[byte]string, guiCh chan *gocui.Gui) {
	var g *gocui.Gui
	for {
		if !dataUpdateFlag {
			g = <-guiCh
			dataUpdateFlag = true
		}
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

func delNewView_(g *gocui.Gui, guiCh chan *gocui.Gui) error {
	newGui = false
	g.Close()

	g, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layoutMainMenu)

	if err := initKeybindings_(g, guiCh); err != nil {
		log.Panicln(err)
	}

	dataUpdateFlag = false
	guiCh <- g

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}
	return nil
}

func newView_(g *gocui.Gui, guiCh chan *gocui.Gui) error {
	if !newGui {
		newGui = true
		dataUpdateFlag = false
		for _, v := range viewArray {
			_ = delView_(g, v.name)
		}

		//inputMenu(g, guiCh)
		//widgets()
		g, err = gocui.NewGui(gocui.OutputNormal)
		if err != nil {
			log.Panicln(err)
		}
		defer g.Close()

		g.SetManagerFunc(layoutNewMenu)
		if err := initKeybindings_(g, guiCh); err != nil {
			log.Panicln(err)
		}

		newGui = true
		if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
			log.Panicln(err)
		}

		/*	g.Update(func(g *gocui.Gui) error {
				return nil
			})
		*/

	}

	return nil
}

func Menu(ch chan map[byte]string) {
	guiCn := make(chan *gocui.Gui)
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
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layoutMainMenu)

	if err := initKeybindings_(g, guiCn); err != nil {
		log.Panicln(err)
	}
	go dataToDisplay(ch, guiCn)
	guiCn <- g
	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
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

func delView_(g *gocui.Gui, name string) error {
	if err := g.DeleteView(name); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}
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

func initKeybindings_(g *gocui.Gui, guiCh chan *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlN, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return newView_(g, guiCh)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlB, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return delNewView_(g, guiCh)
		}); err != nil {
		return err
	}

	return nil
}
