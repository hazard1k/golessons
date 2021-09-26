package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"
)

type client chan<- string

var (
	entering = make(chan client)
	leaving  = make(chan client)
	srvCh    = make(chan string)
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	cfg := net.ListenConfig{
		KeepAlive: time.Minute,
	}
	l, err := cfg.Listen(ctx, "tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()

	go func() {
		input := bufio.NewScanner(os.Stdin)
		for input.Scan() {
			srvCh <- input.Text()
		}
	}()

	wg := &sync.WaitGroup{}
	log.Println("im started!")

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			conn, err := l.Accept()
			if err != nil {
				log.Println(err)
			} else {
				wg.Add(1)
				go handleConn(ctx, conn, wg)
			}
		}
	}()

	<-ctx.Done()

	log.Println("done")
	l.Close()
	wg.Wait()
	log.Println("exit")
}

func handleConn(ctx context.Context, conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()
	ch := make(chan string)
	entering <- ch

	go clientWriter(conn, ch)
	// каждую 1 секунду отправлять клиентам текущее время сервера
	tck := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-tck.C:
			fmt.Fprintf(conn, "now: %s\n", t)
		}
	}
	leaving <- ch
	close(ch)
}

func broadcaster() {
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-srvCh:
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
		}
	}
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}
