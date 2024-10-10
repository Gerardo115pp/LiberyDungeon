package models

import "bytes"

type MediaResponse struct {
	MediaStream *bytes.Buffer
	MimeType    string
	MediaLength int64
	Filename    string
}

type ThumbnailResponse struct {
	MediaResponse
	Resized      bool
	Size         *MediaSize
	OrignialSize *MediaSize
}
