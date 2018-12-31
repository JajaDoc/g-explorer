package ui

import (
	"github.com/JajaDoc/g-explorer/objects"
	"fmt"
	"github.com/lunixbochs/vtclean"
	"io"
)

var Formatting struct {
	Header                func(...interface{}) string
	Selected              func(...interface{}) string
	StatusSelected        func(...interface{}) string
	StatusNormal          func(...interface{}) string
	StatusControlSelected func(...interface{}) string
	StatusControlNormal   func(...interface{}) string
	CompareTop            func(...interface{}) string
	CompareBottom         func(...interface{}) string
	Directory             func(...interface{}) string
}

func printFormatting(w io.Writer, objectList *[]objects.Object) {
	for _, obj := range *objectList {
		if obj.Info.IsDir() {
			fmt.Fprintln(w, Formatting.Directory(vtclean.Clean(obj.Info.Name(), true)))
		} else {
			fmt.Fprintln(w, vtclean.Clean(obj.Info.Name(), false))
		}
	}
}
