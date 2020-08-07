package internal

import (
	"database/sql"
	"fmt"

	// Need this grab side effects to use sql
	_ "github.com/go-sql-driver/mysql"
)

// Db is the exported variable for the connection of the database
var Db *sql.DB

func init() {

	Db, err := sql.Open("mysql", "root:Passw0rd!@tcp(127.0.0.1:3307)/")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("[+] Connected successfully")
	}

	_, err = Db.Exec("CREATE DATABASE Anti")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("[+] Created database")
	}

	_, err = Db.Exec("USE Anti")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("[+] DB selected Successfully")
	}

	fmt.Println("[!] Creating tables now...")

	// **************************************************************************************************************************************************************************************************************************************
	// Pictures Table
	stmtPic, err := Db.Prepare("CREATE TABLE Pictures (ID INT AUTO_INCREMENT PRIMARY KEY,  command VARCHAR(1000), baseimage VARCHAR(255), new_filename VARCHAR(255), album_id VARCHAR(255), album_deletehash VARCHAR(255));")
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmtPic.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}
	// **************************************************************************************************************************************************************************************************************************************

	// **************************************************************************************************************************************************************************************************************************************
	// Albums table
	stmtAlbum, err := Db.Prepare("CREATE TABLE Albums (ID INT AUTO_INCREMENT PRIMARY KEY, Album_Hash VARCHAR(255), Delete_Hash VARCHAR(255), Auth_Type VARCHAR(255), Token VARCHAR(255));")
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmtAlbum.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}
	// **************************************************************************************************************************************************************************************************************************************

	// **************************************************************************************************************************************************************************************************************************************
	// Tasking table
	stmtTask, err := Db.Prepare("CREATE TABLE Tasking (Tasking_Image VARCHAR(255), Tasking_Command VARCHAR(255), Response TEXT, Title VARCHAR(255), Tags VARCHAR(255), Agent VARCHAR(255), Image_Hash VARCHAR(255), Delete_Hash VARCHAR(255), Token VARCHAR(255));")
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmtTask.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}
	// **************************************************************************************************************************************************************************************************************************************

	// **************************************************************************************************************************************************************************************************************************************
	// Agent table
	stmtAgent, err := Db.Prepare("CREATE TABLE Agents (ID INT AUTO_INCREMENT PRIMARY KEY, Status VARCHAR(255), Title VARCHAR(255), Tags VARCHAR(255));")
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmtAgent.Exec()
	if err != nil {
		fmt.Println(err.Error())
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
	} else {
		fmt.Println("[+] Connected successfully")
	}

	_, err = Db.Exec("USE Anti")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("[+] DB selected Successfully")
	}

	insert, err := Db.Prepare("INSERT INTO Pictures (command, baseimage, new_filename) VALUES ( ?, ?, ? )")
	if err != nil {
		fmt.Println(err)
	}
	defer insert.Close()

	insert.Exec(command, baseimage, newfilename)
	fmt.Println("[+] Finally it worked...")
}
