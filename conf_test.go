package main

import (
    "github.com/stretchr/testify/require"
    "log"
    "testing"
    "time"
)

/**
proxy-server:
  host: 1.2.3.4
  port: 5678
  log-file: /path/to/log.file
  log-prefix: test prefix
  downstream:
    url: http://nginx:80
    health-check-url: service/health-check
    paths:
      params:
        id: "[a-z]+[0-9]+[a-z0-9]+"
        num: "[0-9]+"
      allowed:
        - company
        - company/{id}
        - account/{num}
        - "{id}"
        - "{id}/{num}"
        - account/{id}/user/{num}
management-server:
  host: 127.0.0.1
  port: 3344
 */

func TestConfigReader(t *testing.T) {
    cfg, err := NewConfig("./conf_test.yaml")
    if err != nil {
        log.Fatal(err)
    }

    timeout, _ := time.ParseDuration("10s")
    require.Equal(t, timeout, cfg.ProxyServer.Timeout)
    require.Equal(t, "service-1", cfg.ProxyServer.ServiceName)
    require.Equal(t, "/path/to/log.file", cfg.ProxyServer.Logs.ErrorFile)
    require.Equal(t, "/path/to/access.file.log", cfg.ProxyServer.Logs.AccessFile)
    require.Equal(t, "test prefix", cfg.ProxyServer.Logs.Prefix)
    require.Equal(t, "http://nginx:80", cfg.ProxyServer.DownstreamUrl)
    require.Equal(t, "http://something/service/health-check", cfg.ManagementServer.HealthCheckUrls["service-1"])
    require.Equal(t, "127.0.0.1", cfg.ManagementServer.Host)
    require.Equal(t, "3344", cfg.ManagementServer.Port)
    require.Equal(t, []string{
        "company",
        "company/{id}",
        "account/{num}",
        "{id}",
        "{id}/{num}",
        "account/{id}/user/{num}",
    }, cfg.ProxyServer.Paths.Allowed)
    require.Equal(t, "[a-z]+[0-9]+[a-z0-9]+", cfg.ProxyServer.Paths.Params["id"])
    require.Equal(t, "[0-9]+", cfg.ProxyServer.Paths.Params["num"])
}
