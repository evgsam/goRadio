/******************************************/
//Основное меню - вывожу значения, прочитанные из радиостанции
/*****************************************/

package menu

import (
	"errors"
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
	"go.bug.st/serial"
)

const (
	F2_title = "F2 Serial port select"
	F3_title = "F3 Enter freque"
)

func viewArrayFilling() {
	hotkeyViewArray = []viewsStruct{
		{name: "Hotkey for change settings", x0: 0, y0: 0, x1: 50, y1: 7},
		{name: "F1", x0: 1, y0: 1, x1: 8, y1: 3, value: "help"},
		{name: "S", x0: 9, y0: 1, x1: 16, y1: 3, value: "s.port"},
		{name: "F", x0: 17, y0: 1, x1: 24, y1: 3, value: "freque"},
		{name: "M", x0: 25, y0: 1, x1: 32, y1: 3, value: "mode"},
		{name: "A", x0: 33, y0: 1, x1: 40, y1: 3, value: "ATT"},
		{name: "P", x0: 41, y0: 1, x1: 48, y1: 3, value: "preamp"},

		{name: "Q", x0: 1, y0: 4, x1: 8, y1: 6, value: "AF-"},
		{name: "W", x0: 9, y0: 4, x1: 16, y1: 6, value: "AF+"},
		{name: "E", x0: 17, y0: 4, x1: 24, y1: 6, value: "RF-"},
		{name: "R", x0: 25, y0: 4, x1: 32, y1: 6, value: "RF+"},
		{name: "Z", x0: 33, y0: 4, x1: 40, y1: 6, value: "SQL-"},
		{name: "X", x0: 41, y0: 4, x1: 48, y1: 6, value: "SQL+"},
	}

	infoViewArray = []viewsStruct{
		{cmd: mainViews, name: "IC-78Information", x0: 0, y0: 8, x1: 50, y1: 15},
		{cmd: status, name: "status", x0: 1, y0: 9, x1: 16, y1: 11},
		{cmd: mode, name: "mode", x0: 17, y0: 9, x1: 27, y1: 11},
		{cmd: att, name: "ATT", x0: 28, y0: 9, x1: 38, y1: 11},
		{cmd: preamp, name: "preamp", x0: 39, y0: 9, x1: 49, y1: 11},
		{cmd: freqRead, name: "freque", x0: 1, y0: 12, x1: 16, y1: 14},
		{cmd: af, name: "AF", x0: 17, y0: 12, x1: 27, y1: 14},
		{cmd: rf, name: "RF", x0: 28, y0: 12, x1: 38, y1: 14},
		{cmd: sql, name: "SQL", x0: 39, y0: 12, x1: 49, y1: 14},
	}

	for _, v := range infoViewArray {
		viewsNames[byte(v.cmd)] = v.name
	}
}

func MainMenu(portCh chan serial.Port, chRadioSettings chan map[byte]string, chDataSet chan map[byte]string) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	viewArrayFilling()
	g.SetManagerFunc(layoutSetView)

	go dataToDisplay(g, chRadioSettings)
	if err := initKeybindings(g, portCh, chDataSet); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layoutSetView(g *gocui.Gui) error {
	for _, v := range hotkeyViewArray {
		_ = setView(g, v.name, v.x0, v.y0, v.x1, v.y1, v.value)
	}
	for _, v := range infoViewArray {
		_ = setView(g, v.name, v.x0, v.y0, v.x1, v.y1, "")
	}
	for _, v := range inputViewArray {
		_ = setView(g, v.name, v.x0, v.y0, v.x1, v.y1, v.value)
	}
	return nil
}

func setView(g *gocui.Gui, name string, x0, y0, x1, y1 int, value string) error {
	if v, err := g.SetView(name, x0, y0, x1, y1); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = name
		fmt.Fprintln(v, value)
	}
	return nil
}

func initKeybindings(g *gocui.Gui, portCh chan serial.Port, chDataSet chan map[byte]string) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlS, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			if !spMenuActive {
				return spSelectMenu(g, portCh)
			}
			return nil
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlF, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			if !freqMenuActive {
				return freqSetMenu(g, chDataSet)
			}
			return nil
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlM, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return modeSetMenu(chDataSet)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlA, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return attSetMenu(chDataSet)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlP, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return preampSetMenu(chDataSet)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return afMinusMenu(chDataSet)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlW, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return afPlusMenu(chDataSet)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlE, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return rfMinusMenu(chDataSet)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlR, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return rfPlusMenu(chDataSet)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlZ, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return sqlMinusMenu(chDataSet)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlX, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return sqlPlusMenu(chDataSet)
		}); err != nil {
		return err
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
