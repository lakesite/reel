package manager

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/pelletier/go-toml"
	"github.com/gorilla/mux"

	"github.com/lakesite/ls-config/pkg/config"
	"github.com/lakesite/ls-fibre/pkg/service"
)

type DBConfig struct {
	Server   string
	Port     string
	Database string
	User     string
	Password string
	Driver   string
	Source   string
}

type ManagerService struct {
	Config *toml.Tree
	DBConfig map[string]*DBConfig
}

func (ms *ManagerService) RewindMysql(app string) {
	// clear database:
	db, err := sql.Open("mysql", ms.DBConfig[app].User+":"+ms.DBConfig[app].Password+"@/"+ms.DBConfig[app].Database)
	if err != nil {
		log.Fatalf("mysql connection failed: %s\n", err)
	}
	defer db.Close()

	_, err = db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	if err != nil {
		log.Fatalf("Error turning foreign key checks off: %s\n", err)
	}

	query := "SELECT concat('DROP TABLE IF EXISTS ', table_name, ';') FROM information_schema.tables WHERE table_schema = ?"
	rows, err := db.Query(query, ms.DBConfig[app].Database)
	defer rows.Close()
	if err != nil {
		log.Fatalf("Error querying database %s: %s\n", ms.DBConfig[app].Database, err)
	}

	result := ""
	for rows.Next() {
		err := rows.Scan(&result)
		if err != nil {
			log.Fatalf("Error reading result set: %s\n", err)
		}
		_, err = db.Exec(result)
	}

	_, err = db.Exec("SET FOREIGN_KEY_CHECKS = 1")
	if err != nil {
		log.Fatalf("Error turning foreign key checks on: %s\n", err)
	}

	// use mysql client for import (dependency):
	cmd := exec.Command("mysql", "-u", ms.DBConfig[app].User, "-p"+ms.DBConfig[app].Password, ms.DBConfig[app].Database, "-e", "source "+ms.DBConfig[app].Source)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("mysql import failed with status: %s\n", err)
	}
}

func (ms *ManagerService) RewindPostgres(app string) {
	// clear database:
	connection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", ms.DBConfig[app].Server, ms.DBConfig[app].Port, ms.DBConfig[app].User, ms.DBConfig[app].Password, ms.DBConfig[app].Database)
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatalf("postgres connection failed: %s\n", err)
	}
	defer db.Close()

	_, err = db.Exec("DROP SCHEMA " + ms.DBConfig[app].Database + " CASCADE")
	if err != nil {
		log.Printf("Schema did not exist: %s\n", err)
	}
	_, err = db.Exec("CREATE SCHEMA " + ms.DBConfig[app].Database)
	if err != nil {
		log.Printf("Schema creation failed: %s\n", err)
	}

	// required for postgresql 9.3+
	_, err = db.Exec("GRANT ALL ON SCHEMA " + ms.DBConfig[app].Database + " to postgres")
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}
	_, err = db.Exec("GRANT ALL ON SCHEMA " + ms.DBConfig[app].Database + " to " + ms.DBConfig[app].Database)
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}
	_, err = db.Exec("COMMENT ON SCHEMA " + ms.DBConfig[app].Database + " IS 'standard " + ms.DBConfig[app].Database + " schema'")
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}

	// use psql client for import (dependency):
	cmd := exec.Command("psql", "-h", ms.DBConfig[app].Server, "-U", ms.DBConfig[app].User, "-d", ms.DBConfig[app].Database, "-f", ms.DBConfig[app].Source)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "PGPASSWORD=" + ms.DBConfig[app].Password)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("pssql import failed with status: %s\n", err)
	}
}

func (ms *ManagerService) Rewind(app string) {
	// case for driver
	switch ms.DBConfig[app].Driver {
		case "mysql":
			ms.RewindMysql(app)
		case "postgres":
			ms.RewindPostgres(app)
		default:
			log.Fatalf("Unknown/unsupported database driver: %s\n", ms.DBConfig[app].Driver)
	}
	fmt.Println("reel OK")
}

// get the property for app as a string, if property does not exist return err
func (ms *ManagerService) GetAppProperty(app string, property string) (string, error) {
	if ms.Config.Get(app+"."+property) != nil {
		return ms.Config.Get(app + "." + property).(string), nil
	} else {
		return "", fmt.Errorf("Configuration missing '%s' section under [%s] heading.\n", property, app)
	}
}

// initialize an app configuration, return true if successful false otherwise
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

	ms.DBConfig[app].Source, err = ms.GetAppProperty(app, "dbsource")
	if err != nil {
		log.Printf("InitApp Error: %s\n", err)
		success = false
	}

	if _, err = os.Stat(ms.DBConfig[app].Source); os.IsNotExist(err) {
		log.Printf("InitApp Error: Source database dumpfile '%s' does not exist.\n", ms.DBConfig[app].Source)
		success = false
	}

	return success
}

func (ms *ManagerService) Init(cfgfile string) {
	if _, err := os.Stat(cfgfile); os.IsNotExist(err) {
		log.Fatalf("File '%s' does not exist.\n", cfgfile)
	} else {
		ms.Config, _ = toml.LoadFile(cfgfile)
		ms.DBConfig = make(map[string]*DBConfig)
	}
}

func (ms *ManagerService) RewindHandler(w http.ResponseWriter, r *http.Request) {
	// check tokens
	// check app name
	// look in cwd for config or maintain in memory.
	vars := mux.Vars(r)
	if ms.InitApp(vars["app"]) {
		ms.Rewind(vars["app"])
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func (ms *ManagerService) RunManagementService() {
	address := config.Getenv("REEL_HOST", "127.0.0.1") + ":" + config.Getenv("REEL_PORT", "7999")
	ws := service.NewWebService("reel", address)
	// we need to add handlers to rewind, etc.
	ws.Router.HandleFunc("/api/v1/rewind/{app}", ms.RewindHandler)
	ws.RunWebServer()
}
