package proxy

import "time"

type Config struct {
    Timeout time.Duration `yaml:"proxy_timeout"`
    Logs struct {
        ErrorFile  string `yaml:"error_file"`
        AccessFile string `yaml:"access_file"`
        Prefix     string `yaml:"prefix"`
    } `yaml:"logs"`
    DownstreamUrl string `yaml:"downstream_url"`
    Paths         struct {
        Params  map[string]string `yaml:"params"`
        Allowed []string          `yaml:"allowed"`
    } `yaml:"paths"`
}
