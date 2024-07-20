package learn

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gucio32/morse/pkg/generator"
	"golang.org/x/exp/rand"
)

// Lesson represents a particular lesson (lessons map below)
// According to https://morsecode.world/international/timing.html
// it is recommended to work on PARIS=20 and increase InterChar and InterWord breaks
type Lesson struct {
	Letters     string
	InterChar   int
	InterWord   int
	Description string
}

func GetLesson(lessonIdx int) Lesson {
	m := map[int]Lesson{
		1: {"aelv", 9, 21, "First lesson: a, e, l, v higher inter-words and inter-characters breaks"},
		2: {"aelv", 6, 14, "Second lesson: a, e, l, v medium inter-words and inter-characters breaks"},
		3: {"aelvcqst", 9, 21, "Third lesson: a, e, l, v, c, q, s, t higher inter-words and inter-characters breaks"},
		4: {"aelvcqst", 6, 14, "Fourth lesson: a, e, l, v, c, q, s, t medium inter-words and inter-characters breaks"},
		5: {"aelvcqstfgny", 9, 21, "Fifth lesson: a, e, l, v, c, q, s, t, f, g, n, y higher inter-words and inter-characters breaks"},
		6: {"aelvcqstfgny", 6, 14, "Sixth lesson: a, e, l, v, c, q, s, t, f, g, n, y medium inter-words and inter-characters breaks"},
	}

	l, ok := m[lessonIdx]
	if !ok {
		panic("Lesson not found")
	}

	return l
}

func Tutorial(lesson Lesson) {
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
			ncontinue
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
}

func StartLesson(l Lesson, nWords int) {
	// 0 initialize random (unix timestamp)
	rand.Seed(uint64(time.Now().UnixNano()))

	// 1.0 generate nWords 5-letters words from l.Letters
	words := []string{}
	for i := 0; i < nWords; i++ {
		// 1.1 generate a random 5-letters word
		word := ""
		for j := 0; j < 5; j++ {
			word += string(l.Letters[rand.Intn(len(l.Letters))])
		}
		words = append(words, word)
	}

	// 2 generate morse code
	// 2.1 create generator
	g, err := generator.NewGenerator()
	if err != nil {
		panic(err)
	}

	// 2.2 setup generator
	g.SetCustomSeparator(generator.InterCharacter, l.InterChar).
		SetCustomSeparator(generator.InterWord, l.InterWord)

	// 2.3 in goroutine: generate morse code
	go func() {
		g.Play(strings.Join(words, " "))
	}()

	// 3 read user's answer
	// 3.1 prompt user
	fmt.Printf("What do ou hear?: ")
	reader := bufio.NewReader(os.Stdin)
	answer, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	answer = strings.TrimSuffix(answer, "\n")
	// 4 compare user's answer with the generated one (print correct letters in green and incorrect in red - if you donot use linux you have the problem)
	correct := 0
	text := strings.Join(words, " ")
	for i, t := range text {
		if i >= len(answer) { // gray
			fmt.Printf("\033[37m%c\033[0m", t)
		} else if byte(t) != answer[i] {
			fmt.Printf("\033[31m%c\033[0m", t)
		} else {
			fmt.Printf("\033[32m%c\033[0m", t)
			correct++
		}
	}

	fmt.Println()

	fmt.Printf("%d/%d correct\n", correct, len(text))
	if correct == len(text) {
		fmt.Println("Congratulations! You can go ahead!")
	}
}
