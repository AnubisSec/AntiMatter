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
	img, _ := png.Decode(reader)

	sizeOfMessage := stego.GetMessageSizeFromImage(img)

	msg := stego.Decode(sizeOfMessage, img)
	fmt.Println(string(msg))
}
