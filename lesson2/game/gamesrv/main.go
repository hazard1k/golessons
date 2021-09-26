package main

import (
	"bufio"
	"fmt"
	"golessons/lesson2/game/gamesrv/game"
	"log"
	"net"
	"strconv"
	"time"
)

type client chan<- string

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func main() {
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()

	g := &game.Game{}
	g.Start()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn, g)
	}
}

func handleConn(conn net.Conn, g *game.Game) {
	ch := make(chan string)
	input := bufio.NewScanner(conn)

	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are connected from " + who
	ch <- "Введите свое имя:"

	if input.Scan() {
		who = input.Text()
	}

	messages <- who + " has arrived"
	entering <- ch
	messages <- g.CurrentQuestion()

	log.Println(who + " has arrived")

	for input.Scan() {
		in := input.Text()
		messages <- who + ": " + in

		if a, err := strconv.Atoi(in); err == nil && g.IsAnswerCorrect(a) {
			messages <- "правильный ответ: " + g.Answer()
			messages <- who + " is a winner!"
			ch <- "Генерируем новый вопрос..."
			time.Sleep(3 * time.Second)
			messages <- g.NextQuestion()
		}
	}

	leaving <- ch
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func broadcaster() {
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}
