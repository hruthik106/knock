package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"time"
)

const timeout = 5 * time.Second

func main() {
	targets, err := parseArgs(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		printUsage()
		os.Exit(2)
	}

	if len(targets) > 1 {
		fmt.Printf("knocking %d targets \n\n", len(targets))
	}
	exitcode := 0

	for _, url := range targets {
		status, latency, err := knock(url)
		if err != nil {
			fmt.Printf("✖ %s (unreachable)\n", url)
			exitcode = max(exitcode, 3)
			continue
		}
		if status >= 200 && status < 300 {
			fmt.Printf("✔ %s (%d, %s)\n", url, status, latency)
			continue
		}
		fmt.Printf("✖ %s (%d, %s)\n", url, status, latency)
		exitcode = max(exitcode, 1)

	}
	os.Exit(exitcode)

}

func parseArgs(args []string) ([]string, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("invalid arguments ")
	}
	if args[1] == "-f" {
		if len(args) != 3 {
			return nil, fmt.Errorf("missing file")
		}
		return readTargets(args[2])
	}

	if len(args) == 2 {
		return []string{args[1]}, nil
	}

	return nil, fmt.Errorf("invalid argumenst")

}

func readTargets(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		urls = append(urls, line)
	}
	if len(urls) == 0 {
		return nil, fmt.Errorf("no target found")
	}
	return urls, scanner.Err()
}

func knock(url string) (int, string, error) {
	client := http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest(http.MethodHead, url, nil)

	if err != nil {
		return 0, "", err
	}
	start := time.Now()
	resp, err := client.Do(req)
	latency := time.Since(start)
	if err != nil {
		return 0, "", err
	}

	defer resp.Body.Close()

	return resp.StatusCode, formatLatency(latency), nil
}

func formatLatency(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.1fs", d.Seconds())
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "usage : ")
	fmt.Fprintln(os.Stderr, " knock <url>")
	fmt.Fprintln(os.Stderr, " knock -f <file>")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
