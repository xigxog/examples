package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
)

var (
	who  string
	addr string
)

func main() {
	who = os.Getenv("HELLO_WORLD_WHO")
	if who == "" {
		who = "World"
	}

	flag.StringVar(&addr, "addr", "127.0.0.1:3333", "address http server should bind to")
	flag.Parse()

	fmt.Printf("starting http server on '%s'...\n", addr)
	http.HandleFunc("/hello", getHello)
	err := http.ListenAndServe(addr, nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("server closed")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func getHello(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("👋 Hello %s!", who)
	fmt.Println(msg)

	w.Write([]byte(msg))
}
