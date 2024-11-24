package ocr

import "github.com/otiai10/gosseract/v2"

type OCR struct {
	client *gosseract.Client
}

func NewOCR() *OCR {
	client := gosseract.NewClient()
	// 设置语言,chi_sim 表示简体中文
	client.SetLanguage("chi_sim")
	return &OCR{client: client}
}

// RecongizeFromFile 识别图片中的文字
func (o *OCR) RecongizeFromFile(imagePath string) (string, error) {
	defer o.client.Close()

	err := o.client.SetImage(imagePath)
	if err != nil {
		return "", err
	}

	text, err := o.client.Text()
	if err != nil {
		return "", err
	}
	return text, nil
}

// RecongizeFromBytes 识别图片中的文字
func (o *OCR) RecongizeFromBytes(imageBytes []byte) (string, error) {
	defer o.client.Close()

	err := o.client.SetImageFromBytes(imageBytes)
	if err != nil {
		return "", err
	}

	text, err := o.client.Text()
	if err != nil {
		return "", err
	}
	return text, nil
}
