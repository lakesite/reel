package reel

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/lakesite/ls-governor"
)

// Rewind takes an app's driver configuration and runs the appropriate command
// to rewind the database
func Rewind(app string, source string, gapi *governor.API) {
	dbc := gapi.ManagerService.DBConfig[app]

	if source == "" {
		// we need to get the default source
		source = dbc.Meta["source"] // Source
	}

	// case for driver
	switch dbc.Driver {
	case "mysql":
		RewindMysql(app, source, gapi)
	case "postgres":
		RewindPostgres(app, source, gapi)
	default:
		log.Fatalf("Unknown/unsupported database driver: %s\n", dbc.Driver)
	}
	fmt.Println("reel OK")
}

// GetSources uses the app's configuration dbsources to list available sources
func GetSources(app string, gapi *governor.API) ([]string, error) {
	var sources []string
	dbc := gapi.ManagerService.DBConfig[app]

	err := filepath.Walk(dbc.Meta["sources"]/*Sources*/, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			sources = append(sources, path)
		}
		return nil
	})
	return sources, err
}

// PrintSources displays a list of sources via GetSources, using app
// configuration.
func PrintSources(app string, gapi *governor.API) {
	sources, err := GetSources(app, gapi)
	if err != nil {
		log.Printf("Error getting dbsources: %s", err)
	}
	for _, file := range sources {
		fmt.Println(file)
	}
}

// InitApp initializes an app configuration, return true if successful false
// otherwise
func InitApp(app string, gapi *governor.API) bool {
	dbc := gapi.ManagerService.DBConfig[app]
	ms  := gapi.ManagerService
	if dbc == nil {
		gapi.ManagerService.InitDatastore(app)
		dbc = gapi.ManagerService.DBConfig[app]
	}

	err := fmt.Errorf("")
	success := true

	// pull in the database config to DBConfig struct
	dbc.Server, err = ms.GetAppProperty(app, "dbserver")
	if err != nil {
		log.Printf("InitApp Error: %s\n", err)
		success = false
	}

	dbc.Port, err = ms.GetAppProperty(app, "dbport")
	if err != nil {
		log.Printf("InitApp Error: %s\n", err)
		success = false
	}

	dbc.Database, err = ms.GetAppProperty(app, "database")
	if err != nil {
		log.Printf("InitApp Error: %s\n", err)
		success = false
	}

	dbc.User, err = ms.GetAppProperty(app, "dbuser")
	if err != nil {
		log.Printf("InitApp Error: %s\n", err)
		success = false
	}

	dbc.Password, err = ms.GetAppProperty(app, "dbpassword")
	if err != nil {
		log.Printf("InitApp Error: %s\n", err)
		success = false
	}

	dbc.Driver, err = ms.GetAppProperty(app, "dbdriver")
	if err != nil {
		log.Printf("InitApp Error: %s\n", err)
		success = false
	}

	dbc.Meta["sources"], err = ms.GetAppProperty(app, "dbsources")
	if err == nil {
		if _, err = os.Stat(dbc.Meta["sources"]); os.IsNotExist(err) {
			log.Printf("InitApp Error: Source databases dumpfile directory '%s' does not exist.\n", dbc.Meta["sources"])
			success = false
		}
	}

	// if we don't have dbsources defined, we must have a single source defined:
	dbc.Meta["source"]/*Source*/, err = ms.GetAppProperty(app, "dbsource")
	if err != nil && dbc.Meta["sources"] == "" {
		log.Printf("InitApp Error: %s\n", err)
		success = false
	} else {
		if _, err = os.Stat(dbc.Meta["source"]); os.IsNotExist(err) {
			log.Printf("InitApp Error: Source database dumpfile '%s' does not exist.\n", dbc.Meta["source"])
			success = false
		}
	}

	return success
}
