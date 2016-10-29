package database

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
	"errors"
)

const (
	initTables = "CREATE TABLE IF NOT EXISTS `ManagedRole` (`RoleID` TEXT UNIQUE, `Name` TEXT NOT NULL,`Description`TEXT);" +
		"CREATE TABLE IF NOT EXISTS `User` (`JoinedAt` TEXT, `Nick` TEXT, `username` TEXT, `Bot` INTEGER DEFAULT 0, `DiscordID` TEXT, `NickLastChanged` INTEGER DEFAULT 0, `Admin` INTEGER DEFAULT 0);" +
		"CREATE TABLE IF NOT EXISTS `ManagedRolePermissions` (`ManagedRoleID` TEXT NOT NULL, `DiscordRole` TEXT NOT NULL);" //This is a bridging table between discord's roles and the managed roles. It basically just allows us to limit which Discord roles are allowed to join a managed role
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

/*
	Returns the roles a user has permission to join.
 */
func GetManagedRoles(userID string) ([]ManagedRole, error) {
	//TODO: only pull roles that the user has permissions to join
	if db == nil {
		tdb, err := sql.Open("sqlite3", "./botdata.db")
		if err != nil {
			log.Fatal(err)
		}
		db = tdb
	}
	rows, err := db.Query("SELECT RoleID, Name, Description FROM ManagedRole INNTER JOIN ManagedRolePermissions ON RoleID = ManagedRolePermissions.ManagedRoleID")
	if err != nil {
		log.Printf("Unable to retrieve all Managed roles for userID: %s\nError:%s",userID, err)
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


