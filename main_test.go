package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

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
	}

	proxy := Proxy{
		allowedList: allowedList,
	}

	for _, tc := range cases {
		require.Equal(t, tc.expection, proxy.ValidatePath(tc.path), "Test is failing!")
	}
}

func TestProxyHandler(t *testing.T) {
	var cases = map[string]struct {
		method string
		path   string
		status int
		body   string
	}{
		"/company": {"GET", "/company", http.StatusOK, "I am the backend"},
	}

	// downstream server
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tCase := cases[r.RequestURI]

		require.Equal(t, tCase.method, r.Method)
		require.NotEmpty(t, r.Header.Get("X-Forwarded-For"))
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
	}
}
