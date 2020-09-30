package main

import (
    "fmt"
    "fxproxy/management"
    "fxproxy/proxy"
    "log"
    "net/http"
)

func main() {
    cfgPath, err := ParseFlags()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Starting by loading configuration from " + cfgPath)

    cfg, err := NewConfig(cfgPath)
    if err != nil {
        log.Fatal(err)
    }

    prx := proxy.NewProxy(cfg.ProxyServer)
    defer prx.Close()

    manager := management.NewServer(cfg.ManagementServer)
    manager.Run()

    defer manager.Close()

    // Run the web server.
    log.Fatal(http.ListenAndServe(":8888", prx))
}
