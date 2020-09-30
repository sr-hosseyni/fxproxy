package proxy

import (
    "fmt"
    "log"
    "net"
    "net/http"
    "net/http/httputil"
    "net/url"
    "os"
    "regexp"
    "strings"
    "time"
)

const PATH_NOT_FOUND = "Not found!"

type downstream struct {
    address     string
    allowedList []*regexp.Regexp
}

type proxy struct {
    serviceName  string
    downstream   *downstream
    handler      *httputil.ReverseProxy
    timeout      *time.Duration
    logPrefix    string
    errorLogger  *os.File
    accessLogger *os.File

    // very simple index for measuring downstream faults
    failedRequestsCount  int
    totalRequestsCount   int
    invalidRequestsCount int
}

func NewProxy(config Config) *proxy {
    prx := &proxy{
        serviceName: config.ServiceName,
        handler: nil,
        timeout: &config.Timeout,
        downstream: &downstream{
            address:     config.DownstreamUrl,
            allowedList: nil,
        },
        logPrefix:            config.Logs.Prefix,
        errorLogger:          nil,
        accessLogger:         nil,
        failedRequestsCount:  0,
        totalRequestsCount:   0,
        invalidRequestsCount: 0,
    }

    prx.initPaths(config.Paths.Params, config.Paths.Allowed)
    prx.initLog(config.Logs.AccessFile, config.Logs.ErrorFile)
    prx.initHandler()
    return prx
}

func (proxy *proxy) initLog(AccessLogFile string, ErrorLogFile string) {
    if ErrorLogFile != "" {
        f, err := os.OpenFile(ErrorLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
        if err != nil {
            log.Fatalf("error opening file: %v", err)
        }
        proxy.errorLogger = f
    }

    if AccessLogFile != "" {
        f, err := os.OpenFile(AccessLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
        if err != nil {
            log.Fatalf("error opening file: %v", err)
        }
        proxy.accessLogger = f
    }
}

func (proxy *proxy) initHandler() {
    url, _ := url.Parse(proxy.downstream.address)
    proxy.handler = httputil.NewSingleHostReverseProxy(url)
    proxy.handler.ErrorLog = log.New(proxy.errorLogger, proxy.logPrefix, log.LstdFlags)

    // if there is some specific rule for changing url
    if false {
        proxy.handler.Director = func(req *http.Request) {
            //req.Header.Add("X-Forwarded-Host", req.Host)
            //req.Header.Add("X-Origin-Host", proxy.downstream.address)
            //r.URL.Scheme = url.Scheme
            //r.URL.Host = url.Host
            //r.URL.Path = url.Path + r.URL.Path
        }
    }

    if proxy.timeout != nil {
        proxy.handler.Transport = &http.Transport{
            DialContext: (&net.Dialer{
                Timeout: *proxy.timeout,
            }).DialContext,
        }
    }

    proxy.handler.ModifyResponse = func(r *http.Response) error {
        if proxy.accessLogger != nil {
            fmt.Fprintf(
                proxy.accessLogger,
                "%s %s %s %s %s %s\n",
                proxy.logPrefix,
                time.Now().Format(time.RFC3339Nano),
                r.Request.RemoteAddr,
                r.Request.Method,
                r.Request.RequestURI,
                r.Status,
            )
        }

        if r.StatusCode >= 500 {
            proxy.failedRequestsCount++
        }

        return nil
    }

    proxy.handler.ErrorHandler = func(rw http.ResponseWriter, r *http.Request, err error) {
        fmt.Printf("error: %+v", err)

        if proxy.errorLogger != nil {
            fmt.Fprintln(proxy.errorLogger, err)
        }

        rw.WriteHeader(http.StatusInternalServerError)
        rw.Write([]byte(err.Error()))
    }
}

func (proxy *proxy) initPaths(params map[string]string, allowedPaths []string) {
    proxy.downstream.allowedList = make([]*regexp.Regexp, len(allowedPaths))
    for key, path := range allowedPaths {
        for param, value := range params {
            path = strings.ReplaceAll(path, "{"+param+"}", value)
        }
        fmt.Println("Adding new path to allowed list : " + path)
        proxy.downstream.allowedList[key] = regexp.MustCompile(`^` + path + `$`)
    }
}

func (proxy *proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    proxy.totalRequestsCount++
    if proxy.ValidatePath(r.RequestURI) {
        proxy.handler.ServeHTTP(w, r)
        return
    }

    proxy.invalidRequestsCount++
    fmt.Fprintf(
        proxy.accessLogger,
        "%s %s %s %s %s %d Path is not allowed by proxy!\n",
        proxy.logPrefix,
        time.Now().Format(time.RFC3339Nano),
        r.RemoteAddr,
        r.Method,
        r.RequestURI,
        http.StatusNotFound,
    )

    w.WriteHeader(http.StatusNotFound)
    fmt.Fprint(w, PATH_NOT_FOUND)
}

func (proxy *proxy) ValidatePath(path string) bool {
    for _, allowedPath := range proxy.downstream.allowedList {
        if allowedPath.MatchString(strings.Trim(strings.ToLower(strings.Split(path, "?")[0]), "/")) {
            return true
        }
    }

    return false
}

func (proxy *proxy) GetServiceName() string {
    return proxy.serviceName
}

func (proxy *proxy) GetStats() map[string]int {
    return map[string]int{
        "failedRequestsCount":  proxy.failedRequestsCount,
        "totalRequestsCount":   proxy.totalRequestsCount,
        "invalidRequestsCount": proxy.invalidRequestsCount,
    }
}

func (proxy *proxy) Close() {
    proxy.errorLogger.Close()
    proxy.accessLogger.Close()
}
