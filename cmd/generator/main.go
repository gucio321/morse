package main

import (
	"time"

	"github.com/gucio32/morse/pkg/generator"
)

func main() {
	g, err := generator.NewGenerator()
	if err != nil {
		panic(err)
	}

	_ = g
	g.SetPARIS(20)
	g.Play("qc qc qc vvv")
	time.Sleep(1 * time.Second)
}
