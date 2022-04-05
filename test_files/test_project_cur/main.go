package main

import (
	"fmt"
	"io"
	"os"

	"github.com/mattn/go-isatty"
)

func main() {
	var out io.Writer
	var isTerm bool
	if w, ok := out.(*os.File); !ok || os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd())) {
		isTerm = false
	}

	fmt.Println(isTerm)
}
