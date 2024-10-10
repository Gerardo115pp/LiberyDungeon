package helpers

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http"
)

func IsJPEG(file multipart.File) (bool, error) {
	_, image_format, err := image.Decode(file)
	if err != nil {
		return false, err
	}
	file.Seek(0, 0)
	return image_format == "jpeg", err
}

func ConvertMultipartToJPEG(file multipart.File) ([]byte, error) {
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

func GetMimeType(f io.ReadSeeker) (string, error) {
	current_seek_position, err := f.Seek(0, io.SeekCurrent)
	if err != nil {
		return "", err
	}

	file_header := make([]byte, 512)
	f.Read(file_header)

	f.Seek(current_seek_position, io.SeekStart)

	content_type := http.DetectContentType(file_header)

	return content_type, nil
}
