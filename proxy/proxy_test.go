package proxy

import (
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "regexp"
    "strings"
    "testing"

    "github.com/stretchr/testify/require"
)

const ID_FORMAT = "[a-z]+[0-9]+[a-z0-9]+"

var allowedList = []*regexp.Regexp{
    regexp.MustCompile(`^company$`),                                                  // /company/
    regexp.MustCompile(strings.Replace(`^company/{id}$`, "{id}", ID_FORMAT, 1)),      // /company/{id}
    regexp.MustCompile(`^company/account$`),                                          // /company/account
    regexp.MustCompile(`^account$`),                                                  // /account
    regexp.MustCompile(strings.Replace(`^account/{id}$`, "{id}", ID_FORMAT, 1)),      // /account/{id}
    regexp.MustCompile(strings.Replace(`^{id}$`, "{id}", ID_FORMAT, 1)),              // /{id}
    regexp.MustCompile(strings.Replace(`^account/{id}/user$`, "{id}", ID_FORMAT, 1)), // /account/{id}/user
    regexp.MustCompile(`^tenant/account/blocked$`),                                   // /tenant/account/blocked
}

func TestValidator(t *testing.T) {
    var cases = []struct {
        path      string
        expection bool
    }{
        {"company", true},
        {"tenant/sj3co3s4", false},
        {"company/sd45f768", true},
        {"account/acc74850", true},
        {"company/account", true},
        {"acc734340", true},
        {"account/acc234234/user", true},
        {"account/blocked", false},
        {"tenant/account/blocked", true},
        {"tenant/account/acc23849", false},

        {"company/something?foo=bar", false},
        {"company?foo=bar", true},
        {"tenant/account/blocked?foo=bar&x=z", true},
        {"acc23849?xyz", true},
        {"company?", true},
    }

    proxy := proxy{
        downstream: &downstream{
            allowedList: allowedList,
        },
    }

    for _, tc := range cases {
        require.Equal(t, tc.expection, proxy.ValidatePath(tc.path), "Test is failing! %d equal %s", tc.expection, tc.path)
    }
}

func TestProxyHandler(t *testing.T) {
    var cases = map[string]struct {
        method string
        path   string
        status int
        body   string
        cookie *http.Cookie
    }{
        "/company":          {"GET", "/company", http.StatusOK, "I am the backend", &http.Cookie{Name: "someName", Value: "SomeValue"}},
        "/company/sd45f768": {"POST", "/company/sd45f768", http.StatusOK, "Done", nil},
        "/company/abc85033": {"PUT", "/company/abc85033", http.StatusBadRequest, "Bad Request!", nil},
        "/account/acc74850": {"GET", "/account/acc74850", http.StatusNotFound, PATH_NOT_FOUND, nil},
        "/account/acc74850?foo=bar": {"GET", "/account/acc74850?foo=bar&x=y", http.StatusNotFound, PATH_NOT_FOUND, nil},
    }

    // downstream server
    backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tCase := cases[r.RequestURI]

        require.Equal(t, tCase.path, r.URL.Path)
        require.Equal(t, tCase.method, r.Method)
        require.NotEmpty(t, r.Header.Get("X-Forwarded-For"))

        if tCase.cookie != nil {
            http.SetCookie(w, tCase.cookie)
        }
        w.WriteHeader(tCase.status)
        w.Write([]byte(tCase.body))
    }))
    defer backend.Close()

    prx := NewProxy(Config{
        DownstreamUrl: backend.URL,
        Paths: struct {
            Params  map[string]string `yaml:"params"`
            Allowed []string          `yaml:"allowed"`
        }{
            Params: map[string]string{"id": "[a-z]+[0-9]+[a-z0-9]+"},
            Allowed: []string{
                "company",
                "company/{id}",
            },
        },
    })

    // sidecar server
    sidecar := httptest.NewServer(prx)
    defer sidecar.Close()
    sidecarClient := sidecar.Client()

    for _, tCase := range cases {
        getReq, _ := http.NewRequest(tCase.method, sidecar.URL+tCase.path, nil)
        getReq.Close = true
        res, err := sidecarClient.Do(getReq)
        if err != nil {
            t.Fatalf("Get: %v", err)
        }
        require.Equal(t, tCase.status, res.StatusCode)
        body, err := ioutil.ReadAll(res.Body)
        res.Body.Close()
        if err != nil {
            t.Fatalf("reading body: %v", err)
        }
        require.Equal(t, tCase.body, string(body))
        if tCase.cookie != nil {
            require.Equal(t, tCase.cookie.Name, res.Cookies()[0].Name)
            require.Equal(t, tCase.cookie.Value, res.Cookies()[0].Value)
        } else {
            require.Equal(t, 0, len(res.Header["Set-Cookie"])) // len(res.Cookies())
        }
    }
}
