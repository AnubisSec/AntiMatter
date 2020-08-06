package main

import (
	"bufio"
	"fmt"
	"image/png"
	"os"

	stego "github.com/auyer/steganography"
)

func main() {

	inFile, _ := os.Open("/tmp/asdf.jpg")
	defer inFile.Close()

	reader := bufio.NewReader(inFile)
	img, err := png.Decode(reader)
	if err != nil {
		fmt.Println(err)
	}

	sizeOfMessage := stego.GetMessageSizeFromImage(img)
	fmt.Println(sizeOfMessage)

	msg := stego.Decode(sizeOfMessage, img)
	fmt.Println(string(msg))
}
