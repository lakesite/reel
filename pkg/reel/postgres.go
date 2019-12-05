package reel

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/lakesite/ls-governor"
	_ "github.com/lib/pq"
)

// RewindPostgres takes the app instance name and uses the defined configuration
// to delete the tables associated with the database, then runs the postgres
// client to import the source database.
func RewindPostgres(app string, source string, gapi *governor.API) {
	dbconfig := gapi.ManagerService.DBConfig[app]

	// clear database:
	connection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbconfig.Server, dbconfig.Port, dbconfig.User, dbconfig.Password, dbconfig.Database)
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatalf("postgres connection failed: %s\n", err)
	}
	defer db.Close()

	_, err = db.Exec("DROP SCHEMA " + dbconfig.Database + " CASCADE")
	if err != nil {
		log.Printf("Schema did not exist: %s\n", err)
	}
	_, err = db.Exec("CREATE SCHEMA " + dbconfig.Database)
	if err != nil {
		log.Printf("Schema creation failed: %s\n", err)
	}

	// required for postgresql 9.3+
	_, err = db.Exec("GRANT ALL ON SCHEMA " + dbconfig.Database + " to postgres")
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}
	_, err = db.Exec("GRANT ALL ON SCHEMA " + dbconfig.Database + " to " + dbconfig.Database)
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}
	_, err = db.Exec("COMMENT ON SCHEMA " + dbconfig.Database + " IS 'standard " + dbconfig.Database + " schema'")
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}

	// use psql client for import (dependency):
	cmd := exec.Command("psql", "-h", dbconfig.Server, "-U", dbconfig.User, "-d", dbconfig.Database, "-f", source)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "PGPASSWORD="+dbconfig.Password)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("pssql import failed with status: %s\n", err)
	}
}