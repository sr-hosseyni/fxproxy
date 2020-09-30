package management

type Config struct {
    Host             string `yaml:"host"`
    Port             string `yaml:"port"`
    HealthCheckUrl string `yaml:"downstream_health_check_url"`
}
