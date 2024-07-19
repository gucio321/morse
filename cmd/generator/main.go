package main

import (
	"time"

	"github.com/gucio32/morse/pkg/generator"
)

func main() {
	g, err := generator.NewGenerator()
	g.
		SetCustomSeparator(generator.InterCharacter, 6).
		SetCustomSeparator(generator.InterWord, 14)
	if err != nil {
		panic(err)
	}

	g.SetPARIS(20)
	g.Play("qc qc qc vvv")
	time.Sleep(1 * time.Second)
}
