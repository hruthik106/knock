package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

const timeout = 5 * time.Second

func main() {
	url, err := parseArgs(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		printUsage()
		os.Exit(2)
	}

	fmt.Println("knock <", url, ">")

	status, err := knock(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, "✖ unreachable")
		os.Exit(3)
	}

	if status >= 200 && status < 300 {
		fmt.Println("✔ alive", fmt.Sprintf("%d", status))
		os.Exit(0)
	}

	fmt.Println("✖ unhealthy", fmt.Sprintf("(%d)", status))
	os.Exit(1)

}

func parseArgs(args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("invalid arguments ")
	}
	return args[1], nil
}

func knock(url string) (int, error) {
	client := http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest(http.MethodHead, url, nil)

	if err != nil {
		return 0, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	return resp.StatusCode, nil
}

func printUsage() {
	fmt.Println(os.Stderr, "Usage: knock <url>")
}
