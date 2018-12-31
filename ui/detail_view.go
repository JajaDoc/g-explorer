package ui

import (
	"github.com/jroimartin/gocui"
	"github.com/JajaDoc/g-explorer/objects"
	"fmt"
	"strings"
	"github.com/lunixbochs/vtclean"
	"os"
	"log"
	"bufio"
	"path/filepath"
	"path"
)

// DetailView is
type DetailView struct {
	Name              string
	gui               *gocui.Gui
	view              *gocui.View
	header            *gocui.View
	Object            *objects.Object
	Path              string
}

var (
	previewFileSize = 10
	previewInfoSize = 6
)

// NewDetailsView creates a new view object attached the the global [gocui] screen object.
func NewDetailView(name string, gui *gocui.Gui, path string, object *objects.Object) (detailView *DetailView) {
	detailView = new(DetailView)

	// populate main fields
	detailView.Name = name
	detailView.gui = gui
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

// Enter input enter
func (view *DetailView) Enter() error {
	// TODO: impl
	return nil
}

// SetCursor resets the cursor and orients the file tree view based on the given layer index.
func (view *DetailView) SetCursor(layer int) error {
	view.Render()
	return nil
}

// Update refreshes the state objects for future rendering (currently does nothing).
func (view *DetailView) Update() error {
	return nil
}

// Render flushes the state objects to the screen.
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
		fmt.Fprintln(view.header, Formatting.Header(vtclean.Clean(headerStr, false)))

		// update contents
		view.view.Clear()

		_, maxY := g.Size()
		previewFileSize = maxY - 2 - previewInfoSize

		if len(view.Path) != 0 && view.Object != nil {
			if view.Object.Info.IsDir() {
				// TODO: impl enter dir

				// contents
				fmt.Fprintln(view.view, Formatting.Selected(vtclean.Clean("Preview", false)))
				err := view.previewContentsInDir(view.Object.Info.Name())
				if err != nil {
					log.Panicln(err)
					return err
				}

				// info
				fmt.Fprintln(view.view, Formatting.Selected(vtclean.Clean("\nInfo", false)))
				err = view.previewInfo()
				if err != nil {
					log.Panicln(err)
					return err
				}
			} else {
				// contents
				fmt.Fprintln(view.view, Formatting.Selected(vtclean.Clean("Preview", false)))
				err := view.previewFile()
				if err != nil {
					log.Panicln(err)
					return err
				}

				// info
				fmt.Fprintln(view.view, Formatting.Selected(vtclean.Clean("\nInfo", false)))
				err = view.previewInfo()
				if err != nil {
					log.Panicln(err)
					return err
				}
			}
		}

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

func (view *DetailView) selectObject(path string, object *objects.Object) {
	view.Path = path
	view.Object = object
	view.Render()
}

func (view *DetailView) previewFile() error {
	absPath, err := filepath.Abs(path.Join(view.Path, view.Object.Info.Name()))
	if err != nil {
		return err
	}

	fp, err := os.Open(absPath)
	if err != nil {
		return err
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for i := 0; i < previewFileSize && scanner.Scan(); i++ {
		fmt.Fprintln(view.view, scanner.Text())
	}
	return scanner.Err()
}

func (view *DetailView) previewInfo() error {
	t := `name:%s
mode:%s
size:%d
modTime:%s
`

	fmt.Fprintf(view.view, t,
		view.Object.Info.Name(),
		view.Object.Info.Mode().String(),
		view.Object.Info.Size(),
		view.Object.Info.ModTime())
	return nil
}

func (view *DetailView) previewContentsInDir(dir string) error {
	objectList, err := objects.GetObjects(path.Join(view.Path, dir))
	if err != nil {
		return err
	}

	printFormatting(view.view, &objectList)
	return nil
}

func (view *DetailView) enterDir(dir string) error {
	_, err := objects.ChangeDir(path.Join(view.Path, dir))
	if err != nil {
		return err
	}

	return nil
}