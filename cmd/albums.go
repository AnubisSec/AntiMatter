package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// AlbumImages is struct that holds data structs for GetAlbumImages()
type AlbumImages struct {
	Data []struct {
		ID          string   `json:"id"`
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Datetime    int64    `json:"datetime"`
		Imagetype   string   `json:"type"`
		Nsfw        bool     `json:"nsfw"`
		Tags        []string `json:"tags"`
		InGallery   bool     `json:"in_gallery"`
		Link        string   `json:"link"`
	} `json:"data"`
	Success bool `json:"success"`
	Status  int  `json:"status"`
}

// CreateAlbum queries the Imgur API in order to create an anonymous album, using a client ID
func CreateAlbum(title string, clientID string) (albumID, deleteHash interface{}) {

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
	fmt.Println(color.GreenString("albumID:"), albumID, color.GreenString("Album DeleteHash:"), deleteHash)
	fmt.Println(" ")

	return albumID, deleteHash

}

// GetAlbumImages is a function that supposed to retrieve response images
func GetAlbumImages(albumID string, clientID string) { // removed: (imageLink interface{})

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
	if err != nil {
		fmt.Println("[-] Error connecting:", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	var results AlbumImages
	errr := json.Unmarshal([]byte(body), &results)
	if errr != nil {
		fmt.Println("[!] Error unmarshalling::", errr)
	}

	datavalues := results.Data
	if results.Success == true {
		for field := range datavalues {
			fmt.Println("[+] ImageID:", datavalues[field].ID)
			fmt.Println("[+] ImageTitle:", datavalues[field].Title)
			fmt.Println("[+] Description:", datavalues[field].Description)
			fmt.Println("[+] ImageLink:", datavalues[field].Link)
			fmt.Println(" ")
		}

	}

}

// GetLinkClient is a function is to grab the tasking link for the client
func GetLinkClient(albumID string, clientID string) (imageLink string) {

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
	if err != nil {
		fmt.Println("[-] Error connecting:", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	var results AlbumImages
	errr := json.Unmarshal([]byte(body), &results)
	if errr != nil {
		fmt.Println("[!] Error unmarshalling::", errr)
	}

	datavalues := results.Data
	if results.Success == true {
		imageLink = datavalues[0].Link
	}
	return imageLink

}
