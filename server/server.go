package main

import (
	"fmt"
	"net/http"
	"dummy"
	"os"
	"errors"
)

func main() {	
	mux := http.NewServeMux()

	mux.HandleFunc("/", dummy.DummyServe)

	err := http.ListenAndServe(":3333", mux)
  if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}