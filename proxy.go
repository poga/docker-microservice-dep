package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
)

func NewProxy(s Service) (string, *httputil.ReverseProxy) {
	prefix := "/" + s.Name
	urlStr := GetFromEtcd(s)
	fmt.Println(urlStr)
	url, err := url.Parse(urlStr)
	fmt.Printf("%v\n", url)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Director = direct(prefix, url)
	return prefix, proxy
}

func direct(prefix string, target *url.URL) func(req *http.Request) {
	regex := regexp.MustCompile(`^` + prefix)
	fmt.Printf("Proxying %s to %s\n", prefix, target)
	return func(req *http.Request) {
		targetQuery := target.RawQuery
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = regex.ReplaceAllString(singleJoiningSlash(target.Path, req.URL.Path), "")
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
	}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
