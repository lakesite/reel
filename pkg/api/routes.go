// routes handles setting up routes for our API
package api

import (
	"net/http"

	"github.com/lakesite/ls-governor"
)

// SetupRoutes defines and associates routes to handlers.
// Use a wrapper convention to pass a governor API to each handler.
func SetupRoutes(gapi *governor.API) {
	gapi.WebService.Router.HandleFunc(
		"/reel/api/v1/sources/", 
		func(w http.ResponseWriter, r *http.Request) {
			SourcesHandler(w, r, gapi)
		},
	)
	gapi.WebService.Router.HandleFunc(
		"/reel/api/v1/sources/{app}",
		func(w http.ResponseWriter, r *http.Request) {
			SourcesHandler(w, r, gapi)
		},
	)
	gapi.WebService.Router.HandleFunc(
		"/reel/api/v1/rewind/",
		func(w http.ResponseWriter, r *http.Request) {
			RewindHandler(w, r, gapi)
		},
	)
	gapi.WebService.Router.HandleFunc(
		"/reel/api/v1/rewind/{app}",
		func(w http.ResponseWriter, r *http.Request) {
			RewindHandler(w, r, gapi)
		},
	)
	gapi.WebService.Router.HandleFunc(
		"/reel/api/v1/rewind/{app}/{source}",
		func(w http.ResponseWriter, r *http.Request) {
			RewindHandler(w, r, gapi)
		},
	)
	gapi.WebService.Router.HandleFunc(
		"/reel/api/v1/proxy/{app}",
		func(w http.ResponseWriter, r *http.Request) {
			ProxyAppHandler(w, r, gapi)
		},
	)
}