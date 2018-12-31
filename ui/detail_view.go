package ui

import (
	"github.com/jroimartin/gocui"
	"github.com/JajaDoc/g-explorer/objects"
	"fmt"
	"strings"
	"github.com/lunixbochs/vtclean"
)

// DetailView is
type DetailView struct {
	Name              string
	gui               *gocui.Gui
	view              *gocui.View
	header            *gocui.View
	Index             int
	Object            *objects.Objects

}

// NewDetailsView creates a new view object attached the the global [gocui] screen object.
func NewDetailView(name string, gui *gocui.Gui, path string, object *objects.Objects) (detailView *DetailView) {
	detailView = new(DetailView)

	// populate main fields
	detailView.Name = name
	detailView.gui = gui
	detailView.Index = 0
	detailView.Object = object

	return detailView
}

// Setup initializes the UI concerns within the context of a global [gocui] view object.
func (view *DetailView) Setup(v *gocui.View, header *gocui.View) error {

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

	return view.Render()
}

// IsVisible indicates if the layer view pane is currently initialized.
func (view *DetailView) IsVisible() bool {
	if view == nil {
		return false
	}
	return true
}

// CursorDown moves the cursor down in the layer pane (selecting a higher layer).
func (view *DetailView) CursorDown() error {
	return nil
}

// CursorUp moves the cursor up in the layer pane (selecting a lower layer).
func (view *DetailView) CursorUp() error {
	return nil
}

// SetCursor resets the cursor and orients the file tree view based on the given layer index.
func (view *DetailView) SetCursor(layer int) error {
	//Views.Tree.setTreeByLayer(view.getCompareIndexes())
	//Views.Details.Render()
	view.Render()

	return nil
}

// Update refreshes the state objects for future rendering (currently does nothing).
func (view *DetailView) Update() error {

	return nil
}

// Render flushes the state objects to the screen. The layers pane reports:
// 1. the layers of the image + metadata
// 2. the current selected image
func (view *DetailView) Render() error {

	// indicate when selected
	title := ""
	if view.Object != nil {
		title = view.Object.Info.Name()
	}

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


		return nil
	})
	return nil
}

// KeyHelp indicates all the possible actions a user can take while the current pane is selected.
func (view *DetailView) KeyHelp() string {
	return "TODO: Help!"
	//return renderStatusOption(view.keybindingCompareLayer[0].String(), "Show layer changes", view.CompareMode == CompareLayer) +
	//	renderStatusOption(view.keybindingCompareAll[0].String(), "Show aggregated changes", view.CompareMode == CompareAll)
}
