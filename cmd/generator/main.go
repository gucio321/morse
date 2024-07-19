package main

import "github.com/gucio32/morse/pkg/generator"

func main() {
	g, err := generator.NewGenerator()
	if err != nil {
		panic(err)
	}

	_ = g
	g.Play()
}
