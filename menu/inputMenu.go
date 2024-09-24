package menu

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jroimartin/gocui"
)

const helpText = `KEYBINDINGS
Tab: Move between buttons
Enter: Push button
^C: Exit`

type Label struct {
	name string
	x, y int
	w, h int
	body string
}

var (
	text_     string
	portsList = make([]string, 10)
)

type PortWidget struct {
	name string
	x, y int
	w, h int
	body string
}

func NewPortInfoWidget(name string, x, y int, body string) *PortWidget {
	lines := strings.Split(body, "\n")
	w := 0
	for _, l := range lines {
		if len(l) > w {
			w = len(l)
		}
	}
	h := len(lines) + 1
	w = w + 1
	return &PortWidget{name: name, x: x, y: y, w: w, h: h, body: body}
}

func (w *PortWidget) Layout(g *gocui.Gui) error {
	v, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+w.h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprint(v, w.body)
	}
	return nil
}

func NewLabel(name string, x, y int, body string) *Label {
	lines := strings.Split(body, "\n")
	w := 0
	for _, l := range lines {
		if len(l) > w {
			w = len(l)
		}
	}
	h := len(lines) + 1
	w = w + 1

	return &Label{name: name, x: x, y: y, w: w, h: h, body: body}
}

func (l *Label) Layout(g *gocui.Gui) error {
	v, err := g.SetView(l.name, l.x, l.y, l.x+l.w, l.y+l.h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		fmt.Fprint(v, l.body)
	}
	return nil
}

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

func inputMenu( /*ports []string,*/ ch chan string, inputCh chan *inputFormsStruct) {
	p := <-inputCh
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.Cursor = true

	if p.flag {
		t := ""
		for i, s := range p.ports {
			t += "Port №" + strconv.Itoa(i) + ":" + s + "\n"
		}
		t += "Please select port number:"
		portsList = p.ports
		portInfo := NewPortInfoWidget("ports", 1, 1, t)
		input := NewInput("input", 1, len(p.ports)+3, 27, 2)
		focus := gocui.ManagerFunc(SetFocus("input"))
		g.SetManager(portInfo, input, focus)
	}

	if err := initKeybindings(g, ch); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func initKeybindings(g *gocui.Gui, ch chan string) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return gocui.ErrQuit
		}); err != nil {
		return err
	}

	/*
		err = g.SetKeybinding("", '1', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			_, err := g.SetViewOnTop("v1")
			return err
		})
		if err != nil {
			return err
		}
	*/

	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		/*if err := g.DeleteView("ports"); err != nil {
			if err != gocui.ErrUnknownView {
				panic(err)
			}
		}
		if err := g.DeleteView("input"); err != nil {
			if err != gocui.ErrUnknownView {
				panic(err)
			}
		}
		g.Close()
		*/
		_, err := g.SetViewOnBottom("ports")
		_, err = g.SetViewOnBottom("input")
		i, _ := strconv.Atoi(text_)
		ch <- portsList[i]
		return err
	}); err != nil {
		return err
	}

	return nil
}
