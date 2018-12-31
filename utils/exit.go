package utils

import (
	"github.com/jroimartin/gocui"
	"fmt"
	"os"
)

var ui *gocui.Gui

func SetUi(g *gocui.Gui) {
	ui = g
}

func PrintAndExit(args ...interface{}) {
	//logrus.Println(args...)
	Cleanup()
	fmt.Println(args...)
	os.Exit(1)
}

// Note: this should only be used when exiting from non-gocui code
func Exit(rc int) {
	Cleanup()
	os.Exit(rc)
}

func Cleanup() {
	if ui != nil {
		ui.Close()
	}
	//ansi.CursorShow()
}
