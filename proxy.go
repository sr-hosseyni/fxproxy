package main

import (
    "net/http"
    "regexp"
)

type Proxy struct {
    downstream DownStream
    allowedList []*regexp.Regexp
}

func (proxy *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
