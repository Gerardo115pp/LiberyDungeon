package helpers

import (
	"errors"
	_ "image/png" // import PNG format for image.Decode
	"io"
	"os"
	"strings"
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

func FileExists(file_name string) bool {
	if _, err := os.Stat(file_name); os.IsNotExist(err) {
		return false
	}
	return true
}

func RemoveSpecialChars(str string) string {
	var special_chars []string = []string{" ", "!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "+", "=", "[", "]", "{", "}", "|", "\\", "/", "?", ",", "<", ">", "'", "\"", "`", "~", ":", ";"}
	for _, special_char := range special_chars {
		str = strings.Replace(str, special_char, "_", -1)
	}
	return str
}

func GetFileDescriptor(media_path string) (*os.File, error) {
	if !FileExists(media_path) {
		return nil, errors.New("Media file not found")
	}

	return os.Open(media_path)
}
