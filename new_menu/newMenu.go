package newmenu

import (
	"errors"
	"fmt"
	"goRadio/serialDataExchange"
	"log"
	"strconv"

	"github.com/jroimartin/gocui"
)

var (
	infoViewArray   = make([]viewsStruct, 0)
	hotkeyViewArray = make([]viewsStruct, 0)
	inputViewArray  = make([]viewsStruct, 0)
)

const (
	F2_title = "F2 Serial port select"
	F2_input = "Input:"
	F3_title = "F3 Enter freque"
)

type viewsStruct struct {
	name           string
	x0, y0, x1, y1 int
	value          string
	bottomFlag     bool
}

func viewArrayFilling() {
	portsList := serialDataExchange.GetSerialPortList()
	t := ""
	for i, v := range portsList {
		t += "Portn #" + strconv.Itoa(i) + ":" + v + "\n"
	}
	t += "Plese enter port number:"

	inputViewArray = []viewsStruct{
		{name: F2_title, x0: 6, y0: 1, x1: 33, y1: 7, value: t, bottomFlag: true},
		{name: F2_input, x0: 6, y0: 6, x1: 33, y1: 7, value: t, bottomFlag: false},
		{name: F3_title, x0: 17, y0: 1, x1: 40, y1: 6, value: "Freque input", bottomFlag: true},
	}

	hotkeyViewArray = []viewsStruct{
		{name: "Hotkey for change settings", x0: 0, y0: 0, x1: 50, y1: 7},
		{name: "F1", x0: 1, y0: 1, x1: 8, y1: 3, value: "help"},
		{name: "F2", x0: 9, y0: 1, x1: 16, y1: 3, value: "s.port"},
		{name: "F3", x0: 17, y0: 1, x1: 24, y1: 3, value: "freque"},
		{name: "F4", x0: 25, y0: 1, x1: 32, y1: 3, value: "mode"},
		{name: "F5", x0: 33, y0: 1, x1: 40, y1: 3, value: "ATT"},
		{name: "F6", x0: 41, y0: 1, x1: 48, y1: 3, value: "preamp"},

		{name: "F7", x0: 1, y0: 4, x1: 8, y1: 6, value: "AF-"},
		{name: "F8", x0: 9, y0: 4, x1: 16, y1: 6, value: "AF+"},
		{name: "F9", x0: 17, y0: 4, x1: 24, y1: 6, value: "RF-"},
		{name: "F10", x0: 25, y0: 4, x1: 32, y1: 6, value: "RF+"},
		{name: "F11", x0: 33, y0: 4, x1: 40, y1: 6, value: "SQL-"},
		{name: "F12", x0: 41, y0: 4, x1: 48, y1: 6, value: "SQL+"},
	}
	infoViewArray = []viewsStruct{
		{name: "IC-78Information", x0: 0, y0: 8, x1: 50, y1: 15},
		{name: "status", x0: 1, y0: 9, x1: 16, y1: 11},
		{name: "mode", x0: 17, y0: 9, x1: 27, y1: 11},
		{name: "ATT", x0: 28, y0: 9, x1: 38, y1: 11},
		{name: "preamp", x0: 39, y0: 9, x1: 49, y1: 11},
		{name: "freque", x0: 1, y0: 12, x1: 16, y1: 14},
		{name: "AF", x0: 17, y0: 12, x1: 27, y1: 14},
		{name: "RF", x0: 28, y0: 12, x1: 38, y1: 14},
		{name: "SQL", x0: 39, y0: 12, x1: 49, y1: 14},
	}

}

// ===========================
var text_ string

type Input struct {
	name      string
	x, y      int
	w         int
	maxLength int
}

func NewInput(name string, x, y, w, maxLength int) *Input {
	return &Input{name: name, x: x, y: y, w: w, maxLength: maxLength}
}

func (i *Input) Layout(g *gocui.Gui) error {
	v, err := g.SetView(i.name, i.x, i.y, i.x+i.w, i.y+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editor = i
		v.Editable = true
	}
	return nil
}

func (i *Input) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	cx, _ := v.Cursor()
	ox, _ := v.Origin()
	limit := ox+cx+1 > i.maxLength
	switch {
	case ch != 0 && mod == 0 && !limit:
		text_ += string(ch)
		v.EditWrite(ch)
	case key == gocui.KeySpace && !limit:
		text_ += string(ch)
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		text_ = text_[:len(text_)-1]
		v.EditDelete(true)
	}
}

func SetFocus(name string) func(g *gocui.Gui) error {
	return func(g *gocui.Gui) error {
		_, err := g.SetCurrentView(name)
		return err
	}
}

//===========================

func NewMenu() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	viewArrayFilling()

	g.SetManagerFunc(layoutSetView)

	if err := initKeybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layoutSetView(g *gocui.Gui) error {
	for _, v := range hotkeyViewArray {
		_ = setView(g, v.name, v.x0, v.y0, v.x1, v.y1, v.value, false)
	}
	for _, v := range infoViewArray {
		_ = setView(g, v.name, v.x0, v.y0, v.x1, v.y1, "", false)
	}
	for _, v := range inputViewArray {
		_ = setView(g, v.name, v.x0, v.y0, v.x1, v.y1, v.value, v.bottomFlag)
	}
	return nil
}

func setView(g *gocui.Gui, name string, x0, y0, x1, y1 int, value string, flag bool) error {
	if v, err := g.SetView(name, x0, y0, x1, y1); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		if name == F2_input {
			//v.Frame = false
			v.Title = ""
			//fmt.Fprintln(v, "port")
			v.Editable = true
			_, _ = g.SetCurrentView(name)
		} else {
			v.Title = name
		}

		fmt.Fprintln(v, value)
		if flag {
			_, err = g.SetViewOnBottom(name)
		}
	}
	return nil
}

func changeBottomFlag(p *viewsStruct) {
	p.bottomFlag = !p.bottomFlag
}

func viewTopOrBottom(g *gocui.Gui, flag bool, name string) {
	if !flag {
		_, _ = g.SetViewOnTop(name)
	} else {
		_, _ = g.SetViewOnBottom(name)
	}
}

func initKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyF2, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		var ind int
		for i, v := range inputViewArray {
			if v.name == F2_title {
				ind = i
			}
		}
		p := &inputViewArray[ind]
		p.bottomFlag = !p.bottomFlag
		for _, v := range hotkeyViewArray {
			viewTopOrBottom(g, v.bottomFlag, v.name)
		}
		for _, v := range infoViewArray {
			viewTopOrBottom(g, v.bottomFlag, v.name)
		}
		for _, v := range inputViewArray {
			viewTopOrBottom(g, v.bottomFlag, v.name)
		}
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyF3, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		var ind int
		for i, v := range inputViewArray {
			if v.name == F3_title {
				ind = i
			}
		}
		p := &inputViewArray[ind]
		p.bottomFlag = !p.bottomFlag
		for _, v := range hotkeyViewArray {
			viewTopOrBottom(g, v.bottomFlag, v.name)
		}
		for _, v := range infoViewArray {
			viewTopOrBottom(g, v.bottomFlag, v.name)
		}
		for _, v := range inputViewArray {
			viewTopOrBottom(g, v.bottomFlag, v.name)
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func viewsUpdate(g *gocui.Gui, name string, flag bool) error {
	_, err := g.View(name)
	if err != nil {
		return err
	}
	return nil
}
