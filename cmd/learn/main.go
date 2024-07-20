package main

import (
	"flag"
	"fmt"

	"github.com/gucio32/morse/pkg/generator"
	"github.com/gucio32/morse/pkg/learn"
)

func main() {
	tutor := flag.Bool("tutor", false, "Start tutorial mode")
	lessonIdx := flag.Int("lesson", 0, "Lesson index")
	words := flag.Int("words", 3, "Number of words to learn")
	flag.Parse()

	lesson := learn.GetLesson(*lessonIdx)

	if *tutor {
		// print letters in this lesson and their morse code
		for _, letter := range lesson.Letters {
			code, valid := generator.TranslateMorse(rune(letter))
			if !valid {
				panic("Not a valid letter")
			}

			fmt.Printf("%s: %s\n", string(rune(letter)), code)
		}

		var answer string
		fmt.Print("Type letter to hear it or enter to exit: ")
		gen, err := generator.NewGenerator()
		if err != nil {
			panic(err)
		}

		for fmt.Scanln(&answer); answer != ""; fmt.Scanln(&answer) {
			// check if onnly 1 letters
			if len(answer) != 1 {
				fmt.Println("Only one letter allowed")
				continue
			}

			// check if letter is in lesson
			c := false
			for _, letter := range lesson.Letters {
				if answer[0] == byte(letter) {
					c = true
					break
				}
			}

			if !c {
				fmt.Println("Not in lesson")
				continue
			}

			gen.Play(string(rune(answer[0])))

			fmt.Print("Type letter to hear it or enter to exit: ")
		}

		return
	}

	learn.StartLesson(lesson, *words)
}
