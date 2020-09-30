package management

type Config struct {
    Host            string            `yaml:"host"`
    Port            string            `yaml:"port"`
    HealthCheckUrls map[string]string `yaml:"downstream_health_check_urls"`
}
