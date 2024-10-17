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
	"path/filepath"
	"strings"
	"syscall"
)

var supported_mime_types = [...]string{
	"image/jpeg",
	"image/png",
	"image/gif",
	"image/bmp",
	"image/webp",
	"video/mp4",
	"video/webm",
	"video/matroska",
}

var supported_file_extensions = [...]string{
	".jpg",
	".jpeg",
	".png",
	".gif",
	".bmp",
	".webp",
	".mp4",
	".m4v",
	".webm",
	".mkv",
}

func IsSupportedMimeType(mime_type string) bool {
	mime_type = strings.ToLower(mime_type)
	var is_supported bool = false
	for _, supported_mime := range supported_mime_types {
		if supported_mime == mime_type {
			is_supported = true
		}
	}

	return is_supported
}

func IsSupportedFileExtension(filename string) bool {
	filename = strings.ToLower(filename)
	var extension = filepath.Ext(filename)

	var is_supported bool = false
	for _, supported_extension := range supported_file_extensions {
		if supported_extension == extension {
			is_supported = true
		}
	}

	return is_supported
}

func IsVideoFile(filename string) bool {
	var extension = filepath.Ext(filename)

	return extension == ".mp4" || extension == ".m4v" || extension == ".webm" || extension == ".mkv"
}

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

// Checks if two paths are on the same filesystem. expects both paths to exist.
func IsSameFilesystem(path_1, path_2 string) (bool, error) {
	file_stat_1, err := os.Stat(path_1)
	if err != nil {
		return false, err
	}

	file_stat_2, err := os.Stat(path_2)
	if err != nil {
		return false, err
	}

	stat1 := file_stat_1.Sys().(*syscall.Stat_t)
	stat2 := file_stat_2.Sys().(*syscall.Stat_t)

	return stat1.Dev == stat2.Dev && stat1.Ino == stat2.Ino, nil
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

func GetParentDirectory(path_str string) string {
	if path_str[len(path_str)-1] == '/' {
		path_str = path_str[:len(path_str)-1]
	}

	parent_dir, _ := filepath.Split(path_str)

	return parent_dir
}

// Returns a new path with the path_str base element renamed to new_name and always ending with a slash
// E.g. RenameFsPath("/path/to/directory/", "new_directory") -> "/path/to/new_directory/"
// And also RenameFsPath("/path/to/directory", "new_directory") -> "/path/to/new_directory/" (notice the added slash)
func RenameFsPath(path_str, new_name string) string {
	parent_dir := GetParentDirectory(path_str)

	return fmt.Sprintf("%s/", filepath.Join(parent_dir, new_name))
}

// Renames a filename to a new given name keeping the same extension. Does not perform any filesystem operations.
// It expect a extensionless new_name. and a filename with extension.
func RenameFilename(filename, new_name string) string {
	filename_ext := filepath.Ext(filename)

	return fmt.Sprintf("%s%s", new_name, filename_ext)
}

// Returns true if child_path is a child of parent_path
func IsChildPath(parent_path, child_path string) bool {
	return strings.HasPrefix(child_path, parent_path)
}

// Returns a platform normilized path e.g would turn '/home/el_maligno/SoftwareProjects/LiberyDungeon' into '/home/el_maligno/SoftwareProjects/LiberyDungeon/'
func NormalizePath(path string) string {
	normalized_path := filepath.Clean(path)

	if normalized_path != "/" {
		normalized_path = fmt.Sprintf("%s/", normalized_path)
	}

	return normalized_path
}
