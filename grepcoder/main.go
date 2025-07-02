package main

import (
	"fmt"
	"os"
	"sync"
	"net/http"
	"time"
	"strings"
	"strconv"
	"golang.org/x/net/html"
	"io"
)

// Help
func PrintHelp() {
	fmt.Println("Usage: \n\tgo run main.go <startId> <endId> <comma-separated keywords>")
}

// Counter for random User-Agent
var (
	counter     int
	userAgents  = []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Safari/605.1.15",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Linux; Android 10; SM-A505FN) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36",
	}
	problemIdABC = []string { "c", "d", "e", "f", "g" }		// Add Problem a, b, c, g
	problemIdARC = []string { "a", "b", "c" }
	prefixABC  = "https://atcoder.jp/contests/abc"
	prefixARC  = "https://atcoder.jp/contests/arc"
	midfixABC  = "/tasks/abc"
	midfixARC  = "/tasks/arc"
	suffix  = "/editorial"
	EdUrl   []string
	keyword []string

	// Synchronization
	wg sync.WaitGroup
	mu sync.Mutex
)

// Generate URLs for contests
func generateProblemURLs(startId, endId int) []string {
	var ids []string
	for i := startId; i <= endId; i++ {
		for _, Id := range problemIdABC{
			ids = append(ids, fmt.Sprintf("%s%d%s%d_%s%s", prefixABC, i, midfixABC, i, Id, suffix))
		}
		for _, Id := range problemIdARC{
			ids = append(ids, fmt.Sprintf("%s%d%s%d_%s%s", prefixARC, i, midfixARC, i, Id, suffix))
		}
	}
	return ids
}

// Process CLI arguments
// Contest range and keywords to search
func processArgs(args []string) (int, int) {
	startId, err1 := strconv.Atoi(strings.TrimSpace(args[0]))
	endId, err2 := strconv.Atoi(strings.TrimSpace(args[1]))

	if err1 != nil || err2 != nil {
		fmt.Println("Error: startId and endId must be integers")
		os.Exit(1)
	}

	if startId > endId {
		startId, endId = endId, startId
	}

	// Process keywords
	tempKeyword := strings.Join(args[2:], " ")
	for _, str := range strings.Split(tempKeyword, ",") {
		keyword = append(keyword, strings.ToLower(strings.TrimSpace(str)))
	}

	return startId, endId
}

// Extract editorial links
func extractLinks(body io.Reader) string {
	tokenizer := html.NewTokenizer(body)
	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			return ""
		case html.StartTagToken:
			token := tokenizer.Token()
			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						link := attr.Val
						if strings.Contains(link, "editorial") && !strings.Contains(link, "?") {
							afterSplit := strings.Split(link, "/")
							if len(afterSplit) == 5 {
								return "https://atcoder.jp" + link
							}
						}
					}
				}
			}
		}
	}
}

// Make an HTTP request
func makeRequest(url string, con bool) {
	defer wg.Done()

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	// Select random user-agent
	mu.Lock()
	userAgent := userAgents[counter%len(userAgents)]
	counter++
	mu.Unlock()

	req.Header.Set("User-Agent", userAgent)

	// Make request
	resp, err := client.Do(req)

	if err != nil || resp == nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	// Extract links and keyword search
	if con {
		link := extractLinks(resp.Body)
		if link != "" {
			mu.Lock()
			EdUrl = append(EdUrl, link)
			mu.Unlock()
		}
	} else {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return
		}
		bodyStr := strings.ToLower(string(body))
		for _, key := range keyword {
			if strings.Contains(bodyStr, key) {
				fmt.Println(key, " :: ", url)
			}
		}
	}
}

func main() {
	// Check arguments
	args := os.Args[1:]
	if len(args) < 2 {
		PrintHelp()
		os.Exit(1)
	}

	// Parse arguments
	startId, endId := processArgs(args)
	urls := generateProblemURLs(startId, endId)

	if len(keyword) == 0 {
		fmt.Println("Error: No keywords provided")
		os.Exit(1)
	}

	// Rate limiter: 2 requests per second
	rate := time.Second / 2
	limiter := time.NewTicker(rate)
	defer limiter.Stop()

	// Extract editorial links
	for _, url := range urls {
		<-limiter.C
		wg.Add(1)
		go makeRequest(url, true)
	}
	wg.Wait()

	// Wait 5 Minute
	fmt.Printf("\n\nWaiting 5 Minute\n\n")
	time.Sleep(5 * time.Minute)

	// Search for keywords in extracted links
	for _, url := range EdUrl {
		<-limiter.C
		wg.Add(1)
		go makeRequest(url, false)
	}
	wg.Wait()
}
