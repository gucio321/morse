package main

import (
	"flag"

	"github.com/gucio32/morse/pkg/learn"
)

func main() {
	tutor := flag.Bool("tutor", false, "Start tutorial mode")
	lessonIdx := flag.Int("lesson", 0, "Lesson index")
	words := flag.Int("words", 3, "Number of words to learn")
	flag.Parse()

	lesson := learn.GetLesson(*lessonIdx)

	if *tutor {
		learn.Tutorial(lesson)
		return
	}

	learn.StartLesson(lesson, *words)
}
