package main

import (
	"fmt"
	"os"
)

func main() {
	var status int
	if err := _main(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		status = 1
	}
	os.Exit(status)
}

func _main() error {
	text, image := render_statusline()

        if len(os.Args) < 2 {
          image()
        }

	if os.Args[1] == "text" {
		fmt.Print(text)
	} else {
		image()
	}

	return nil
}
