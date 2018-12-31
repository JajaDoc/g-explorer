package ui

import (
	"github.com/jroimartin/gocui"
	"log"
	"github.com/JajaDoc/g-explorer/utils"
	"github.com/pkg/errors"
	"github.com/JajaDoc/g-explorer/objects"
	"github.com/fatih/color"
)

const debug = false

var Formatting struct {
	Header                func(...interface{}) string
	Selected              func(...interface{}) string
	StatusSelected        func(...interface{}) string
	StatusNormal          func(...interface{}) string
	StatusControlSelected func(...interface{}) string
	StatusControlNormal   func(...interface{}) string
	CompareTop            func(...interface{}) string
	CompareBottom         func(...interface{}) string
}

// Views contains all rendered UI panes.
var Views struct {
	Pain1    *PainView
	Pain2    *PainView
	//Pain3  *StatusView
	Detail  *DetailView
	lookup  map[string]View
}

// View defines the a renderable terminal screen pane.
type View interface {
	Setup(*gocui.View, *gocui.View) error
	CursorDown() error
	CursorUp() error
	Enter() error
	Render() error
	Update() error
	KeyHelp() string
	IsVisible() bool
}

// CursorDown moves the cursor down in the currently selected gocui pane, scrolling the screen as needed.
func CursorDown(g *gocui.Gui, v *gocui.View) error {
	cx, cy := v.Cursor()

	// if there isn't a next line
	line, err := v.Line(cy + 1)
	if err != nil {
		// todo: handle error
	}
	if len(line) == 0 {
		return errors.New("unable to move cursor down, empty line")
	}
	if err := v.SetCursor(cx, cy+1); err != nil {
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}
	return nil
}

// CursorUp moves the cursor up in the currently selected gocui pane, scrolling the screen as needed.
func CursorUp(g *gocui.Gui, v *gocui.View) error {
	ox, oy := v.Origin()
	cx, cy := v.Cursor()
	if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
		if err := v.SetOrigin(ox, oy-1); err != nil {
			return err
		}
	}
	return nil
}

// isNewView determines if a view has already been created based on the set of errors given (a bit hokie)
func isNewView(errs ...error) bool {
	for _, err := range errs {
		if err == nil {
			return false
		}
		if err != nil && err != gocui.ErrUnknownView {
			return false
		}
	}
	return true
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	debugWidth := 0
	if debug {
		debugWidth = maxX / 4
	}
	debugCols := maxX - debugWidth

	bottomRows := 1
	headerRows := 1
	splitCols := maxX / 3

	var view, header *gocui.View
	var viewErr, headerErr, err error

	// Debug pane
	if debug {
		if _, err := g.SetView("debug", debugCols, -1, maxX, maxY-bottomRows); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
		}
	}

	view, viewErr = g.SetView(Views.Pain1.Name, 0, -1+headerRows, splitCols, maxY-1)
	header, headerErr = g.SetView(Views.Pain1.Name+"header", 0, -1, splitCols, headerRows)
	if isNewView(viewErr, headerErr) {
		Views.Pain1.Setup(view, header)

		if _, err = g.SetCurrentView(Views.Pain1.Name); err != nil {
			return err
		}
		// since we are selecting the view, we should rerender to indicate it is selected
		Views.Pain1.Render()
	}

	view, viewErr = g.SetView(Views.Pain2.Name, splitCols, -1+headerRows, splitCols * 2, maxY-1)
	header, headerErr = g.SetView(Views.Pain2.Name+"header", splitCols, -1, splitCols * 2, headerRows)
	if isNewView(viewErr, headerErr) {
		Views.Pain2.Setup(view, header)
	}

	view, viewErr = g.SetView(Views.Detail.Name, splitCols * 2, -1+headerRows, maxX-1, maxY-1)
	header, headerErr = g.SetView(Views.Detail.Name+"header", splitCols * 2, -1, maxX-1, headerRows)
	if isNewView(viewErr, headerErr) {
		Views.Detail.Setup(view, header)
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// Update refreshes the state objects for future rendering.
func Update() {
	for _, view := range Views.lookup {
		view.Update()
	}
}

// Render flushes the state objects to the screen.
func Render() {
	for _, view := range Views.lookup {
		if view.IsVisible() {
			view.Render()
		}
	}
}

// keyBindings registers global key press actions, valid when in any pane.
func keyBindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	return nil
}

func Run() {
	Formatting.Selected = color.New(color.ReverseVideo, color.Bold).SprintFunc()
	Formatting.Header = color.New(color.Bold).SprintFunc()
	Formatting.StatusSelected = color.New(color.BgMagenta, color.FgWhite).SprintFunc()
	Formatting.StatusNormal = color.New(color.ReverseVideo).SprintFunc()
	Formatting.StatusControlSelected = color.New(color.BgMagenta, color.FgWhite, color.Bold).SprintFunc()
	Formatting.StatusControlNormal = color.New(color.ReverseVideo, color.Bold).SprintFunc()
	Formatting.CompareTop = color.New(color.BgMagenta).SprintFunc()
	Formatting.CompareBottom = color.New(color.BgGreen).SprintFunc()

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	utils.SetUi(g)
	defer g.Close()

	path := "./"
	objs, err := objects.GetObjects(path)
	if err != nil {
		log.Panicln(err)
		utils.Exit(0)
	}

	Views.lookup = make(map[string]View)
	Views.Pain1 = NewPain1View("pain1", g,1, path, objs)
	Views.lookup[Views.Pain1.Name] = Views.Pain1
	Views.Pain2 = NewPain1View("pain2", g, 1, "", []objects.Objects{})
	Views.lookup[Views.Pain2.Name] = Views.Pain2
	Views.Detail = NewDetailView("detail", g, "", nil)
	Views.lookup[Views.Detail.Name] = Views.Detail

	g.Cursor = true
	g.Mouse  = true
	g.SetManagerFunc(layout)

	Update()
	Render()

	if err := keyBindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
	utils.Exit(0)
}
