package main

import (
	"log"

	"github.com/jessebrands/triforceblitz/python"
	"github.com/jessebrands/triforceblitz/randomizer"
)

func main() {
	interpreter, err := python.Find()
	if err != nil {
		log.Fatalln("Could not find Python interpreter:", err)
	}
	version, err := interpreter.Version()
	if err != nil {
		log.Fatalln("Could not get Python version:", err)
	}
	log.Println("Found Python version", version, "at", interpreter.Path())
	svc := randomizer.NewService(randomizer.DefaultConfig())
	log.Println(svc)
}

