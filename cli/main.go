package main

import (
	_ "conv"
	"conv/converter"
	"conv/input"
	"fmt"
	"log"
)

func main() {
	src := input.ReadSeq()
	c := converter.Get("mulit_line_escape")
	b, err := c.From(src)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v", b)
}
