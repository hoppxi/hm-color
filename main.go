package main

import (
	"log"

	"github.com/hoppxi/recolor/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}