package main

import (
	"os"
	"strconv"

	"github.com/gucio32/morse/pkg/learn"
)

func main() {
	idx, _ := strconv.Atoi(os.Args[1])
	learn.StartLesson(learn.GetLesson(idx), 3)
}
