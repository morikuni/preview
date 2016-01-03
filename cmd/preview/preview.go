package main

import (
	"fmt"
	"os"

	"github.com/morikuni/preview"
	_ "github.com/morikuni/preview/image"
	_ "github.com/morikuni/preview/text"
)

func main() {
	conf, args, err := preview.NewConfig(os.Args[1:])

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	for _, fname := range args {
		file, err := os.Open(fname)
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, fname)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		defer file.Close()

		r, err := preview.NewRenderer(file, conf)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		r.Render(os.Stdout)
	}
}
