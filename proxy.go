package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
    "regexp"
    "strings"
)

const PROXY_NOT_FOUND = "Not found!"

type Proxy struct {
    downstream DownStream
    allowedList []*regexp.Regexp
    logger io.Writer
    logPrefix string
}

func (proxy *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if proxy.ValidatePath(r.RequestURI) {
        url, _ := url.Parse(proxy.downstream.Address)
        handler := httputil.NewSingleHostReverseProxy(url)
        handler.ErrorLog = log.New(proxy.logger, proxy.logPrefix, log.LstdFlags)
        handler.ServeHTTP(w, r)
        return
    }

    w.WriteHeader(http.StatusNotFound)
    fmt.Fprint(w, PROXY_NOT_FOUND)
}

func (proxy *Proxy)ValidatePath(path string) bool {
    for _, allowedPath := range proxy.allowedList {
        if allowedPath.MatchString(strings.Trim(strings.ToLower(path), "/")) {
            return true
        }
    }

    return false
}
