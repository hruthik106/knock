package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		printUsage()
		os.Exit(2)
	}

	url := os.Args[1]
	fmt.Println("knocking on : ", url)

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(os.Stderr, "unreachable url : ", url)
		os.Exit(3)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Println("successfully knocked on : ", url)
		fmt.Println("responded with : ", resp.Status)
		os.Exit(0)
	}

	fmt.Println("X ", resp.Status)
	os.Exit(1)

}

func printUsage() {
	fmt.Println(os.Stderr, "Usage: knock <url>")
}
