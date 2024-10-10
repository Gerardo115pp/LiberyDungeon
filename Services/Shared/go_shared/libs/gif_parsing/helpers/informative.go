package helpers

import (
	"fmt"
	"io"
)

func PrintFileOffset(opened_file io.ReadSeeker) {
	file_offset, err := opened_file.Seek(0, io.SeekCurrent)
	if err != nil {
		fmt.Printf("Error while getting file offset: %s\n", err.Error())
		return
	}

	fmt.Printf("Current file offset: %#x\n", file_offset)
}
