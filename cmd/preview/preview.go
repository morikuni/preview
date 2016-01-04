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

	for _, path := range args {
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, path)
		if err := preview.Preview(path, os.Stdout, conf); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
