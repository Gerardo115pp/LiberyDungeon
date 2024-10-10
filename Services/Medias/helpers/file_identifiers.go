package helpers

import (
	"io"
	"os"
	"strings"
)

// Checks if mime_type starts with "video"
func IsVideoMime(mime_type string) bool {
	return strings.HasPrefix(mime_type, "video")
}

// Checks if mime_type starts with "image"
func IsImageMime(mime_type string) bool {
	return strings.HasPrefix(mime_type, "image")
}

// ===== File content checks =====

func IsGIF(rs *os.File) (bool, error) {
	current_seek_position, err := rs.Seek(0, io.SeekCurrent)
	if err != nil {
		return false, err
	}

	file_header := make([]byte, 6)
	_, err = rs.Read(file_header)
	if err != nil {
		return false, err
	}

	_, err = rs.Seek(current_seek_position, io.SeekStart)
	if err != nil {
		return false, err
	}

	if file_header[0] != 0x47 || file_header[1] != 0x49 || file_header[2] != 0x46 { // Checks for G, I, F
		return false, nil
	}

	if file_header[3] != 0x38 || (file_header[4] != 0x37 && file_header[4] != 0x39) || file_header[5] != 0x61 { // Checks for 8, 7 or 9, a
		return false, nil
	}

	return true, nil
}
