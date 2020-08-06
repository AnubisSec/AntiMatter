package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

func main() {

	url := "https://api.imgur.com/3/album/{{albumDeleteHash}}"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("deletehashes[]", "{{imageDeleteHash}}")
	_ = writer.WriteField("deletehashes[]", "{{imageDeleteHash2}}")
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Authorization", "Client-ID {{clientId}}")

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
}
