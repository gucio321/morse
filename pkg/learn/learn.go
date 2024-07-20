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
	Letters   []rune
	InterChar int
	InterWord int
}

func GetLesson(lessonIdx int) Lesson {
	m := map[int]Lesson{
		1: {[]rune{'a', 'e', 'l', 'v'}, 9, 21},
		2: {[]rune{'a', 'e', 'l', 'v'}, 6, 14},
		3: {[]rune{'a', 'e', 'l', 'v', 'c', 'q', 's', 't'}, 9, 21},
		4: {[]rune{'a', 'e', 'l', 'v', 'c', 'q', 's', 't'}, 6, 14},
	}

	l, ok := m[lessonIdx]
	if !ok {
		panic("Lesson not found")
	}

	return l
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
