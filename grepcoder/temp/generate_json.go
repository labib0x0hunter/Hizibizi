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

var tags = map[string][]string{
	"rooted tree":                  {"rooted tree"},
	"tree":                         {"tree"},
	"vertices":                     {"graph"},
	"vertex":                       {"graph"},
	"virtual vertex":               {"graph"},
	"dfs":                          {"graph", "dfs"},
	"bfs":                          {"graph", "bfs"},
	"graph":                        {"graph"},
	"dijkstra":                     {"graph", "dijkstra", "shortest-path"},
	"shortest path":                {"graph", "shortest-path"},
	"warshall-floyd":               {"shortest-path", "graph", "floyd-warshall"},
	"floyd-warshall":               {"shortest-path", "graph", "floyd-warshall"},
	"bellman-ford":                 {"shortest-path", "graph", "bellman-ford"},
	"strongly connected component": {"graph", "scc"},
	"strongly-connected component": {"graph", "scc"},
	"bipartite graph":              {"graph", "bipartite-graph"},
	"maxflow":                      {"graph", "flow", "maxflow"},
	"maximum flow":                 {"graph", "flow", "maxflow"},
	"minimum cost flow":            {"graph", "flow", "mincostflow"},
	"mincostflow":                  {"graph", "flow", "mincostflow"},
	"atcoder::mcf_graph":           {"graph", "flow"},
	"dinic’s algorithm":            {"graph", "dinic", "flow"},
	"topological sort":             {"graph", "topsort"},
	"topological order":            {"graph", "topsort"},
	"bridge":                       {"graph", "bridge"},
	"articulation point":           {"graph", "articulation-point"},
	"2-sat":                        {"graph", "2-sat"},
	"trémaux tree":                 {"tree"},
	"binary tree":                  {"tree", "binary-tree"},
	"functional graph":             {"graph", "func-graph"},
	"permutation graph":            {"graph", "perm-graph"},
	"namori graph":                 {"graph", "namori-graph"},
	"minimum spanning tree":        {"tree", "mst"},
	"kruskal":                      {"tree", "mst", "kruskal"},
	"01 on tree":                   {"tree", "01-on-tree"},
	"disjoint set union":           {"dsu"},
	"dsu":                          {"dsu"},
	"union-find":                   {"dsu"},
	"cumulative sum":               {"prefix-sum"},
	"cumulative-sum":               {"prefix-sum"},
	"prefix sum":                   {"prefix-sum"},
	"prefix-sum":                   {"prefix-sum"},
	"difference array":             {"difference-array"},
	"imos method":                  {"prefix-sum", "imos-method"},
	"heap":                         {"ds", "priority-queue"},
	"priority_queue":               {"ds", "priority-queue"},
	"stack":                        {"ds", "stack"},
	"queue":                        {"queue", "ds"},
	"deque":                        {"deque", "ds"},
	"multiset":                     {"ds", "multiset"},
	"map":                          {"ds", "map"},
	"set":                          {"ds", "set"},
	"associative array":            {"ds", "associative-array"},
	"square root decomposition":    {"adv-ds", "sqrt-decomp"},
	"mo’s algorithm":               {"adv-ds", "sqrt-decomp", "mo-algo"},
	"segment tree":                 {"adv-ds", "segment-tree"},
	"segtree":                      {"adv-ds", "segment-tree"},
	"fenwick tree":                 {"adv-ds", "bi-tree"},
	"binary index tree":            {"adv-ds", "bi-tree"},
	"lazy segment tree":            {"adv-ds", "lazy-segment-tree"},
	"lazyseg":                      {"adv-ds", "lazy-segment-tree"},
	"lca":                          {"tree", "lca", "adv-ds"},
	"level ancestor":               {"tree", "adv-ds", "level-ancestor", "lca"},
	"euler tour":                   {"tree", "euler-tour"},
	"virtual tree":                 {"tree", "virtual-tree", "adv-ds"},
	"heavy-light decomposition":    {"tree", "adv-ds", "hld"},
	"centroid decomposition":       {"tree", "adv-ds", "cent-decomp"},
	"binary lifting":               {"tree", "adv-ds", "binary-lifting"},
	"persistent segment tree":      {"adv-ds", "persistent-segment-tree"},
	"wavelet tree":                 {"adv-ds", "wavelet-tree"},
	"zobrist hash":                 {"adv-ds", "zobrist-hash", "hashing"},
	"dynamic programming":          {"dp"},
	"knapsack":                     {"dp", "knapsack"},
	"digit dp":                     {"dp", "digit-dp"},
	"subset dp":                    {"dp", "subset-dp"},
	"inline dp":                    {"dp", "inline-dp"},
	"bitmask dp":                   {"dp", "bitmask-dp"},
	"bit dp":                       {"dp", "bitmask-dp"},
	"tree dp":                      {"dp", "tree-dp"},
	"matrix exponent":              {"matrix-expo"},
	"dp optimizations":             {"dp"},
	"convex hull trick":            {"dp", "convex-hull"},
	"divide and conquer dp":        {"dp", "dnc"},
	"traveling salesman":           {"dp", "tsp-dp"},
	"memorization in recursion":    {"dp"},
	"number theory":                {"num-theory"},
	"gcd":                          {"gcd"},
	"greatest common divisor":      {"gcd"},
	"lcm":                          {"lcm"},
	"multiple":                     {"lcm"},
	"modular arithmetic":           {"mod"},
	"modular inverse":              {"mod-inv"},
	"combinatorics":                {"comb"},
	"binomial coefficient":         {"comb", "bino-coef"},
	"factorial":                    {"factorial"},
	"prime":                        {"prime"},
	"sieve":                        {"sieve"},
	"factorization":                {"factorization"},
	"totient":                      {"phi"},
	"matrix":                       {"matrix"},
	"fft":                          {"fft"},
	"ntt":                          {"ntt"},
	"polynomial":                   {},
	"mobius function":              {"mobius"},
	"chinese remainder theorem":    {"crt"},
	"pigeonhole principle":         {"pigeonhole"},
	"inclusion-exclusion":          {"inc-exc"},
	"inclusion exclusion":          {"inc-exc"},
	"contribution technique":       {"contribution-technique"},
	"binomial theorem":             {"binomial-theorem"},
	"lucas number":                 {},
	"harmonic series":              {"hermonic-series"},
	"kirchhoff’s theorem":          {},
	"laplacian matrix":             {},
	"square number":                {"square-number"},
	"grundy number":                {},
	"combinatorial interpretation": {},
	"eulerian number":              {},
	"discrete logarithm":           {},
	"expected value":               {"expected-value"},
	"geometry":                     {},
	"convex hull":                  {"convex-hull"},
	"line sweep":                   {"sweep-line"},
	"sweep line":                   {"sweep-line"},
	"closest pair":                 {},
	"circle":                       {},
	"polygon":                      {"polygon"},
	"cross product":                {},
	"dot product":                  {},
	"sliding window":               {"sliding-window"},
	"binary search":                {"binary-search"},
	"ternary search":               {"ternary-search"},
	"upperbound":                   {"binary-search"},
	"lowerbound":                   {"binary-search"},
	"meet in the middle":           {"mitm"},
	"meet-in-the-middle":           {"mitm"},
	"two pointer":                  {"two-pointer"},
	"divide and conquer":           {"dnc"},
	"next_permutation":             {"permutation"},
	"permutation":                  {"permutation"},
	"inversion number":             {"permutation"},
	"bitset":                       {"bitset"},
	"hashing":                      {"string", "hashing"},
	"rolling hash":                 {"string", "hashing"},
	"string matching":              {"string"},
	"kmp":                          {"string", "kmp"},
	"knuth-morris-pratt":           {"string", "kmp"},
	"z algorithm":                  {"string", "z-algo"},
	"z-algorithm":                  {"string", "z-algo"},
	"manacher":                     {"string", "manacher", "palindrome"},
	"aho-corasick":                 {"string", "aho-corasick"},
	"aho corasick":                 {"string", "aho-corasick"},
	"suffix array":                 {"string", "suffix-array"},
	"sa-is algorithm":              {"string", "sa-is"},
	"prefix function":              {"string", "prefix-func"},
	"suffix automaton":             {"string", "suffix-auto"},
	"trie":                         {"string", "trie"},
	"palindrome":                   {"string", "palindrome"},
	"coordinate compression":       {"coordinate compression"},
	"compressing coordinate":       {"coordinate compression"},
	"run length encoding":          {"rle"},
	"run-length encoding":          {"rle"},
	"bubble sort":                  {"bubble-sort"},
	"merge sort":                   {"merge-sort"},
	"heap sort":                    {"heap-sort"},
	"monotonic":                    {"monotonic-ds"},
	"doubling technique":           {"doubling-technique"},
	"data structure":               {"ds"},
	"clique decision problem":      {"clique-decision"},
	"np-complete":                  {"np-complete"},
	"baby-step giant-step":         {"baby-step-giant-step"},
	"__builtin_popcount":           {"bitmask"},
	"bitwise bruteforcing":         {"bitmask"},
	"independent bitwise":          {"bitmask"},
	"bit exhaustive search":        {"bitmask"},
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
	for key, value := range tags {
		if bytes.Contains(body, []byte(key)) {
			tg = append(tg, value...)
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
