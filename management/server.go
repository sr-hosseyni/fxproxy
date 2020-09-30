package management

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
)

type server struct {
    *http.Server
    stats           []StatsResourceInterface
    healthCheckUrls map[string]string
}

type healthResponse struct {
    ServiceStatus string                    `json:"service_status"`
    Stats         map[string]map[string]int `json:"stats"`
}

func NewServer(config Config) *server {
    server := &server{
        Server: &http.Server{
            Addr: config.Host + ":" + config.Port,
        },
        stats:           []StatsResourceInterface{},
        healthCheckUrls: config.HealthCheckUrls,
    }

    server.initRouter()
    return server
}

func (server *server) initRouter() {

    router := http.NewServeMux()

    router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Fxproxy management system!")
    })

    router.HandleFunc("/health-check", func(w http.ResponseWriter, r *http.Request) {

        fmt.Println("Health check request!")

        respBody, err := json.Marshal(server.getStats())

        if err != nil {
            fmt.Fprintln(w, err)
            log.Println(err)
        }

        fmt.Println(string(respBody))
        w.Write(respBody)
    })

    server.Handler = router
}

func (server server) callHealthCheck(serviceName string) int {
    resp, err := http.Get(server.healthCheckUrls[serviceName])
    if err != nil {
        log.Fatalln(err)
    }

    return resp.StatusCode
}

func (server server) getStats() *healthResponse {
    res := &healthResponse{
        ServiceStatus: "OK",
        Stats:         map[string]map[string]int{},
    }

    for _, resource := range server.stats {
        name := resource.GetServiceName()
        res.Stats[name] = resource.GetStats()
        res.Stats[name]["health-check"] = server.callHealthCheck(name)
        if res.Stats[name]["health-check"] != http.StatusOK {
            res.ServiceStatus = "unhealthy"
        }
    }

    return res
}

func (server *server) AddStatsResource(resource StatsResourceInterface) {
    server.stats = append(server.stats, resource)
}

func (server *server) Run() {
    fmt.Printf("Management server is starting on %s\n", server.Addr)

    go func() {
        log.Fatal(server.ListenAndServe())
    }()
}
