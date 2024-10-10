package workflows

import (
	"fmt"
	"io"
	gif_models "libery-dungeon-libs/libs/gif_parsing/models"
	"os"
)

func OpenReadableFile(filename string) (*os.File, error) {

	gif_stat, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("Gif file does not exist")
	} else if err != nil {
		return nil, fmt.Errorf("Error while checking gif file: %s", err.Error())
	}

	if gif_stat.Mode().Perm()&0400 == 0 {
		return nil, fmt.Errorf("Gif file is not readable")
	}

	gif_file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Error while opening gif file: %s", err.Error())
	}

	return gif_file, nil
}

func ReadGifFile(gif_file *os.File) (*gif_models.ParsedGif, error) {
	var gif_data *gif_models.ParsedGif

	gif_data, err := gif_models.ParseGifGlobalData(gif_file)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing gif global data: %s", err.Error())
	}

	gif_data.Filename = gif_file.Name()

	err = gif_models.ParseGifBlocks(gif_file, gif_data)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing gif blocks: %s", err.Error())
	}

	current_seek_position, err := gif_file.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, fmt.Errorf("Error while getting current seek position: %s", err.Error())
	}

	gif_data.TrailerFilePosition = current_seek_position

	return gif_data, nil
}
