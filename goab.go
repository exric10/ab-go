package main

import (
        "fmt"
        "io"
		"os"
        "io/ioutil"
        "net/http"
		"github.com/tcnksm/go-httpstat"
		"net/http/httptrace"
        "log"
        "time"
		"strconv"
		"sync"
)

var latency time.Duration
var errors int = 0
var kareq int = 0

func readingArgs(args [] string) (int, int, bool, string) {
	var n, c int = 1, 1
	var k bool = false
	var link string = args[len(args)-1]

	if len(args) == 2 {
		return n, c, k , link
	}

	for index, element := range args[0:len(args)-1] {
		if element == "-n" {
			n, _= strconv.Atoi(args[index+1])
		} else if element == "-c" {
			c, _ = strconv.Atoi(args[index+1])
		} else if element == "-k" {
			k = true
		}
	}
	return n, c, k, link
}

func testWeb (n, c int, k bool, link string) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	if k {
		for i := 0; i < n; i += c {
			if i+c > n {
				c = n%c
			}
			wg.Add(c)
			for j := 0; j < c; j++ {
				go getRequestWithKeepAlive(link, &wg, &mu)
			}
			wg.Wait()
		}
	} else {
		for i := 0; i < n; i += c {
			if i+c > n {
				c = n%c
			}
			wg.Add(c)
			for j := 0; j < c; j++ {
				go getRequestWithoutKeepAlive(link, &wg, &mu)
			}
			wg.Wait()
		}
	}
}

func getRequestWithKeepAlive(link string, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()

	req, err := http.NewRequest("GET",link,nil)
	if err != nil {
		log.Fatal(err)
	}

	var result httpstat.Result
	ctx := httpstat.WithHTTPStat(req.Context(), &result)
	req = req.WithContext(ctx)

	trace := &httptrace.ClientTrace {
		GotConn: func(connInfo httptrace.GotConnInfo) {
			if connInfo.Reused {
				kareq += 1
			}
			//print used for checking if Keep-Alive feature works
			//fmt.Printf("Reused Conn: %+v\n", connInfo.Reused)
		},
	}

	mu.Lock()
	req = req.WithContext(httptrace.WithClientTrace(req.Context(),trace))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		errors += 1
	}

	if _, err := io.Copy(ioutil.Discard, resp.Body); err != nil {
		log.Fatal(err)
	}

	resp.Body.Close()
	
	latency += result.Total(time.Now())
	mu.Unlock()
}

func getRequestWithoutKeepAlive(link string, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()

	req, err := http.NewRequest("GET",link,nil)
	if err != nil {
		log.Fatal(err)
	}

	var result httpstat.Result
	ctx := httpstat.WithHTTPStat(req.Context(), &result)
	
	req = req.WithContext(ctx)

	trace := &httptrace.ClientTrace {
		GotConn: func(connInfo httptrace.GotConnInfo) {
			if connInfo.Reused {
				kareq += 1
			}
			//print used for checking if Keep-Alive feature works
			//fmt.Printf("Reused Conn: %+v\n", connInfo.Reused)
		},
	}

	mu.Lock()
	req = req.WithContext(httptrace.WithClientTrace(req.Context(),trace))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		errors += 1
	}

	resp.Body.Close()

	latency += result.Total(time.Now())
	mu.Unlock()
}

func main() {
	args := os.Args
	n, c, k, link := readingArgs(args)

	testWeb(n,c,k,link)

	fmt.Printf("Concurrency level: %d\n", c)
	fmt.Printf("Time taken for tests: %g\n", latency.Seconds())
	fmt.Printf("Completed requests: %d\n", n-errors)
	fmt.Printf("Errored responses: %d (amount), %g (percentage)\n",errors, (float64(errors)/float64(n)) * 100)
	fmt.Printf("Keep-Alive requests: %d\n", kareq)
	fmt.Printf("Transactions per second: %g [#/sec] (mean)\n", float64(n)/latency.Seconds())
	fmt.Printf("Time per request: %g [ms] (mean)\n", ((float64(c)*latency.Seconds()*1000)/float64(n)))
	fmt.Printf("Time per request: %g [ms] (mean,across all concurrent requests)\n", ((latency.Seconds()*1000)/float64(n)))
}