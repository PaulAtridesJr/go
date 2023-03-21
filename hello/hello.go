package main

import (
    "fmt"
    "log"
    "example.com/greetings"
)

func main() {
// Set properties of the predefined Logger, including
    // the log entry prefix and a flag to disable printing
    // the time, source file, and line number.
    log.SetPrefix("greetings: ")
    log.SetFlags(log.LstdFlags)

var a [4]int
a[0] = 1
i := a[0]
log.Print(i)

//letters := []string{"a", "b", "c", "d"}
//s := make([]byte, 5)

    // Request a greeting message.
    message, err := greetings.Hello("hh")
    // If an error was returned, print it to the console and
    // exit the program.
    if err != nil {
        log.Fatal(err)
    }

    // If no error was returned, print the returned message
    // to the console.
    fmt.Println(message)

	 // A slice of names.
	 names := []string{"Gladys", "Samantha", "Darrin"}

	 // Request greeting messages for the names.
	 messages, err := greetings.Hellos(names)
	 if err != nil {
		 log.Fatal(err)
	 }
	 // If no error was returned, print the returned map of
	 // messages to the console.
	 fmt.Println(messages)
}