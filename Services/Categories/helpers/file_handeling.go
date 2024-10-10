package helpers

import (
	"bytes"
	"fmt"
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

// MoveFile moves a file from source to dest even if different file systems are involved.
func MoveFile(source_file, dest_file string) error {
	inputFile, err := os.Open(source_file)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(dest_file)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(source_file)
	if err != nil {
		return fmt.Errorf("Failed removing original file: %s", err)
	}
	return nil
}

func IsDirectoryEmpty(path string) (bool, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}

	return len(files) == 0, nil
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

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
