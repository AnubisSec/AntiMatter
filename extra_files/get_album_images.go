package main

import (
	"bytes"
	//"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"reflect"
	"strings"
	//"strconv"
)

// AlbumImages is a test struct
type AlbumImages struct {
	ImageID     string `json:"id"`
	ImageTitle  string `json:"title"`
	Description string `json:"description"`
	ImageLink   string `json:"link"`
}

func main() {

	// This hash is the albumID hash
	url := "https://api.imgur.com/3/album/3OR2oAy/images"
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
	req.Header.Add("Authorization", "Client-ID bc8c890066d6157")

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	//var result map[string]interface{}

	blah := strings.NewReplacer(`{"data":[`, "", "]", "", ":[", ":0", `"}`, `"`, `\`, "")
	content := AlbumImages{}

	test := blah.Replace(string(body))
	//test := `{"id":"K1s5Q3k","title":"blahh","description":"within this response there is a key","datetime":1594869138,"type":"image/png","animated":false,"width":200,"height":200,"size":7348,"views":0,"bandwidth":0,"vote":null,"favorite":false,"nsfw":null,"section":null,"account_url":null,"account_id":null,"is_ad":false,"in_most_viral":false,"has_sound":false,"tags":0,"ad_type":0,"ad_url":"","edited":"0","in_gallery":false,"success":true,"status":200}`
	//fmt.Println(test)

	json.Unmarshal([]byte(test), &content)

	v := reflect.ValueOf(content)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fmt.Printf("[+] %s: %v\n", typeOfS.Field(i).Name, v.Field(i).Interface())
	}

	/*
		fmt.Printf("%+v\n", content)
		fmt.Printf("%T\n", content)
	*/
	//fmt.Println(result)

	//fmt.Printf("%T\n", test)
	//fmt.Printf("%+v\n", string(body))
	//fmt.Println(test)
	//fmt.Println(nestedMap)

}
