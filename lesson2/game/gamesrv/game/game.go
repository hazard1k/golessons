package game

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	Sum = iota
	Minus
	Dev
	Multiply
)

type Game struct {
	question   string
	operString string
	answer     int
}

func (g *Game) Start() {
	rand.Seed(time.Now().UnixNano())
	g.generateQuestion()
}

func (g *Game) randomOperation() int {
	return rand.Intn(4)
}

func (g *Game) CurrentQuestion() string {
	return g.question
}

func (g *Game) generateQuestion() string {

	min := 0
	max := 10
	l := rand.Intn(max-min+1) + min
	r := rand.Intn(max-min+1) + min

	switch g.randomOperation() {
	case Sum:
		g.answer = l + r
		g.operString = "+"
	case Minus:
		g.answer = l - r
		g.operString = "-"
	case Dev:
		g.answer = l / r
		g.operString = "/"
	case Multiply:
		g.answer = l * r
		g.operString = "*"
	}
	g.question = fmt.Sprintf("%d %s %d", l, g.operString, r)

	return g.question
}

func (g *Game) NextQuestion(ch chan<- string) {

	if g.currentQuestion() != "" {
		ch <- "Генерируем новый вопрос..."
		time.Sleep(3 * time.Second)
	}
	ch <- g.generateQuestion()

}

func (g *Game) currentQuestion() string {
	return g.question
}

func (g *Game) IsAnswerCorrect(userAnswer int) bool {
	if g.answer == userAnswer {
		return true
	}
	return false
}
