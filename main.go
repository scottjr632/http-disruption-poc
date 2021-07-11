package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sync"
)

func main() {

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		log.Printf("REQUEST \t %s %s %v", r.Method, r.URL.Path, r.Header)
		if r.Method == http.MethodPost {
			log.Printf("DROPPING REQUEST TO %s\n", r.URL.Host+r.URL.Path)
			return
		}

		var scheme string
		if r.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}

		// Hacky way to get rid of the port since our ip table rules only filter on
		// http or https traffic we don't have to worry about a different port.
		re := regexp.MustCompile("(:.*$)")
		host := re.ReplaceAllString(r.Host, "")
		proxyURI := fmt.Sprintf("%s://%s%s", scheme, host, r.URL.Path)
		log.Printf("Proxying to %s\n", proxyURI)

		req, err := http.NewRequest(http.MethodPost, proxyURI, nil)
		if err != nil {
			log.Fatal(err)
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
			return
		}

		rw.Write(body)
	})

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		log.Println("Starting TLS server on :9090")
		log.Fatal(http.ListenAndServeTLS(":9090", "./google.com+5.pem", "./google.com+5-key.pem", nil))
		wg.Done()
	}()

	go func() {
		log.Println("Starting HTTP server on :8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
		wg.Done()
	}()

	wg.Wait()
}
