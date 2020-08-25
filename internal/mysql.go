package internal

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	// Need this grab side effects to use sql
	_ "github.com/go-sql-driver/mysql"

	"github.com/fatih/color"
)

// Db is the exported variable for the connection of the database
var Db *sql.DB

func init() {

	Db, err := sql.Open("mysql", "root:Passw0rd!@tcp(127.0.0.1:3307)/")
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("[-] Unable to connect to MySQL instance, exiting...")
		os.Exit(1)
	} else {
		fmt.Println(color.GreenString("[+]"), "Connected to MySQL instance successfully")
	}

	_, err = Db.Exec("CREATE DATABASE Anti")
	if strings.Contains(err.Error(), "database exists") {
		fmt.Println("[*] Database already exists, skipping")
	} else {
		fmt.Println("[+] Created new database")
	}

	_, err = Db.Exec("USE Anti")
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("[-] Unable to use database, exiting...")
		os.Exit(1)
	} else {
		fmt.Println(color.GreenString("[+]"), "DB selected Successfully")
	}

	fmt.Println("[*] Creating tables now...")

	time.Sleep(1 * time.Second)

	// **************************************************************************************************************************************************************************************************************************************
	// Pictures Table
	stmtPic, err := Db.Prepare("CREATE TABLE Pictures (ID INT AUTO_INCREMENT PRIMARY KEY,  command VARCHAR(1000), baseimage VARCHAR(255), new_filename VARCHAR(255), album_id VARCHAR(255), album_deletehash VARCHAR(255));")
	if err != nil {
		fmt.Println(err)
	}
	defer stmtPic.Close()

	_, err = stmtPic.Exec()

	if err != nil {
		if strings.Contains(err.Error(), "Error 1050") {
			fmt.Println(color.YellowString("[!]"), "Pictures table already exists, skipping...")
		}
	} else {
		fmt.Println(color.GreenString("[+]"), "Succesfully created Pictures table")
	}
	// **************************************************************************************************************************************************************************************************************************************

	time.Sleep(1 * time.Second)

	// **************************************************************************************************************************************************************************************************************************************
	// Albums table
	stmtAlbum, err := Db.Prepare("CREATE TABLE Albums (ID INT AUTO_INCREMENT PRIMARY KEY, Album_Hash VARCHAR(255), Delete_Hash VARCHAR(255), Auth_Type VARCHAR(255), Token VARCHAR(255));")
	if err != nil {
		fmt.Println(err)
	}
	_, err = stmtAlbum.Exec()

	if err != nil {
		if strings.Contains(err.Error(), "Error 1050") {
			fmt.Println(color.YellowString("[!]"), "Albums table already exists, skipping...")
		}
	} else {
		fmt.Println(color.GreenString("[+]"), "Succesfully created Albums table")
	}
	// **************************************************************************************************************************************************************************************************************************************

	time.Sleep(1 * time.Second)

	// **************************************************************************************************************************************************************************************************************************************
	// Tasking table

	// Removed: Response TEXT, Agent VARCHAR(255), Image_Hash VARCHAR(255),
	stmtTask, err := Db.Prepare("CREATE TABLE Tasking (Tasking_Image VARCHAR(255), Tasking_Command VARCHAR(255), Title VARCHAR(255), Tags VARCHAR(255), Delete_Hash VARCHAR(255));")
	if err != nil {
		fmt.Println(err)
	}
	_, err = stmtTask.Exec()
	if err != nil {
		if strings.Contains(err.Error(), "Error 1050") {
			fmt.Println(color.YellowString("[!]"), "Tasking table already exists, skipping...")
		}
	} else {
		fmt.Println(color.GreenString("[+]"), "Succesfully created Tasking table")
	}

	// **************************************************************************************************************************************************************************************************************************************

	time.Sleep(1 * time.Second)

	// **************************************************************************************************************************************************************************************************************************************
	// Agent table
	stmtAgent, err := Db.Prepare("CREATE TABLE Agents (ID INT AUTO_INCREMENT PRIMARY KEY, Status VARCHAR(255), Title VARCHAR(255), Tags VARCHAR(255));")
	if err != nil {
		fmt.Println(err)
	}
	_, err = stmtAgent.Exec()

	if err != nil {
		if strings.Contains(err.Error(), "Error 1050") {
			fmt.Println(color.YellowString("[!]"), "Agents table already exists, skipping...")
		}
	} else {
		fmt.Println(color.GreenString("[+]"), "Succesfully created Agents table")
	}
	// **************************************************************************************************************************************************************************************************************************************

	fmt.Println(color.GreenString("[+]"), "Successfully created all the tables")
}

