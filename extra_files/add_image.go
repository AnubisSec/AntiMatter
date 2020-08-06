package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
)

func main() {

	// THis hash is the album deleteHash
	url := "https://api.imgur.com/3/album/RXtLCumUBxGxmhl"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	// This hash is the image deleteHash
	_ = writer.WriteField("deletehashes[]", "u0lhTkgvLOx5k41")
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
	var result map[string]interface{}

	json.NewDecoder(res.Body).Decode(&result)

	success := result["success"]
	status := result["status"]

	fmt.Println("Was this a success:", success, "Status:", status)
}
