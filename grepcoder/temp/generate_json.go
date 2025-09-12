package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/projectdiscovery/useragent"
)

const kenkooApiUrl = "https://kenkoooo.com/atcoder/resources/problems.json"

var tags = []string{
	// Graphs
	"tree",
	"dfs",
	"bfs",
	"graph",
	"dijkstra",
	"shortest path",
	"warshall-floyd",
	"strongly connected component",
	"bipartite graph",
	"maxflow",
	"mincostflow",
	"atcoder::mcf_graph",
	"topological sort",
	"bridge",
	"articulation point",
	"2-sat",

	// DSU / Union-Find
	"disjoint set union",
	"dsu",

	// Prefix / Accumulation
	"cumulative sum",
	"prefix sum",
	"prefix-sum",
	"difference array",
	"imos method",

	// Basic DS
	"priority_queue",
	"stack",
	"queue",
	"deque",
	"multiset",
	"map",
	"unordered_map",

	// Advanced DS
	"square root decomposition",
	"segment tree",
	"fenwick tree",
	"lazy segment tree",
	"lca",
	"euler tour",
	"virtual tree",
	"heavy-light decomposition",
	"centroid decomposition",
	"binary lifting",
	"persistent segment tree",
	"wavelet tree",

	// DP
	"dynamic programming",
	"knapsack",
	"digit dp",
	"subset dp",
	"inline dp",
	"bitmask dp",
	"matrix exponentiation",
	"dp optimizations",
	"convex hull trick",
	"divide and conquer dp",

	// Math
	"number theory",
	"gcd",
	"lcm",
	"modular arithmetic",
	"modular inverse",
	"combinatorics",
	"binomial coefficient",
	"factorial",
	"prime sieve",
	"totient",
	"matrix",
	"fft",
	"ntt",
	"polynomial",
	"mobius function",
	"chinese remainder theorem",
	"pigeonhole principle",

	// Geometry
	"geometry",
	"convex hull",
	"line sweep",
	"closest pair",
	"circle",
	"polygon",
	"area",
	"cross product",
	"dot product",

	// General Algorithms
	"sliding window",
	"binary search",
	"ternary search",
	"upperbound",
	"lowerbound",
	"meet in the middle",
	"two pointers",
	"greedy",

	// Misc
	"next_permutation",
	"bitset",
	"hashing",
	"rolling hash",
	"string matching",
	"kmp",
	"z algorithm",
	"manacher",
	"aho-corasick",
	"aho corasick",
	"suffix array",
	"suffix automaton",
	"trie",
}

type problemModel struct {
	Id           string   `json:"id"`
	Name         string   `json:"name"`
	ContestId    string   `json:"contest_id"`
	EditorialUrl string   `json:"edi_url"`
	Tags         []string `json:"tag"`
}

func (p *problemModel) setEditorialUrl(url string) {
	p.EditorialUrl = url
}

func (p *problemModel) setTag(tag []string) {
	p.Tags = tag
}

var finalResult []problemModel
var client *http.Client
var rate time.Duration
var limiter *time.Ticker

func init() {
	finalResult = make([]problemModel, 0, 20000)
	client = &http.Client{Timeout: 10 * time.Second}
	rate = time.Second / 2
	limiter = time.NewTicker(rate)

}

