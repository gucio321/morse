package learn

import (
	"fmt"

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

func GetLessons() map[int]Lesson {
	return map[int]Lesson{
		1: {{'q', 'v', 'e', 'a'}, 9, 21},
		2: {{'q', 'v', 'e', 'a'}, 6, 14},
	}
}

func StartLesson(l Lesson, nWords int) {
	// 1.0 generate nWords 5-letters words from l.Letters
	words := []string{}
	for i := 0; i < nWords; i++ {
		// 1.1 generate a random 5-letters word
		word := ""
		for j := 0; j < 5; j++ {
			word += l.Letters[rand.Intn(len(l.Letters))]
		}
		words = append(words, word)
	}

	fmt.Println(words)
}
