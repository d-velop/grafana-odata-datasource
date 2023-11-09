package plugin

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var (
	oc          ODataClientImpl
	handler     func(w http.ResponseWriter, r *http.Request)
	handlerPath string
)

func TestMain(m *testing.M) {
	fmt.Println("mocking server")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if handlerPath == "*" {
			handler(w, r)
			return
		}

		switch strings.TrimSpace(r.URL.Path) {
		case handlerPath:
			handler(w, r)
		default:
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	}))

	fmt.Println("mocking odata server")
	oc = ODataClientImpl{
		httpClient: &http.Client{
			Timeout: 1 * time.Second,
		},
		baseUrl: server.URL,
	}

	fmt.Println("run tests")
	m.Run()
}

func GetOC(hp string, h func(w http.ResponseWriter, r *http.Request)) ODataClient {
	handler = h
	handlerPath = hp
	return &oc
}