// InsertClientID is a function that just puts the client-id into the table for others to grab
func InsertClientID(clientID string) {

	// Had the hardest time with this, but forgot to load this connection routine at the beginning of this function
	Db, err := sql.Open("mysql", "root:Passw0rd!@tcp(127.0.0.1:3307)/")
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = Db.Exec("USE Anti")
	if err != nil {
		fmt.Println(err.Error())
	}

	insert, err := Db.Prepare("INSERT IGNORE INTO Pictures (clientID) VALUES(?)")
	if err != nil {
		fmt.Println(err)
	}
	insert.Exec(clientID)

}

// GetClientID is a function that just puts the client-id into the table for others to grab
func GetClientID() (clientid string) {

	// Had the hardest time with this, but forgot to load this connection routine at the beginning of this function
	Db, err := sql.Open("mysql", "root:Passw0rd!@tcp(127.0.0.1:3307)/")
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = Db.Exec("USE Anti")
	if err != nil {
		fmt.Println(err.Error())
	}

	blah, _ := Db.Query("SELECT clientID FROM Pictures WHERE ID=1")
	defer blah.Close()
	for blah.Next() {
		err := blah.Scan(&clientid)
		if err != nil {
			fmt.Println(err)
		}

	}
	return clientid
}

// InsertImages is a function to help insert image options into the database
func InsertImages(clientid, command, baseimage, newfilename string) {

	// Had the hardest time with this, but forgot to load this connection routine at the beginning of this function
	Db, err := sql.Open("mysql", "root:Passw0rd!@tcp(127.0.0.1:3307)/")
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = Db.Exec("USE Anti")
	if err != nil {
		fmt.Println(err.Error())
	}

	insert, err := Db.Prepare("INSERT IGNORE INTO Pictures (clientID, command, baseimage, new_filename) VALUES ( ?, ?, ?, ? )")
	if err != nil {
		fmt.Println(err)
	}
	defer insert.Close()

	insert.Exec(clientid, command, baseimage, newfilename)
}

// InsertAlbum is a function that will insert the data from the album module into the database
func InsertAlbum(title, albumid, deletehash string) {

	// Had the hardest time with this, but forgot to load this connection routine at the beginning of this function
	Db, err := sql.Open("mysql", "root:Passw0rd!@tcp(127.0.0.1:3307)/")
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = Db.Exec("USE Anti")
	if err != nil {
		fmt.Println(err)
	}

	insert, err := Db.Prepare("INSERT INTO Albums (title, AlbumID, DeleteHash) VALUES ( ?, ?, ? )")
	if err != nil {
		fmt.Println(err)
	}
	defer insert.Close()

	insert.Exec(title, albumid, deletehash)
}

// InsertTask is a function that will insert the data from the task module into the database
func InsertTask(taskingimage, title, description, imageID, deleteHash string) {

	// Had the hardest time with this, but forgot to load this connection routine at the beginning of this function
	Db, err := sql.Open("mysql", "root:Passw0rd!@tcp(127.0.0.1:3307)/")
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = Db.Exec("USE Anti")
	if err != nil {
		fmt.Println(err)
	}

	insert, err := Db.Prepare("INSERT INTO Tasking (Tasking_Image, Title, Description, Image_Hash, Delete_Hash) VALUES ( ?, ?, ?, ?, ? )")
	if err != nil {
		fmt.Println(err)
	}
	defer insert.Close()

	insert.Exec(taskingimage, title, description, imageID, deleteHash)

}

// GetImages is a function that will query the database for the current images that are encoded
func GetImages() {

	var (
		id          int
		newFilename string
		command     string
	)

	// Had the hardest time with this, but forgot to load this connection routine at the beginning of this function
	Db, err := sql.Open("mysql", "root:Passw0rd!@tcp(127.0.0.1:3307)/")
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = Db.Exec("USE Anti")
	if err != nil {
		fmt.Println(err)
	}

	rows, err := Db.Query("SELECT ID, new_filename, command FROM Pictures")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &newFilename, &command)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("|", color.GreenString("ID:"), id, "|", color.GreenString("File Name:"), newFilename, "|", color.GreenString("Encoded Command:"), command, "|")
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
	}

}

// GetAlbums is a function that will query the database for the current albums that have been created
func GetAlbums() {

	var (
		id        int
		title     string
		albumID   string
		albumHash string
	)

	// Had the hardest time with this, but forgot to load this connection routine at the beginning of this function
	Db, err := sql.Open("mysql", "root:Passw0rd!@tcp(127.0.0.1:3307)/")
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = Db.Exec("USE Anti")
	if err != nil {
		fmt.Println(err)
	}

	rows, err := Db.Query("SELECT ID, title, AlbumID, DeleteHash FROM Albums")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &title, &albumID, &albumHash)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("|", color.GreenString("ID:"), id, "|", color.GreenString("Album Title:"), title, "|", color.GreenString("Album ID:"), albumID, "|", color.GreenString("Album Delete Hash:"), albumHash, "|")

	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
	}

}
