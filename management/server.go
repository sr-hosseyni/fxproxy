package management

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
)

type server struct {
    *http.Server
}

type healthResponse struct {
    ServiceStatus int `json:"service_status"`
    FailureIndex  int `json:"failure_index"`
}

func NewServer(config Config) *server {
    return &server{
        Server: &http.Server{
            Addr:    config.Host + ":" + config.Port,
            Handler: initRouter(config.HealthCheckUrl),
        },
    }
}

func initRouter(HealthCheckUrl string) *http.ServeMux {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Fxproxy management system!")
    })

    http.HandleFunc("/health-check", func(w http.ResponseWriter, r *http.Request) {

        fmt.Println("Health check request!")
        resp, err := http.Get(HealthCheckUrl)
        if err != nil {
            log.Fatalln(err)
        }

        respBody, err := json.Marshal(&healthResponse{
            ServiceStatus: resp.StatusCode,
            FailureIndex:  0, // Dummy value
        })

        if err != nil {
            fmt.Fprintln(w, err)
            log.Println(err)
        }

        fmt.Println(string(respBody))
        w.Write(respBody)
    })

    return http.DefaultServeMux
}

func (server *server) Run() {
    fmt.Printf("Management server is starting on %s\n", server.Addr)

    go func() {
        log.Fatal(server.ListenAndServe())
    }()
}
