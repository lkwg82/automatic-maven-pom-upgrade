package main

import (
	"os"
	"log"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)
}

func main() {
	log.Print("test")
	log.Print("test")
	log.Print("test")
	log.Print("test")
	log.Print("test")
}