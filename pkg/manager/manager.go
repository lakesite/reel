package manager

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"

	"github.com/lakesite/ls-fibre/pkg/service"
)

// ManagerService has a toml Config property which contains reel specific directives,
// a DBConfig array of app database configurations, and a pointer to the web
// service.
type ManagerService struct {
	Config     *toml.Tree
	DBConfig   map[string]*DBConfig
	WebService *service.WebService
}

// Rewind takes an app's driver configuration and runs the appropriate command
// to rewind the database
func (ms *ManagerService) Rewind(app string, source string) {
	if source == "" {
		// we need to get the default source
		source = ms.DBConfig[app].Source
	}

	// case for driver
	switch ms.DBConfig[app].Driver {
	case "mysql":
		ms.RewindMysql(app, source)
	case "postgres":
		ms.RewindPostgres(app, source)
	default:
		log.Fatalf("Unknown/unsupported database driver: %s\n", ms.DBConfig[app].Driver)
	}
	fmt.Println("reel OK")
}

// GetSources uses the app's configuration dbsources to list available sources
func (ms *ManagerService) GetSources(app string) ([]string, error) {
	var sources []string

	err := filepath.Walk(ms.DBConfig[app].Sources, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			sources = append(sources, path)
		}
		return nil
	})
	return sources, err
}

// PrintSources displays a list of sources via GetSources, using app
// configuration.
func (ms *ManagerService) PrintSources(app string) {
	sources, err := ms.GetSources(app)
	if err != nil {
		log.Printf("Error getting dbsources: %s", err)
	}
	for _, file := range sources {
		fmt.Println(file)
	}
}

// InitApp initializes an app configuration, return true if successful false
// otherwise
func (ms *ManagerService) InitApp(app string) bool {
	// should we re-init every time?
	if ms.DBConfig[app] == nil {
		ms.DBConfig[app] = &DBConfig{}
	}

	err := fmt.Errorf("")
	success := true

	// pull in the database config to DBConfig struct
	ms.DBConfig[app].Server, err = ms.GetAppProperty(app, "dbserver")
	if err != nil {
		log.Printf("InitApp Error: %s\n", err)
		success = false
	}

	ms.DBConfig[app].Port, err = ms.GetAppProperty(app, "dbport")
	if err != nil {
		log.Printf("InitApp Error: %s\n", err)
		success = false
	}

	ms.DBConfig[app].Database, err = ms.GetAppProperty(app, "database")
	if err != nil {
		log.Printf("InitApp Error: %s\n", err)
		success = false
	}

	ms.DBConfig[app].User, err = ms.GetAppProperty(app, "dbuser")
	if err != nil {
		log.Printf("InitApp Error: %s\n", err)
		success = false
	}

	ms.DBConfig[app].Password, err = ms.GetAppProperty(app, "dbpassword")
	if err != nil {
		log.Printf("InitApp Error: %s\n", err)
		success = false
	}

	ms.DBConfig[app].Driver, err = ms.GetAppProperty(app, "dbdriver")
	if err != nil {
		log.Printf("InitApp Error: %s\n", err)
		success = false
	}

	ms.DBConfig[app].Sources, err = ms.GetAppProperty(app, "dbsources")
	if err == nil {
		if _, err = os.Stat(ms.DBConfig[app].Sources); os.IsNotExist(err) {
			log.Printf("InitApp Error: Source databases dumpfile directory '%s' does not exist.\n", ms.DBConfig[app].Sources)
			success = false
		}
	}

	// if we don't have dbsources defined, we must have a single source defined:
	ms.DBConfig[app].Source, err = ms.GetAppProperty(app, "dbsource")
	if err != nil && ms.DBConfig[app].Sources == "" {
		log.Printf("InitApp Error: %s\n", err)
		success = false
	} else {
		if _, err = os.Stat(ms.DBConfig[app].Source); os.IsNotExist(err) {
			log.Printf("InitApp Error: Source database dumpfile '%s' does not exist.\n", ms.DBConfig[app].Source)
			success = false
		}
	}

	return success
}

// Init is required to initialize the manager service via a config file.
func (ms *ManagerService) Init(cfgfile string) {
	if _, err := os.Stat(cfgfile); os.IsNotExist(err) {
		log.Fatalf("File '%s' does not exist.\n", cfgfile)
	} else {
		ms.Config, _ = toml.LoadFile(cfgfile)
		ms.DBConfig = make(map[string]*DBConfig)
	}
}
