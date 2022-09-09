package conv

import (
	"flag"
)

func Parse() {
	flag.StringVar(&conf.InputType, "i", "plain", "input type")
	flag.StringVar(&conf.OutputType, "o", "", "output type")

	flag.Parse()
}
