package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var csvFile string
var score int
var shuffle bool
var timer int

// QuestionAnswer is a question/answer pair for quiz questions
type QuestionAnswer struct {
	Question, Answer string
}

// Determine flag values for quiz
func init() {
	flag.StringVar(&csvFile, "f", "problems.csv", "name of csv file for quiz problems")
	flag.IntVar(&timer, "t", 30, "duration of quiz timer")
	flag.BoolVar(&shuffle, "s", true, "t/f for shuffling the questions")
	flag.Parse()
}

// ParseCSV opens and parses CSV file for quiz
func ParseCSV(f string) []QuestionAnswer {
	file, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var quiz []QuestionAnswer
	csvReader := csv.NewReader(bufio.NewReader(file))
	for {
		line, error := csvReader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		quiz = append(quiz, QuestionAnswer{
			Question: line[0],
			Answer:   line[1],
		})
	}

	if shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(quiz), func(i, j int) { quiz[i], quiz[j] = quiz[j], quiz[i] })
	}

	return quiz
}

// Play runs the Q&A functionality of the quiz
func Play(quiz []QuestionAnswer) {
	quizTimer := time.NewTimer(time.Duration(timer) * time.Second)
	score = 0
game:
	for i, q := range quiz {
		answer := make(chan string)
		go func() {
			input := ""
			fmt.Printf("Problem #%d: %s = ", i+1, q.Question)
			fmt.Scanln(&input)
			answer <- strings.TrimSpace(input)
		}()

		select {
		case a := <-answer:
			if a == q.Answer {
				score++
			}
		case <-quizTimer.C:
			break game
		}
	}
	return
}

func main() {
	quiz := ParseCSV(csvFile)
	Play(quiz)
	fmt.Printf("\nYou got %d out of %d.\n", score, len(csvFile)+1)
}
