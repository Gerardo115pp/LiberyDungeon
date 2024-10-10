package helpers

import "io"

func ReadPreservingOffset(r io.ReadSeeker, byte_count int) ([]byte, error) {
	current_seek_offset, err := r.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	data := make([]byte, byte_count)
	_, err = r.Read(data)
	if err != nil {
		return data, err
	}

	_, err = r.Seek(current_seek_offset, io.SeekStart)
	if err != nil {
		return data, err
	}

	return data, nil
}
