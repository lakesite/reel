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
	DBConfig *DBConfig
	Config *toml.Tree
}

func (ms *ManagerService) RewindMysql() {
	// clear database:
	db, err := sql.Open("mysql", ms.DBConfig.User+":"+ms.DBConfig.Password+"@/"+ms.DBConfig.Database)
	if err != nil {
		log.Fatalf("mysql connection failed: %s\n", err)
	}
	defer db.Close()

	_, err = db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	if err != nil {
		log.Fatalf("Error turning foreign key checks off: %s\n", err)
	}

	query := "SELECT concat('DROP TABLE IF EXISTS ', table_name, ';') FROM information_schema.tables WHERE table_schema = ?"
	rows, err := db.Query(query, ms.DBConfig.Database)
	defer rows.Close()
	if err != nil {
		log.Fatalf("Error querying database %s: %s\n", ms.DBConfig.Database, err)
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
	cmd := exec.Command("mysql", "-u", ms.DBConfig.User, "-p"+ms.DBConfig.Password, ms.DBConfig.Database, "-e", "source "+ms.DBConfig.Source)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("mysql import failed with status: %s\n", err)
	}
}

func (ms *ManagerService) RewindPostgres() {
	// clear database:
	connection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", ms.DBConfig.Server, ms.DBConfig.Port, ms.DBConfig.User, ms.DBConfig.Password, ms.DBConfig.Database)
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatalf("postgres connection failed: %s\n", err)
	}
	defer db.Close()

	_, err = db.Exec("DROP SCHEMA " + ms.DBConfig.Database + " CASCADE")
	if err != nil {
		log.Printf("Schema did not exist: %s\n", err)
	}
	_, err = db.Exec("CREATE SCHEMA " + ms.DBConfig.Database)
	if err != nil {
		log.Printf("Schema creation failed: %s\n", err)
	}

	// required for postgresql 9.3+
	_, err = db.Exec("GRANT ALL ON SCHEMA " + ms.DBConfig.Database + " to postgres")
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}
	_, err = db.Exec("GRANT ALL ON SCHEMA " + ms.DBConfig.Database + " to " + ms.DBConfig.Database)
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}
	_, err = db.Exec("COMMENT ON SCHEMA " + ms.DBConfig.Database + " IS 'standard " + ms.DBConfig.Database + " schema'")
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}

	// use psql client for import (dependency):
	cmd := exec.Command("psql", "-h", ms.DBConfig.Server, "-U", ms.DBConfig.User, "-d", ms.DBConfig.Database, "-f", ms.DBConfig.Source)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "PGPASSWORD=" + ms.DBConfig.Password)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("pssql import failed with status: %s\n", err)
	}
}

func (ms *ManagerService) Rewind() {
	// case for driver
	switch ms.DBConfig.Driver {
		case "mysql":
			ms.RewindMysql()
		case "postgres":
			ms.RewindPostgres()
		default:
			log.Fatalf("Unknown/unsupported database driver: %s\n", ms.DBConfig.Driver)
	}
	fmt.Println("reel OK")
}

func (ms *ManagerService) Init(cfgfile string, app string) {
	if _, err := os.Stat(cfgfile); os.IsNotExist(err) {
		log.Fatalf("File '%s' does not exist.\n", cfgfile)
	} else {
		ms.Config, _ = toml.LoadFile(cfgfile)
		ms.DBConfig = &DBConfig{}

		// pull in the database config to DBConfig struct
		if ms.Config.Get(app+".dbserver") != nil {
			ms.DBConfig.Server = ms.Config.Get(app + ".dbserver").(string)
		} else {
			log.Fatalf("Configuration missing dbserver section under [%s] heading.\n", app)
		}

		if ms.Config.Get(app+".dbport") != nil {
			ms.DBConfig.Port = ms.Config.Get(app + ".dbport").(string)
		} else {
			log.Fatalf("Configuration missing dbport section under [%s] heading.\n", app)
		}

		if ms.Config.Get(app+".database") != nil {
			ms.DBConfig.Database = ms.Config.Get(app + ".database").(string)
		} else {
			log.Fatalf("Configuration missing database section under [%s] heading.\n", app)
		}

		if ms.Config.Get(app+".dbuser") != nil {
			ms.DBConfig.User = ms.Config.Get(app + ".dbuser").(string)
		} else {
			log.Fatalf("Configuration missing dbuser section under [%s] heading.\n", app)
		}

		if ms.Config.Get(app+".dbpassword") != nil {
			ms.DBConfig.Password = ms.Config.Get(app + ".dbpassword").(string)
		} else {
			log.Fatalf("Configuration missing dbpassword section under [%s] heading.\n", app)
		}

		if ms.Config.Get(app+".dbdriver") != nil {
			ms.DBConfig.Driver = ms.Config.Get(app + ".dbdriver").(string)
		} else {
			log.Fatalf("Configuration missing dbdriver section under [%s] heading.\n", app)
		}

		if ms.Config.Get(app+".dbsource") != nil {
			ms.DBConfig.Source = ms.Config.Get(app + ".dbsource").(string)
		} else {
			log.Fatalf("Configuration missing dbsource section under [%s] heading.\n", app)
		}

		if _, err = os.Stat(ms.DBConfig.Source); os.IsNotExist(err) {
			log.Fatalf("Source database dumpfile '%s' does not exist.\n", ms.DBConfig.Source)
		}
	}
}

func (ms *ManagerService) RewindHandler(w http.ResponseWriter, r *http.Request) {
	// check tokens
	// check app name
	// look in cwd for config or maintain in memory.
	w.WriteHeader(http.StatusNoContent)
}

func (ms *ManagerService) RunManagementService() {
	address := config.Getenv("REEL_HOST", "127.0.0.1") + ":" + config.Getenv("REEL_PORT", "7999")
	ws := service.NewWebService("reel", address)
	// we need to add handlers to rewind, etc.
	ws.Router.HandleFunc("/api/v1/rewind", ms.RewindHandler)
	ws.RunWebServer()
}
