package ocr

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestOCR(t *testing.T) {
	file, err := os.Open("/Users/zhijundu/Library/Mobile Documents/com~apple~Preview/Documents/861732456487_.pic_副本.jpeg")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	var imageBytes bytes.Buffer
	io.Copy(&imageBytes, file)

	ocr := NewOCR()
	text, _ := ocr.RecongizeFromBytes(imageBytes.Bytes())

	fmt.Println(text)
}
