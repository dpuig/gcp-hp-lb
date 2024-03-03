# algorithm 

```go
func loadBalanceRoundRobin(targets []string) *httputil.ReverseProxy {
	// Convert targets to *url.URL and wrap them in our custom struct
	targetURLs := make([]*target, len(targets))
	for i, rawurl := range targets {
		u, _ := url.Parse(rawurl)
		targetURLs[i] = &target{url: u}
	}

	// Initialize current index
	currentIndex := 0

	director := func(req *http.Request) {
		// Use a mutex to ensure that only one request can select a target at a time
		mutex := &sync.Mutex{}
		mutex.Lock()
		defer mutex.Unlock()

		// Select the target using round-robin
		currentIndex = (currentIndex + 1) % len(targetURLs)

		// Rewrite the request to be sent to the selected target
		req.URL.Scheme = targetURLs[currentIndex].url.Scheme
		req.URL.Host = targetURLs[currentIndex].url.Host
		req.URL.Path = singleJoiningSlash(targetURLs[currentIndex].url.Path, req.URL.Path)
		if targetURLs[currentIndex].url.RawQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetURLs[currentIndex].url.RawQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetURLs[currentIndex].url.RawQuery + "&" + req.URL.RawQuery
		}
	}

	return &httputil.ReverseProxy{Director: director}
}
```