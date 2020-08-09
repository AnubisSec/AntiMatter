package internal

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	// Need this grab side effects to use sql
	_ "github.com/go-sql-driver/mysql"
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
		fmt.Println("[+] Connected to MySQL instance successfully")
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
		fmt.Println("[+] DB selected Successfully")
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
	if strings.Contains(err.Error(), "Error 1050") {
		fmt.Println("[!] Pictures table already exists, skipping...")
	}
	// **************************************************************************************************************************************************************************************************************************************

	time.Sleep(1 * time.Second)

	// **************************************************************************************************************************************************************************************************************************************
	// Albums table
	stmtAlbum, err := Db.Prepare("CREATE TABLE Albums (ID INT AUTO_INCREMENT PRIMARY KEY, Album_Hash VARCHAR(255), Delete_Hash VARCHAR(255), Auth_Type VARCHAR(255), Token VARCHAR(255));")
	if err != nil {
		fmt.Println("Skipping...")
	}
	_, err = stmtAlbum.Exec()
	if strings.Contains(err.Error(), "Error 1050") {
		fmt.Println("[!] Albums table already exists, skipping...")
	}
	// **************************************************************************************************************************************************************************************************************************************

	time.Sleep(1 * time.Second)

	// **************************************************************************************************************************************************************************************************************************************
	// Tasking table
	stmtTask, err := Db.Prepare("CREATE TABLE Tasking (Tasking_Image VARCHAR(255), Tasking_Command VARCHAR(255), Response TEXT, Title VARCHAR(255), Tags VARCHAR(255), Agent VARCHAR(255), Image_Hash VARCHAR(255), Delete_Hash VARCHAR(255), Token VARCHAR(255));")
	if err != nil {
		fmt.Println("Skipping...")
	}
	_, err = stmtTask.Exec()
	if strings.Contains(err.Error(), "Error 1050") {
		fmt.Println("[!] Tasking table already exists, skipping...")
	}
	// **************************************************************************************************************************************************************************************************************************************

	time.Sleep(1 * time.Second)

	// **************************************************************************************************************************************************************************************************************************************
	// Agent table
	stmtAgent, err := Db.Prepare("CREATE TABLE Agents (ID INT AUTO_INCREMENT PRIMARY KEY, Status VARCHAR(255), Title VARCHAR(255), Tags VARCHAR(255));")
	if err != nil {
		fmt.Println("Skipping...")
	}
	_, err = stmtAgent.Exec()
	if strings.Contains(err.Error(), "Error 1050") {
		fmt.Println("[!] Agents table already exists, skipping...")
	}
	// **************************************************************************************************************************************************************************************************************************************

	fmt.Println("[+] Successfully created all the tables")

	//defer Db.Close()
}

// InsertImages is a function to help insert image options into the database
func InsertImages(command, baseimage, newfilename string) {

	// Had the hardest time with this, but forgot to load this connection routine at the beginning of this function
	Db, err := sql.Open("mysql", "root:Passw0rd!@tcp(127.0.0.1:3307)/")
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = Db.Exec("USE Anti")
	if err != nil {
		fmt.Println(err.Error())
	}

	insert, err := Db.Prepare("INSERT INTO Pictures (command, baseimage, new_filename) VALUES ( ?, ?, ? )")
	if err != nil {
		fmt.Println(err)
	}
	defer insert.Close()

	insert.Exec(command, baseimage, newfilename)
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
