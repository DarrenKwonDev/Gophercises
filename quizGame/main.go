package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	// ./quiz -csv=test.csv 꼴로 flag를 붙인 값을 받아올 수 있음
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "time limit")

	flag.Parse() // Must be called after all flags are defined and before flags are accessed by the program.

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the csv file: %s\n", *csvFilename))
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll() // csv 크기가 작으므로 한번에 다 읽어버리자
	if err != nil {
		exit("Failed to parse the provided csv file")
	}

	problems := parseLines(lines)
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correct := 0

	fmt.Printf("start game. time limit is %d \n", *timeLimit)

problemLoop:
	for i, p := range problems {
		fmt.Printf("problem #%d: %s = \n", i+1, p.q)
		answerCh := make(chan string)

		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
			case <- timer.C:
				fmt.Printf("timeout. You scored %d\n", correct)
				break problemLoop
			case answer := <- answerCh:
				if answer != p.a {
					fmt.Printf("You're wrong! The correct answer is %s\n", p.a)
					fmt.Printf("You scored %d\n", correct)
					break problemLoop
				}
				correct++
				if correct == len(problems) {
					fmt.Println("well done. you got it all right")
					break problemLoop
				}
		}
	}
}

func parseLines(lines [][]string) []problem {
	// 왜 빈 슬라이스 만들고 그냥 append하지 않는가?
	// -> 우리는 이미 len을 알 수 있으므로 cap을 늘리고 하는 부가 동작을 막고 빠르게 사용하기 위해서임.
	problems := make([]problem, len(lines))

	for i, line := range lines {
		problems[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}

	return problems
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1) // something wrong
}