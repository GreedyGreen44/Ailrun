package main

import (
	"log"
	"os"
)

var (
	mainLog  *log.Logger
	errorLog *log.Logger
)

func main() {

	mainLog = log.New(os.Stdout, "Main: ", log.LstdFlags|log.Lmicroseconds)
	errorLog = log.New(os.Stderr, "Error: ", log.LstdFlags|log.Lmicroseconds)

	mainLog.Println("Starting program")

	clArgs := os.Args[1:]

	if len(clArgs) != 1 {
		mainLog.Println("Error! Incorrect number of arguments. You only need to provide path to config file")
		return
	}
}
