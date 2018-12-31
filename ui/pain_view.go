package ui

import (
	"github.com/jroimartin/gocui"
	"github.com/JajaDoc/g-explorer/objects"
	"fmt"
	"strings"
	"github.com/lunixbochs/vtclean"
)

// PainView is
type PainView struct {
	PainNo            int
	Name              string
	gui               *gocui.Gui
	view              *gocui.View
	header            *gocui.View
	Index             int
	Objects           []objects.Objects
	Path              string

	//keybindingCompareAll   []keybinding.Key
	//keybindingCompareLayer []keybinding.Key
}

// NewDetailsView creates a new view object attached the the global [gocui] screen object.
func NewPain1View(name string, gui *gocui.Gui, painNo int, path string, objects []objects.Objects) (pain1View *PainView) {
	pain1View = new(PainView)

	// populate main fields
	pain1View.PainNo = painNo
	pain1View.Name = name
	pain1View.gui = gui
	pain1View.Path = path
	pain1View.Index = 0
	pain1View.Objects = objects

	return pain1View
}

// Setup initializes the UI concerns within the context of a global [gocui] view object.
func (view *PainView) Setup(v *gocui.View, header *gocui.View) error {

	// set view options
	view.view = v
	view.view.Editable = false
	view.view.Wrap = false
	view.view.Frame = true

	view.header = header
	view.header.Editable = false
	view.header.Wrap = false
	view.header.Frame = false

	// set keybindings
	if err := view.gui.SetKeybinding(view.Name, gocui.KeyArrowDown, gocui.ModNone, func(*gocui.Gui, *gocui.View) error { return view.CursorDown() }); err != nil {
		return err
	}
	if err := view.gui.SetKeybinding(view.Name, gocui.KeyArrowUp, gocui.ModNone, func(*gocui.Gui, *gocui.View) error { return view.CursorUp() }); err != nil {
		return err
	}
	if err := view.gui.SetKeybinding(view.Name, gocui.KeyEnter, gocui.ModNone, func(*gocui.Gui, *gocui.View) error { return view.Enter() }); err != nil {
		return err
	}

	return view.Render()
}

// IsVisible indicates if the layer view pane is currently initialized.
func (view *PainView) IsVisible() bool {
	if view == nil {
		return false
	}
	return true
}

// CursorDown moves the cursor down in the layer pane (selecting a higher layer).
func (view *PainView) CursorDown() error {
	if view.Index < len(view.Objects) {
		err := CursorDown(view.gui, view.view)
		if err == nil {
			view.SetCursor(view.Index + 1)
			view.Index++
		}
	}
	return nil
}

// CursorUp moves the cursor up in the layer pane (selecting a lower layer).
func (view *PainView) CursorUp() error {
	if view.Index > 0 {
		err := CursorUp(view.gui, view.view)
		if err == nil {
			view.SetCursor(view.Index - 1)
			view.Index--
		}
	}
	return nil
}

// Enter input enter
func (view *PainView) Enter() error {
	// TODO: impl
	return nil
}

// SetCursor resets the cursor and orients the file tree view based on the given layer index.
func (view *PainView) SetCursor(layer int) error {
	view.Render()
	return nil
}

// Update refreshes the state objects for future rendering (currently does nothing).
func (view *PainView) Update() error {

	return nil
}

// Render flushes the state objects to the screen. The layers pane reports:
// 1. the layers of the image + metadata
// 2. the current selected image
func (view *PainView) Render() error {

	// indicate when selected
	title := view.Path
	if view.gui.CurrentView() == view.view {
		title = "● " + title
	}

	view.gui.Update(func(g *gocui.Gui) error {
		// update header
		view.header.Clear()
		width, _ := g.Size()
		headerStr := fmt.Sprintf("[%s]%s\n", title, strings.Repeat("─", width*2))
		headerStr += fmt.Sprintf("%s", "TODO")
		fmt.Fprintln(view.header, Formatting.Header(vtclean.Clean(headerStr, false)))

		// update contents
		view.view.Clear()
		for _, obj := range view.Objects {
			fmt.Fprintln(view.view, obj.Info.Name())
		}

		return nil
	})
	return nil
}

// KeyHelp indicates all the possible actions a user can take while the current pane is selected.
func (view *PainView) KeyHelp() string {
	return "TODO: Help!"
	//return renderStatusOption(view.keybindingCompareLayer[0].String(), "Show layer changes", view.CompareMode == CompareLayer) +
	//	renderStatusOption(view.keybindingCompareAll[0].String(), "Show aggregated changes", view.CompareMode == CompareAll)
}
