// handlers contains the handlers to manage API endpoints
package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lakesite/ls-fibre"
	"github.com/gorilla/mux"
	"github.com/lakesite/ls-governor"

	"github.com/lakesite/reel/pkg/job"
	"github.com/lakesite/reel/pkg/reel"

)

// Handle requests to rewind an app via Rewind
func RewindHandler(w http.ResponseWriter, r *http.Request, gapi *governor.API) {
	var app string
	var source string
	var err error = nil

	vars := mux.Vars(r)

	if vars["app"] == "" {
		app, err = gapi.ManagerService.GetAppProperty("reel", "default_app")
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
		if reel.InitApp(app, gapi) {
			// submit a rewind job
			job.ReelQueue <- job.ReelJob{App: app, Source: source, Gapi: gapi}
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	}
}

// Handle requests to list available database sources
func SourcesHandler(w http.ResponseWriter, r *http.Request, gapi *governor.API) {
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
		app, err = gapi.ManagerService.GetAppProperty("reel", "default_app")
	}

	if err == nil {
		if reel.InitApp(app, gapi) {
			sources, err = reel.GetSources(app, gapi)
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

func ProxyAppHandler(w http.ResponseWriter, r *http.Request, gapi *governor.API) {
	vars := mux.Vars(r)
	if reel.InitApp(vars["app"], gapi) {
		proxyDest, err := gapi.ManagerService.GetAppProperty(vars["app"], "proxy_dest")
		if err == nil {
			// default proxy route
			cfg := []fibre.ProxyConfig{
				fibre.ProxyConfig{
					Path: "/" + vars["app"],
					Host: proxyDest,
					Override: fibre.ProxyOverride{
						Match: "/" + vars["app"],
						Path:  "/",
					},
				},
			}
			gapi.WebService.Proxy(cfg)
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	}
}