func extractEditorialUrl(body io.ReadCloser, ContestId string) string {
	// Extract editorial links
	// tokenizer := html.NewTokenizer(body)
	// for {
	// 	tokenType := tokenizer.Next()
	// 	switch tokenType {
	// 	case html.ErrorToken:
	// 		return ""
	// 	case html.StartTagToken:
	// 		token := tokenizer.Token()
	// 		if token.Data == "a" {
	// 			for _, attr := range token.Attr {
	// 				if attr.Key == "href" {
	// 					link := attr.Val

	// 					pattern := fmt.Sprintf(`<a href="/contests/%s/editorial/`, ContestId)

	// 					if ok, err := regexp.MatchString(pattern, link); err == nil && ok {
	// 						return link
	// 					} else if err != nil {
	// 						return ""
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	// ptr := fmt.Sprintf(`<a href="/contests/%s/editorial/`, ContestId)
	// pattern := []byte(ptr)
	content, err := io.ReadAll(body)
	if err != nil {
		return ""
	}

	// for i := 0; i+len(pattern) < len(content); i++ {
	// 	if bytes.Equal(content[i:i+len(pattern)], pattern) {
	// 		valid := true
	// 		edId := ""
	// 		for j := i + len(pattern); j < len(content); j++ {
	// 			if content[j] == 34 {
	// 				break
	// 			}
	// 			if content[j] < 48 || content[j] > 57 {
	// 				// fmt.Println("content[j]=", content[j])
	// 				valid = false
	// 				break
	// 			}
	// 			edId += string(content[j])
	// 		}
	// 		// fmt.Println("Valid=", valid)
	// 		if valid {
	// 			// fmt.Println("VALis")
	// 			return fmt.Sprintf("/contests/%s/editorial/%s", ContestId, edId)
	// 		}
	// 	}
	// }

	re := regexp.MustCompile(fmt.Sprintf(`/contests/%s/editorial/\d+`, ContestId))
	matches := re.Find(content)
	if matches != nil {
		return "https://atcoder.jp" + string(matches)
	}
	return ""
}

func findTags(pl problemModel) (tg []string) {
	// resp, err := http.Get(pl.EditorialUrl)
	// if err != nil {
	// 	panic(err)
	// }

	// client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", pl.EditorialUrl, nil)
	if err != nil {
		return
	}

	userAgent := useragent.PickRandom()
	req.Header.Set("User-Agent", userAgent.Raw)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	body = bytes.ToLower(body)
	for _, key := range tags {
		if bytes.Contains(body, []byte(key)) {
			tg = append(tg, key)
		}
	}
	return
}

func getEditorialUrl(pl problemModel) {
	url := fmt.Sprintf("https://atcoder.jp/contests/%s/tasks/%s/editorial", pl.ContestId, pl.Id)

	// resp, err := http.Get(url)
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()

	// client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	userAgent := useragent.PickRandom()
	req.Header.Set("User-Agent", userAgent.Raw)

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	url = extractEditorialUrl(resp.Body, pl.ContestId)
	if url == "" {
		finalResult = append(finalResult, pl)
		return
	}

	pl.setEditorialUrl(url)

	tags := findTags(pl)
	if len(tags) == 0 {
		finalResult = append(finalResult, pl)
		return
	}

	pl.setTag(tags)

	fmt.Println("Got=", pl)

	finalResult = append(finalResult, pl)
}

func main() {

	defer limiter.Stop()

	// client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", kenkooApiUrl, nil)
	if err != nil {
		return
	}

	userAgent := useragent.PickRandom()
	req.Header.Set("User-Agent", userAgent.Raw)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	pl := make([]problemModel, 20000)
	if err := json.Unmarshal(content, &pl); err != nil {
		panic(err)
	}

	for i := range pl {
		if strings.HasPrefix(pl[i].Id, "abc") || strings.HasPrefix(pl[i].Id, "arc") {
			// // <-limiter.C
			// if pl[i].Id == "abc375_d" {
			fmt.Println("Process=", pl[i].Id)
			getEditorialUrl(pl[i])
			time.Sleep(1 * time.Second)
			// }
		}
	}

	out := "gstore.json"
	f, err := os.OpenFile(out, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	b, err := json.Marshal(finalResult)
	if err != nil {
		panic(err)
	}

	if _, err = f.Write(b); err != nil {
		panic(err)
	}

	// content, err := io.ReadAll(f)
	// if err != nil {
	// 	panic(err)
	// }
	// if err := json.Unmarshal(content, &finalResult); err != nil {
	// 	panic(err)
	// }

	// fmt.Println(finalResult[0])
}
