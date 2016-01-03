package preview

import "flag"

// Color is enum of color set.
type Color int

const (
	// Color8 is config to use 8 color.
	Color8 = Color(0)

	// Color256 is config to use 256 color.
	Color256 = Color(1)
)

// Config represents configuration data.
type Config struct {
	Width  uint
	Height uint
	Color  Color
}

// NewConfig parse args and create Config.
func NewConfig(args []string) (*Config, []string, error) {
	flags := flag.NewFlagSet("preview", flag.ContinueOnError)
	width := flags.Uint("width", 0, "max width (default terminal width)")
	height := flags.Uint("height", 128, "max height")
	color := flags.Bool("c8", false, "use only 8 color (default 256 color)")

	if err := flags.Parse(args); err != nil {
		return nil, nil, err
	}

	if *width == 0 {
		ws, err := GetWindowSize()
		if err != nil {
			return nil, nil, err
		}
		*width = uint(ws.Width)
	}

	var c Color
	if *color {
		c = Color8
	} else {
		c = Color256
	}

	conf := &Config{
		*width,
		*height,
		c,
	}

	return conf, flags.Args(), nil
}
