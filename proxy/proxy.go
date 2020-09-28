package proxy

import (
    "fmt"
    "io"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
    "os"
    "regexp"
    "strings"
)

const PATH_NOT_FOUND = "Not found!"

type Downstream struct {
    address      string
    allowedList  []*regexp.Regexp
}

type Proxy struct {
    downstream Downstream
    logPrefix  string
    logger     io.Writer
}

func NewProxy(config Config) *Proxy {
    prx := &Proxy{
        downstream: Downstream{
            address:      config.DownstreamUrl,
            allowedList:  nil,
        },
        logPrefix: config.LogPrefix,
        logger:    nil,
    }

    prx.initPaths(config.Paths.Params, config.Paths.Allowed)
    prx.initLog(config.LogFile)
    return prx
}

func (proxy *Proxy) initLog(logFile string) {
    f, err := os.OpenFile(logFile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("error opening file: %v", err)
    }
    defer f.Close()
    proxy.logger = f
}

func (proxy *Proxy) initPaths(params map[string]string, allowedPaths []string) {
    proxy.downstream.allowedList = make([]*regexp.Regexp, len(allowedPaths))
    for key, path := range allowedPaths {
        for param, value := range params {
            path = strings.ReplaceAll(path, "{" + param + "}", value)
        }
        fmt.Println("Adding new path to allowed list : " + path)
        proxy.downstream.allowedList[key] = regexp.MustCompile(`^` + path + `$`)
    }
}

func (proxy *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if proxy.ValidatePath(r.RequestURI) {
        url, _ := url.Parse(proxy.downstream.address)
        handler := httputil.NewSingleHostReverseProxy(url)
        handler.ErrorLog = log.New(proxy.logger, proxy.logPrefix, log.LstdFlags)
        handler.ServeHTTP(w, r)
        return
    }

    w.WriteHeader(http.StatusNotFound)
    fmt.Fprint(w, PATH_NOT_FOUND)
}

func (proxy *Proxy) ValidatePath(path string) bool {
    for _, allowedPath := range proxy.downstream.allowedList {
        if allowedPath.MatchString(strings.Trim(strings.ToLower(path), "/")) {
            return true
        }
    }

    return false
}
