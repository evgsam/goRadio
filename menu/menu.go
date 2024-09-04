// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package menu

import (
	"errors"
	"fmt"
	datastruct "goRadio/dataStruct"
	component "goRadio/gocui-component"
	"log"

	"github.com/jroimartin/gocui"
	//"github.com/awesome-gocui/gocui"
)

func Menu(ch chan *datastruct.RadioSettings) {
	//var myRadioSettings radioSettings
	myRadioSettings := <-ch
	fmt.Println(myRadioSettings.Mode)
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}
}

/*
	func quit(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}
*/
func layout(g *gocui.Gui) error {
	//maxX, maxY := g.Size()

	// Overlap (front)

	if v, err := g.SetView("v0", 0, 0, 50, 7); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " IC-78 Information "
	}

	if v, err := g.SetView("v1", 1, 1, 16, 3); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " status "
		//v.BgColor = gocui.GetColor("#FFAA55")
		//v.TitleColor = gocui.GetColor("#FFAA55")
		//fmt.Fprint(v, "\n")
		fmt.Fprint(v, "Disconnected!")
	}

	if v, err := g.SetView("v2", 17, 1, 27, 3); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " mode "
		fmt.Fprint(v, "RTTY")
	}

	if v, err := g.SetView("v3", 28, 1, 38, 3); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " ATT "
		fmt.Fprint(v, "YES")
	}

	if v, err := g.SetView("v4", 39, 1, 49, 3); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " preamp "
		fmt.Fprint(v, "P.AMP")
	}

	if v, err := g.SetView("v5", 1, 4, 16, 6); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " freque "
		fmt.Fprint(v, "2999999 Hz")
	}

	if v, err := g.SetView("v6", 17, 4, 27, 6); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " AF "
		fmt.Fprint(v, "100 %")
	}

	if v, err := g.SetView("v7", 28, 4, 38, 6); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " RF "
		fmt.Fprint(v, "100 %")
	}

	if v, err := g.SetView("v8", 39, 4, 49, 6); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " SQL "
		fmt.Fprint(v, "100 %")
	}

	component.NewSelect(g, "Mode:", 2, 9, 0, 5).
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

	/*if v, err := g.SetView("v3", 0, 4, 30, 7, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = "Transiver mode"
	}
	*/
	return nil
}
