package main

import (
	"conv"
	"conv/converter"
	"conv/input"
	"fmt"
	"log"
)

func main() {
	conv.Parse()
	fmt.Printf("config: %#v\n", conv.GetConf())
	src := input.ReadSeq()
	fmt.Printf("src: %#v\n", src)
	c := converter.Get("mulit_line_escape")
	b, err := c.From(src)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v", b)
}
