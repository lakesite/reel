package manager

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pelletier/go-toml"

	"github.com/lakesite/ls-config/pkg/config"
	"github.com/lakesite/ls-fibre/pkg/service"
)

type DBConfig struct {
	Server string
	Port string
	Database string
	User string
	Password string
	Driver string
	Source string
}

func RewindMysql(dbconfig *DBConfig) {
	// clear database:
	db, err := sql.Open("mysql", dbconfig.User + ":" + dbconfig.Password + "@/" + dbconfig.Database)
	defer db.Close()

	_, err = db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	if err != nil {
		log.Fatalf("Error turning foreign key checks off: %s\n", err)
	}

	query := "SELECT concat('DROP TABLE IF EXISTS ', table_name, ';') FROM information_schema.tables WHERE table_schema = ?"
	rows, err := db.Query(query, dbconfig.Database)
	defer rows.Close()
	if err != nil {
		log.Fatalf("Error querying database %s: %s\n", dbconfig.Database, err)
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
	cmd := exec.Command("mysql", "-u", dbconfig.User, "-p" + dbconfig.Password, dbconfig.Database, "-e", "source " + dbconfig.Source)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("mysql import failed with status: %s\n", err)
	}
}

// func RewindPostgres(dbconfig DBConfig) {
// }

func Rewind(cfgfile string, app string) {
	if _, err := os.Stat(cfgfile); os.IsNotExist(err) {
		log.Fatalf("File '%s' does not exist.\n", cfgfile)
	} else {
		config, _ := toml.LoadFile(cfgfile)
		dbconfig := &DBConfig{}

		// pull in the database config to DBConfig struct
		if config.Get(app + ".dbserver") != nil {
			dbconfig.Server =	config.Get(app + ".dbserver").(string)
		} else {
			log.Fatalf("Configuration missing dbserver section under [%s] heading.\n", app)
		}

		if config.Get(app + ".dbport") != nil {
			dbconfig.Port =	config.Get(app + ".dbport").(string)
		} else {
			log.Fatalf("Configuration missing dbport section under [%s] heading.\n", app)
		}

		if config.Get(app + ".database") != nil {
			dbconfig.Database = config.Get(app + ".database").(string)
		} else {
			log.Fatalf("Configuration missing database section under [%s] heading.\n", app)
		}

		if config.Get(app + ".dbuser") != nil {
			dbconfig.User = config.Get(app + ".dbuser").(string)
		} else {
			log.Fatalf("Configuration missing dbuser section under [%s] heading.\n", app)
		}

		if config.Get(app + ".dbpassword") != nil {
			dbconfig.Password = config.Get(app + ".dbpassword").(string)
		} else {
			log.Fatalf("Configuration missing dbpassword section under [%s] heading.\n", app)
		}

		if config.Get(app + ".dbdriver") != nil {
			dbconfig.Driver = config.Get(app + ".dbdriver").(string)
		} else {
			log.Fatalf("Configuration missing dbdriver section under [%s] heading.\n", app)
		}

		if config.Get(app + ".dbsource") != nil {
			dbconfig.Source = config.Get(app + ".dbsource").(string)
		} else {
			log.Fatalf("Configuration missing dbsource section under [%s] heading.\n", app)
		}

		if _, err = os.Stat(dbconfig.Source); os.IsNotExist(err) {
			log.Fatalf("Source database dumpfile '%s' does not exist.\n", dbconfig.Source)
		}

		// case for driver
		switch dbconfig.Driver {
			case "mysql":
				RewindMysql(dbconfig)
			//case "postgres":
			//	RewindPostgres(dbconfig)
			default:
				log.Fatalf("Unknown/unsupported database driver: %s\n", dbconfig.Driver)
		}
		fmt.Println("reel OK")
	}
}

func RunManagementService() {
	address := config.Getenv("REEL_HOST", "127.0.0.1") + ":" + config.Getenv("REEL_PORT", "7999")
	ws := service.NewWebService("reel", address)
	// we need to add handlers to rewind, etc.
	ws.RunWebServer()
}
