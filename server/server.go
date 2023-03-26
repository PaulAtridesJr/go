package main

import (
	"fmt"
	"net/http"
	"dummy"
	"os"
	"errors"
	"flag"
	"advanced"
)

func main() {	
	// 10.2.7.24:9000
	// :9000
	// 192.168.4.91:9000
	serverIP := flag.String("a", ":9000", "server IP")
	var mode int
	flag.IntVar(&mode, "m", 0, "0 - dummy (default), 1 - full")
	flag.Parse()

	fmt.Printf("Server IP: %s\n", *serverIP)
	var h func(w http.ResponseWriter, r *http.Request)

	switch mode {
		case 0:
			fmt.Printf("Server mode: dummy\n")
			h = dummy.DummyServe
		case 1:
			fmt.Printf("Server mode: full\n")
			h = advanced.AdvancedServe(true)
		break
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", h)

	err := http.ListenAndServe(*serverIP, mux)
  if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}