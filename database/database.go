package database

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
	"errors"
)

const (
	initTables = "CREATE TABLE IF NOT EXISTS `ManagedRole` (`RoleID` TEXT UNIQUE, `Name` TEXT NOT NULL,`Description`TEXT);"
)

var db *sql.DB

//Initializes the database. Run at startup
func Init() {
	tdb, err := sql.Open("sqlite3", "./botdata.db")
	if err != nil {
		log.Fatal(err)
	}
	db = tdb
	//defer db.Close()
	_, err = db.Exec(initTables)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Done initializing database")
}

//Defer call this right after calling init so that we close this connection gracefully
func End() {
	db.Close()
}

func GetManagedRoles(userID string) ([]ManagedRole, error) {
	//TODO: only pull roles that the user has permissions to join
	if db == nil {
		tdb, err := sql.Open("sqlite3", "./botdata.db")
		if err != nil {
			log.Fatal(err)
		}
		db = tdb
	}
	rows, err := db.Query("SELECT RoleID, Name, Description FROM MangedRoles")
	if err != nil {
		log.Print("Unable to retrieve all Managed roles")
		return nil, errors.New("Unable to retrieve managed roles")
	}
	output := make([]ManagedRole, 0)
	defer rows.Close()
	for rows.Next() {
		var RoleID string
		var Name string
		var Description string
		err = rows.Scan(&RoleID, &Name, &Description)
		if err != nil {
			log.Print("Unable to retrieve all Managed roles")
			continue
		}
		output = append(output, ManagedRole{RoleID: RoleID, Name:Name, Description: Description})
	}
	return output, nil
}


type ManagedRole struct {
	RoleID      string
	Name        string
	Description string
}


