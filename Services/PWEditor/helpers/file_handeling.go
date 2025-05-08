package helpers

import (
	"bytes"
	"image"
	"image/jpeg"
	_ "image/png" // import PNG format for image.Decode
	"io"
	"mime/multipart"
	"os"
)

func CopyFile(from_name string, to_name string) (err error) {
	from, err := os.Open(from_name)
	if err != nil {
		return
	}
	defer from.Close()

	to, err := os.OpenFile(to_name, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	return
}

func IsJPEG(file multipart.File) (bool, error) {
	_, image_format, err := image.Decode(file)
	if err != nil {
		return false, err
	}
	file.Seek(0, 0)
	return image_format == "jpeg", err
}

func ConvertToJPEG(file multipart.File) ([]byte, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	// Create a new buffer
	var buf bytes.Buffer

	// Encode the image as a JPEG image
	err = jpeg.Encode(&buf, img, nil)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
