package models

import (
	"fmt"
	"io"
	dungeon_models "libery-dungeon-libs/models"
)

const (
	WatchPointSize       uint8 = 46
	WatchPointIntroducer uint8 = 0x2C
)

type WatchPoint struct {
	Ord       uint16 `json:"ord"`
	MediaUUID string `json:"media_uuid"`
	StartTime uint32 `json:"start_time"`
}

func (w WatchPoint) Bytes() []byte {
	var ord_bytes = make([]byte, 2)
	ord_bytes[0] = byte(w.Ord)
	ord_bytes[1] = byte(w.Ord >> 8)

	media_uuid_bytes := []byte(w.MediaUUID)

	var start_time_bytes = make([]byte, 4)

	start_time_bytes[3] = byte(w.StartTime >> 24)
	start_time_bytes[2] = byte(w.StartTime >> 16)
	start_time_bytes[1] = byte(w.StartTime >> 8)
	start_time_bytes[0] = byte(w.StartTime)

	var watch_point_bytes = make([]byte, 46)

	copy(watch_point_bytes[0:2], ord_bytes)
	copy(watch_point_bytes[2:42], media_uuid_bytes)
	copy(watch_point_bytes[42:46], start_time_bytes)

	return watch_point_bytes
}

func (w WatchPoint) String() string {
	return fmt.Sprintf("WatchPoint{Ord: %d\n,\tMediaUUID: %s\n,\tStartTime: %d\b}", w.Ord, w.MediaUUID, w.StartTime)
}

func WatchPointFromBytes(watch_point_bytes []byte) (*WatchPoint, *dungeon_models.LabeledError) {
	if len(watch_point_bytes) != int(WatchPointSize) {
		return nil, dungeon_models.NewLabeledError(fmt.Errorf("Invalid watch point size"), fmt.Sprintf("Wanted 46 bytes, got %d", len(watch_point_bytes)), ErrInvalidWatchPointSize)
	}

	var new_watch_point *WatchPoint = new(WatchPoint)

	new_watch_point.Ord = uint16(watch_point_bytes[0]) | (uint16(watch_point_bytes[1]) << 8)
	new_watch_point.MediaUUID = string(watch_point_bytes[2:42])
	new_watch_point.StartTime = uint32(watch_point_bytes[42]) | (uint32(watch_point_bytes[43]) << 8) | (uint32(watch_point_bytes[44]) << 16) | (uint32(watch_point_bytes[45]) << 24)

	return new_watch_point, nil
}

type WatchPointStream struct {
	WatchPoints []byte
	SeekPoint   int
}

func NewWatchPointStream() *WatchPointStream {
	var new_watch_point_stream *WatchPointStream = new(WatchPointStream)

	new_watch_point_stream.WatchPoints = make([]byte, 0)

	return new_watch_point_stream
}

func (w *WatchPointStream) AddWatchPoint(media_uuid string, start_time uint32) *WatchPoint {
	var new_watch_point = WatchPoint{
		Ord:       uint16(w.WatchPointsCount()),
		MediaUUID: media_uuid,
		StartTime: start_time,
	}

	var new_watch_point_bytes []byte = make([]byte, WatchPointSize+1)
	new_watch_point_bytes[0] = WatchPointIntroducer
	copy(new_watch_point_bytes[1:], new_watch_point.Bytes())

	w.Write(new_watch_point_bytes)

	return &new_watch_point
}

func (w *WatchPointStream) LoadWatchPoints(watch_points []byte) error {
	var lerr *dungeon_models.LabeledError

	for h := 0; h < len(watch_points); h += (int(WatchPointSize) + 1) {
		if watch_points[h] != WatchPointIntroducer {
			lerr = dungeon_models.NewLabeledError(fmt.Errorf("Invalid bytes passed to LoadWatchPoints"), fmt.Sprintf("Invalid introducer at byte %#x", h), dungeon_models.ErrProcessError)
			return lerr
		}

		var watch_point_bytes = watch_points[h+1 : h+int(WatchPointSize)+1]

		_, lerr = WatchPointFromBytes(watch_point_bytes)
		if lerr != nil {
			return lerr
		}
	}

	w.WatchPoints = watch_points
	return nil
}

