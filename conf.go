package main

import (
    "flag"
    "fmt"
    "fxproxy/management"
    "fxproxy/proxy"
    "gopkg.in/yaml.v3"
    "os"
)

type Config struct {
    ProxyServer struct {
        proxy.Config `yaml:",inline"`
        Host         string `yaml:"host"`
        Port         string `yaml:"port"`
        DSHealthCheckUrl string `yaml:"downstream_health_check_url"`
    } `yaml:"proxy_server"`
    ManagementServer management.ManagementServerConfig `yaml:"management_server"`
}

func NewConfig(configPath string) (*Config, error) {
    config := &Config{}

    file, err := os.Open(configPath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    d := yaml.NewDecoder(file)

    if err := d.Decode(&config); err != nil {
        return nil, err
    }

    return config, nil
}

// Makes sure that the path provided is a file and can be read
func ValidateConfigPath(path string) error {
    s, err := os.Stat(path)
    if err != nil {
        return err
    }
    if s.IsDir() {
        return fmt.Errorf("'%s' is a directory, not a normal file", path)
    }
    return nil
}

// ParseFlags will create and parse the CLI flags and return the path to be used elsewhere
func ParseFlags() (string, error) {
    var configPath string

    flag.StringVar(&configPath, "config", "/etc/fxproxy.yml", "path to config file")
    flag.Parse()

    if err := ValidateConfigPath(configPath); err != nil {
        return "", err
    }

    return configPath, nil
}
