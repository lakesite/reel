// api.go
// This file defines all the web service handlers, and the web service API
// endpoints.

package manager

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/lakesite/ls-config/pkg/config"
	"github.com/lakesite/ls-fibre/pkg/service"
)

// Handle requests to rewind an app via Rewind
func (ms *ManagerService) RewindHandler(w http.ResponseWriter, r *http.Request) {
	var app string
	var source string
	var err error = nil

	vars := mux.Vars(r)

	if vars["app"] == "" {
		app, err = ms.GetAppProperty("reel", "default_app")
		if err != nil {
			// app is not defined and no default app is set:
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	} else {
		app = vars["app"]
	}

	source = ""
	if vars["source"] != "" {
		source = vars["source"]
	}

	if err == nil {
		if ms.InitApp(app) {
			ms.Rewind(app, source)
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	}
}

// Handle requests to list available database sources
func (ms *ManagerService) SourcesHandler(w http.ResponseWriter, r *http.Request) {
	var app string
	var sources []string
	var err error = nil
	type JR struct {
		Sources []string
	}

	// do we have an app specified?
	vars := mux.Vars(r)
	if vars["app"] != "" {
		app = vars["app"]
	} else {
		app, err = ms.GetAppProperty("reel", "default_app")
	}

	if err == nil {
		if ms.InitApp(app) {
			sources, err = ms.GetSources(app)
			if err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
			} else {
				w.WriteHeader(http.StatusOK)
				response, _ := json.MarshalIndent(JR{sources}, "", " ")
				fmt.Fprintf(w, string(response))
			}
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func (ms *ManagerService) ProxyAppHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if ms.InitApp(vars["app"]) {
		proxyDest, err := ms.GetAppProperty(vars["app"], "proxy_dest")
		if err == nil {
			// default proxy route
			cfg := []service.ProxyConfig{
				service.ProxyConfig{
					Path: "/" + vars["app"],
					Host: proxyDest,
					Override: service.ProxyOverride{
						Match: "/" + vars["app"],
						Path:  "/",
					},
				},
			}
			ms.WebService.Proxy(cfg)
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	}
}

// setupRoutes defines and associates routes to handlers.
func (ms *ManagerService) setupRoutes(ws *service.WebService) {
	ws.Router.HandleFunc("/reel/api/v1/sources/", ms.SourcesHandler)
	ws.Router.HandleFunc("/reel/api/v1/sources/{app}", ms.SourcesHandler)
	ws.Router.HandleFunc("/reel/api/v1/rewind/", ms.RewindHandler)
	ws.Router.HandleFunc("/reel/api/v1/rewind/{app}", ms.RewindHandler)
	ws.Router.HandleFunc("/reel/api/v1/rewind/{app}/{source}", ms.RewindHandler)
	ws.Router.HandleFunc("/reel/api/v1/proxy/{app}", ms.ProxyAppHandler)
}

// RunManagementService sets up the web service and defines routes for the API.
func (ms *ManagerService) RunManagementService() {
	address := config.Getenv("REEL_HOST", "127.0.0.1") + ":" + config.Getenv("REEL_PORT", "7999")
	ms.WebService = service.NewWebService("reel", address)
	ms.setupRoutes(ms.WebService)
	ms.WebService.RunWebServer()
}
