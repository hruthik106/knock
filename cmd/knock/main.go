package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		printUsage()
		os.Exit(2)
	}

	url := os.Args[1]
	fmt.Println("knocking on : ", url)

}

func printUsage() {
	fmt.Println(os.Stderr, "usage : knock <url>")
}
