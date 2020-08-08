package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"image/png"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	stego "github.com/auyer/steganography"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

// CreateImage is a function that uses the stego lib to encode an image you define and then write it to a new image
// if the original pic doesn't exist this will break...hard
func CreateImage(command string, origPic string, newPic string) {
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

func UploadImage(imageFile string, title string, album string, description string, clientID string) (imageID, deleteHash interface{}) {
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

func AddImage(albumDeleteHash string, clientID string, imgDeleteHash string) (success, status interface{}) {
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

func GetImage() bool {
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

// DecodeImage is a test function to decode stego
func DecodeImage() {
	// If this image doesn't exist it will break
	inFile, _ := os.Open("/tmp/asdf.jpg")
	defer inFile.Close()

	reader := bufio.NewReader(inFile)
	img, _ := png.Decode(reader)

	sizeOfMessage := stego.GetMessageSizeFromImage(img)

	msg := stego.Decode(sizeOfMessage, img)
	fmt.Println(string(msg))
}
