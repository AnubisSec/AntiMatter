package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image/png"
	"os"

	stego "github.com/auyer/steganography"
)

func main() {
	// Open file
	inFile, _ := os.Open("dali.png")
	// Buffer reader
	reader := bufio.NewReader(inFile)
	// Decoding to Golang's image.Image()
	img, _ := png.Decode(reader)

	// Buffer that will recive the results
	w := new(bytes.Buffer)
	// Encode the message into the image
	err := stego.Encode(w, img, []byte("This is a Dali test"))
	if err != nil {
		fmt.Printf("Error encoding file %v", err)
		return
	}
	// Create new image
	outFile, _ := os.Create("out_file.png")
	// Write buffer to it
	w.WriteTo(outFile)
	outFile.Close()
}
