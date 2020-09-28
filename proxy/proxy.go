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
    address     string
    allowedList []*regexp.Regexp
}

type Proxy struct {
    downstream   Downstream
    handler      *httputil.ReverseProxy
    logPrefix    string
    errorLogger  io.Writer
    accessLogger io.Writer
}

func NewProxy(config Config) *Proxy {
    prx := &Proxy{
        downstream: Downstream{
            address:     config.DownstreamUrl,
            allowedList: nil,
        },
        logPrefix:    config.Logs.Prefix,
        errorLogger:  nil,
        accessLogger: nil,
    }

    prx.initPaths(config.Paths.Params, config.Paths.Allowed)
    prx.initLog(config.Logs.AccessFile, config.Logs.ErrorFile)
    prx.initHandler()
    return prx
}

func (proxy *Proxy) initLog(AccessLogFile string, ErrorLogFile string) {
    f, err := os.OpenFile(ErrorLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("error opening file: %v", err)
    }
    defer f.Close()
    proxy.errorLogger = f

    f, err = os.OpenFile(AccessLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("error opening file: %v", err)
    }
    defer f.Close()
    proxy.accessLogger = f
}

func (proxy *Proxy) initHandler() {
    url, _ := url.Parse(proxy.downstream.address)
    proxy.handler = httputil.NewSingleHostReverseProxy(url)
    proxy.handler.ErrorLog = log.New(proxy.errorLogger, proxy.logPrefix, log.LstdFlags)
}

func (proxy *Proxy) initPaths(params map[string]string, allowedPaths []string) {
    proxy.downstream.allowedList = make([]*regexp.Regexp, len(allowedPaths))
    for key, path := range allowedPaths {
        for param, value := range params {
            path = strings.ReplaceAll(path, "{"+param+"}", value)
        }
        fmt.Println("Adding new path to allowed list : " + path)
        proxy.downstream.allowedList[key] = regexp.MustCompile(`^` + path + `$`)
    }
}

func (proxy *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if proxy.ValidatePath(r.RequestURI) {
        proxy.handler.ServeHTTP(w, r)
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
