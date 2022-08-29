package conv

import (
	"flag"
)

func Parse() {
	flag.StringVar(&conf.InputType, "InputType", "plain", "")
	flag.StringVar(&conf.OutputType, "OutputType", "", "")

	flag.Parse()
}
