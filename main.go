package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sync"
)

const (
	tlsPort  = ":9090"
	httpPort = ":8080"

	certFile = "./localhost+2.pem"
	keyFile  = "./localhost+2-key.pem"
)

func main() {

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		log.Printf("REQUEST \t %s %s %v", r.Method, r.URL.Path, r.Header)
		if r.Method == http.MethodPost {
			log.Printf("DROPPING REQUEST TO %s\n", r.Host+r.URL.Path)
			return
		}

		var scheme string
		if r.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}

		// Hacky way to get rid of the port in our host since our ip table rules only filter
		// on ports 5000 and 4000. This prevents a looping scenario where the requests are
		// consintely looped back to the proxy due to the ip table rules. Outside of the POC
		// the ip table rule will be more specific to where we will not have to worry about
		// any looping.
		re, err := regexp.Compile("(:.*$)")
		if err != nil {
			log.Fatal(err)
			return
		}

		host := re.ReplaceAllString(r.Host, "")
		proxyURI := fmt.Sprintf("%s://%s%s", scheme, host, r.URL.Path)
		log.Printf("Proxying to %s\n", proxyURI)

		req, err := http.NewRequest(r.Method, proxyURI, r.Body)
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
		log.Printf("Starting TLS server on %s\n", tlsPort)
		log.Fatal(http.ListenAndServeTLS(tlsPort, certFile, keyFile, nil))
		wg.Done()
	}()

	go func() {
		log.Printf("Starting HTTP server on %s\n", httpPort)
		log.Fatal(http.ListenAndServe(httpPort, nil))
		wg.Done()
	}()

	wg.Wait()
}
