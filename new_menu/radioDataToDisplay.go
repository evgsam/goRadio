package newmenu

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

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
