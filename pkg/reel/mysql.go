package reel

import (
	"database/sql"
	"log"
	"os/exec"

	"github.com/lakesite/ls-governor"
	_ "github.com/go-sql-driver/mysql"
)

// RewindMysql takes the app instance name and uses the defined configuration
// to delete the tables associated with the database, then runs the mysql client
// to import the source database.
func RewindMysql(app string, source string, gapi *governor.API) {
	dbconfig := gapi.ManagerService.DBConfig[app]

	// clear database:
	db, err := sql.Open("mysql", dbconfig.User+":"+dbconfig.Password+"@/"+dbconfig.Database)
	if err != nil {
		log.Fatalf("mysql connection failed: %s\n", err)
	}
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
	cmd := exec.Command("mysql", "-u", dbconfig.User, "-p"+dbconfig.Password, dbconfig.Database, "-e", "source "+source)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("mysql import failed with status: %s\n", err)
	}
}
