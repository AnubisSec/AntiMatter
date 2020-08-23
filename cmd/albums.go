package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// AlbumImages is a struct to hold the relevant image info after uploading to an album
type AlbumImages struct {
	ImageID     string `json:"id"`
	ImageTitle  string `json:"title"`
	Description string `json:"description"`
	ImageLink   string `json:"link"`
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
func GetAlbumImages(albumID string, clientID string) (imageLink interface{}) {

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
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	// A dirty way of stripping uneeded json garbage since I can't figure out how to do it another way
	stripResponse := strings.NewReplacer(`{"data":[`, "", "]", "", ":[", ":0", `"}`, `"`, `\`, "")
	//Init the AlbumImages struct
	content := AlbumImages{}

	newResponse := stripResponse.Replace(string(body))

	json.Unmarshal([]byte(newResponse), &content)

	v := reflect.ValueOf(content)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fmt.Printf(color.GreenString("[+]")+" %s: %v\n", typeOfS.Field(i).Name, v.Field(i).Interface())
	}

	link := v.Field(3).Interface()

	return link

}
