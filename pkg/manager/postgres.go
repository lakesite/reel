package manager

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"

	_ "github.com/lib/pq"
)

// RewindPostgres takes the app instance name and uses the defined configuration
// to delete the tables associated with the database, then runs the postgres
// client to import the source database.
func (ms *ManagerService) RewindPostgres(app string, source string) {
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
	cmd := exec.Command("psql", "-h", ms.DBConfig[app].Server, "-U", ms.DBConfig[app].User, "-d", ms.DBConfig[app].Database, "-f", source)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "PGPASSWORD="+ms.DBConfig[app].Password)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("pssql import failed with status: %s\n", err)
	}
}