func (w *WatchPointStream) Read(p []byte) (n int, err error) {
	if w.SeekPoint >= len(w.WatchPoints) {
		return 0, nil
	}

	requested_chunk := len(p)

	if requested_chunk >= len(w.WatchPoints)-w.SeekPoint {
		err = io.EOF
	}

	n = copy(p, w.WatchPoints[w.SeekPoint:])

	w.SeekPoint += n

	return n, err
}

func (w WatchPointStream) ReadWatchPoint(ord uint16) (*WatchPoint, *dungeon_models.LabeledError) {
	var watch_point_size int = int(WatchPointSize)
	var watch_point_offset int = int(ord) * (watch_point_size + 1)

	if watch_point_offset >= len(w.WatchPoints) {
		return nil, dungeon_models.NewLabeledError(fmt.Errorf("Watch point not found"), fmt.Sprintf("Watch point with ord %d not found", ord), dungeon_models.ErrProcessError)
	}

	if w.WatchPoints[watch_point_offset] != WatchPointIntroducer {
		return nil, dungeon_models.NewLabeledError(fmt.Errorf("Invalid watch point introducer"), fmt.Sprintf("Watch point with ord %d has invalid introducer", ord), dungeon_models.ErrProcessError)
	}

	watch_point_offset++

	var watch_point_bytes = w.WatchPoints[watch_point_offset : watch_point_offset+watch_point_size]

	return WatchPointFromBytes(watch_point_bytes)
}

func (w *WatchPointStream) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		w.SeekPoint = int(offset)
	case io.SeekCurrent:
		w.SeekPoint += int(offset)
	case io.SeekEnd:
		w.SeekPoint = len(w.WatchPoints) + int(offset)
	}

	return int64(w.SeekPoint), nil
}

func (w *WatchPointStream) UpdateWatchPointTime(ord uint16, start_time uint32) *dungeon_models.LabeledError {
	var watch_point_size int = int(WatchPointSize)
	var watch_point_offset int = int(ord) * (watch_point_size + 1)

	if watch_point_offset >= len(w.WatchPoints) {
		return dungeon_models.NewLabeledError(fmt.Errorf("Watch point not found"), fmt.Sprintf("Watch point with ord %d not found", ord), dungeon_models.ErrProcessError)
	}

	if w.WatchPoints[watch_point_offset] != WatchPointIntroducer {
		return dungeon_models.NewLabeledError(fmt.Errorf("Invalid watch point introducer"), fmt.Sprintf("Watch point with ord %d has invalid introducer", ord), dungeon_models.ErrProcessError)
	}

	watch_point_offset++

	var new_start_time_bytes = make([]byte, 4)

	new_start_time_bytes[3] = byte(start_time >> 24)
	new_start_time_bytes[2] = byte(start_time >> 16)
	new_start_time_bytes[1] = byte(start_time >> 8)
	new_start_time_bytes[0] = byte(start_time)

	copy(w.WatchPoints[watch_point_offset+42:watch_point_offset+46], new_start_time_bytes)

	return nil
}

func (w *WatchPointStream) Write(p []byte) (n int, err error) {
	w.WatchPoints = append(w.WatchPoints, p...)

	return len(p), nil
}

func (w WatchPointStream) WatchPointsCount() int {
	var watch_point_size = int(WatchPointSize) + 1 // Data size + 1 byte for the introducer

	return len(w.WatchPoints) / watch_point_size
}

func (w *WatchPointStream) Yield() (*WatchPoint, *dungeon_models.LabeledError) {
	if w.SeekPoint >= len(w.WatchPoints) {
		return nil, dungeon_models.NewLabeledError(io.EOF, "End of stream", ErrEndOfStream)
	}

	var watch_point_size int = int(WatchPointSize)

	if w.WatchPoints[w.SeekPoint] != WatchPointIntroducer {
		return nil, dungeon_models.NewLabeledError(fmt.Errorf("Invalid watch point introducer"), fmt.Sprintf("Watch point with ord %d has invalid introducer", w.SeekPoint), dungeon_models.ErrProcessError)
	}

	w.SeekPoint++

	var watch_point_bytes = w.WatchPoints[w.SeekPoint : w.SeekPoint+watch_point_size]

	w.SeekPoint += watch_point_size

	return WatchPointFromBytes(watch_point_bytes)
}
