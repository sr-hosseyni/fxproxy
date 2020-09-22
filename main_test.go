package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProxyHandler(t *testing.T) {
	var cases = map[string]struct {
		method string
		path   string
		status int
		body   string
		cookie *http.Cookie
	}{
		"/company": {"GET", "/company", http.StatusOK, "I am the backend", &http.Cookie{Name: "someName", Value: "SomeValue"}},
		"/company/sd45f768": {"POST", "/company/sd45f768", http.StatusOK, "Done", nil},
		"/company/abc85033": {"PUT", "/company/abc85033", http.StatusBadRequest, "Bad Request!", nil},
		"/account/acc74850": {"GET", "/account/acc74850", http.StatusNotFound, PATH_NOT_FOUND, nil},
	}

	// downstream server
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tCase := cases[r.RequestURI]

		require.Equal(t, tCase.method, r.Method)
		require.NotEmpty(t, r.Header.Get("X-Forwarded-For"))

		if tCase.cookie != nil {
			http.SetCookie(w, tCase.cookie)
		}
		w.WriteHeader(tCase.status)
		w.Write([]byte(tCase.body))
	}))
	defer backend.Close()

	proxy := &Proxy{
		downstream: DownStream{
			Address: backend.URL,
		},
		allowedList: []*regexp.Regexp{
			regexp.MustCompile(`^company$`),
			regexp.MustCompile(`^company/[a-z]+[0-9]+[a-z0-9]+$`),
		},
	}

	// sidecar server
	sidecar := httptest.NewServer(proxy)
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
