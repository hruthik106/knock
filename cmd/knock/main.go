package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

const timeout = 5 * time.Second

type FilterType int

const (
	FilterNone FilterType = iota
	FilterAlive
	FilterUnhealthy
	FilterUnreachable
)

type Config struct {
	Targets []string
	Timeout time.Duration
	Method  string
	only    FilterType
}

func main() {
	cfg, err := parseConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		printUsage()
		os.Exit(2)
	}
	if len(cfg.Targets) > 1 {
		fmt.Printf("knocking %d targets \n\n", len(cfg.Targets))
	}
	exitcode := 0

	for _, url := range cfg.Targets {
		status, latency, err := knock(url, cfg)
		result := classify(status, err)

		if !shouldPrint(result, cfg.only) {
			exitcode = max(exitcode, exitFromResult(result))
			continue
		}
		printResult(url, status, latency, result)
		exitcode = max(exitcode, exitFromResult(result))

	}
	os.Exit(exitcode)

}

type ResultType int

const (
	ResultAlive ResultType = iota
	ResultUnhealthy
	ResultUnreachable
)

func classify(status int, err error) ResultType {
	if err != nil {
		return ResultUnreachable
	}
	if status >= 200 && status < 300 {
		return ResultAlive
	}
	return ResultUnhealthy
}

func shouldPrint(r ResultType, f FilterType) bool {
	if f == FilterNone {
		return true
	}
	return int(r) == int(f-1)
}

func exitFromResult(r ResultType) int {
	switch r {
	case ResultAlive:
		return 0
	case ResultUnhealthy:
		return 1
	default:
		return 3
	}
}

func printResult(url string, status int, latency string, r ResultType) {

	switch r {
	case ResultAlive:
		fmt.Printf("✔ %s (%d, %s)\n", url, status, latency)
	case ResultUnhealthy:
		fmt.Printf("✖ %s (%d, %s)\n", url, status, latency)
	case ResultUnreachable:
		fmt.Printf("✖ %s (unreachable)\n", url)

	}
}

func parseConfig() (*Config, error) {
	cfg := &Config{
		Timeout: 5 * time.Second,
		Method:  http.MethodHead,
		only:    FilterNone,
	}

	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	file := fs.String("f", "", "")
	fs.StringVar(file, "file", "", "")

	timeout := fs.Duration("t", cfg.Timeout, "")
	fs.DurationVar(timeout, "timeout", cfg.Timeout, "")

	method := fs.String("method", cfg.Method, "")

	only := fs.String("o", "", "")
	fs.StringVar(only, "only", "", "")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	//validate method
	m := strings.ToUpper(*method)
	if m != http.MethodHead && m != http.MethodGet {
		return nil, fmt.Errorf("invalid method : ", *method)
	}

	cfg.Method = m
	cfg.Timeout = *timeout

	//parse filter
	if *only != "" {
		switch strings.ToLower(*only) {
		case "al", "alive":
			cfg.only = FilterAlive
		case "uh", "unhealthy":
			cfg.only = FilterUnhealthy
		case "ur", "unreachable":
			cfg.only = FilterUnreachable
		default:
			return nil, fmt.Errorf("invalid filter ", *only)
		}
	}

	//input source
	args := fs.Args()
	if *file != "" && len(args) > 0 {
		return nil, fmt.Errorf("use either a file or a url , not both")
	}
	if *file != "" {
		targets, err := readTargets(*file)
		if err != nil {
			return nil, err
		}
		cfg.Targets = targets
		return cfg, nil
	}
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid arguments ")
	}
	cfg.Targets = []string{args[0]}
	return cfg, nil
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

func knock(url string, cfg *Config) (int, string, error) {
	client := http.Client{
		Timeout: cfg.Timeout,
	}
	req, err := http.NewRequest(cfg.Method, url, nil)

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
	fmt.Fprintln(os.Stderr, " knock <url> [flags]")
	fmt.Fprintln(os.Stderr, " knock -f <file> [flags]")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "flags:")
	fmt.Fprintln(os.Stderr, "  -f , --file <path>      read targets from file")
	fmt.Fprintln(os.Stderr, "  -t , --timeout <dur>    request timeout (default 5s)")
	fmt.Fprintln(os.Stderr, "  --method <HEAD|GET>     http method (default HEAD)")
	fmt.Fprintln(os.Stderr, "  -o, --only <filter>     al|uh|ur or alive|unhealthy|unreachable")

}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
