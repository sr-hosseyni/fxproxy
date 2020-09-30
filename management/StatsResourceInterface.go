package management

type StatsResourceInterface interface {
    GetServiceName() string
    GetStats() map[string]int
}
