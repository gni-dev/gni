package main

import (
	"log"
	"os"

	"gni.dev/gni/cmd"
)

func main() {
	log.SetFlags(0)

	if len(os.Args) < 2 {
		log.Fatal("No command-line arguments specified.")
	}

	switch os.Args[1] {
	case "gen":
		cmd.Gen(os.Args[2:])
	default:
		log.Fatalf("'%s' is not valid command.\n", os.Args[1])
	}
}
