package proxy

type Config struct {
    LogFile       string `yaml:"log_file"`
    LogPrefix     string `yaml:"log_prefix"`
    DownstreamUrl string `yaml:"downstream_url"`
    Paths         struct {
        Params  map[string]string `yaml:"params"`
        Allowed []string          `yaml:"allowed"`
    } `yaml:"paths"`
}
