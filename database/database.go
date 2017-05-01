package database

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
	"git.chiefnoah.tech/chiefnoah/GoGetMeHentai/config"
	"github.com/bwmarrin/discordgo"
	"strings"
)

const (
	initTables = "CREATE TABLE IF NOT EXISTS `ManagedRole` (`RoleID` TEXT UNIQUE, `Name` TEXT NOT NULL,`Description`TEXT);" +
		"CREATE TABLE IF NOT EXISTS `User` (`JoinedAt` TEXT, `Nick` TEXT, `username` TEXT, `Bot` INTEGER DEFAULT 0, `DiscordID` TEXT, `NickLastChanged` INTEGER DEFAULT 0, `Admin` INTEGER DEFAULT 0);" +
		"CREATE TABLE IF NOT EXISTS `ManagedRolePermissions` (`ManagedRoleID` TEXT NOT NULL, `DiscordRoleID` TEXT NOT NULL);" //This is a bridging table between discord's roles and the managed roles. It basically just allows us to limit which Discord roles are allowed to join a managed role
)

var db *sql.DB
var cfg *config.ConfigWrapper

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
	cfg = config.LoadConfig()
	log.Print("Done initializing database")
}

//Defer call this right after calling init so that we close this connection gracefully
func End() {
	db.Close()
}

/*
	Returns the roles a user has permission to join.
 */
func GetManagedRoles(m *discordgo.Member) ([]ManagedRole, error) {
	//TODO: only pull roles that the user has permissions to join
	if db == nil {
		tdb, err := sql.Open("sqlite3", "./botdata.db")
		if err != nil {
			log.Fatal(err)
		}
		db = tdb
	}
	log.Printf("Checking if current roles: %+v is in the list of allowed managed roles...", m.Roles)
	rolesParams := strings.Join(m.Roles, ", ")
	rows, err := db.Query("SELECT RoleID, Name, Description FROM ManagedRole INNTER JOIN ManagedRolePermissions ON RoleID = ManagedRolePermissions.ManagedRoleID WHERE DiscordRoleID IN (?)", rolesParams)
	if err != nil {
		log.Print("Unable to query database: " , err)
		return nil, err
	}
	defer rows.Close()
	output := make([]ManagedRole, 0)
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

func AddManagedRole(discordRoldID, managedRoleID string, guildRoles []*discordgo.Role) error {


	return nil
}

func createManagedRole() {

}

type ManagedRole struct {
	RoleID      string
	Name        string
	Description string
}


