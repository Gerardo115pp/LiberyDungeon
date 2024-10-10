package watch_points_database

import (
	"context"
	"fmt"
	"io"
	dungeon_models "libery-dungeon-libs/models"
	service_models "libery-metadata-service/models"
	"os"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

type WatchPointDatabase struct {
	mediaToOrd map[string]uint16
	stream     *service_models.WatchPointStream
}

func NewWatchPointDatabase() *WatchPointDatabase {
	var new_watch_point_database *WatchPointDatabase = new(WatchPointDatabase)
	new_watch_point_database.stream = service_models.NewWatchPointStream()

	var loaded_watch_points []byte

	loaded_watch_points, err := loadWatchPoints()
	if err != nil {
		echo.EchoWarn(fmt.Sprintf("Error loading watch points: %s", err))
	}

	new_watch_point_database.mediaToOrd = make(map[string]uint16)

	if loaded_watch_points != nil {
		err = new_watch_point_database.stream.LoadWatchPoints(loaded_watch_points)
		if err != nil {
			echo.EchoWarn(fmt.Sprintf("Error loading watch points: %s", err))
		}

		new_watch_point_database.populateMediaToOrd()
	}

	return new_watch_point_database
}

func (wdb *WatchPointDatabase) populateMediaToOrd() {
	wdb.stream.Seek(0, io.SeekStart)

	for watch_p, err := wdb.stream.Yield(); err == nil; watch_p, err = wdb.stream.Yield() {
		echo.EchoDebug(fmt.Sprintf("Populating mediaToOrd with media_uuid<%s> -> ord<%d>: time<%d>", watch_p.MediaUUID, watch_p.Ord, watch_p.StartTime))
		wdb.mediaToOrd[watch_p.MediaUUID] = watch_p.Ord
	}
}

func (wdb WatchPointDatabase) GetWatchPointByMediaID(ctx context.Context, media_uuid string) (*service_models.WatchPoint, *dungeon_models.LabeledError) {
	requested_ord, exists := wdb.mediaToOrd[media_uuid]

	if !exists {
		var target_watch_point *service_models.WatchPoint

		wdb.stream.Seek(0, io.SeekStart)

		for watch_p, err := wdb.stream.Yield(); err == nil; watch_p, err = wdb.stream.Yield() {
			echo.EchoDebug(fmt.Sprintf("Looking for media_uuid<%s> in watch point with media_uuid<%s>", media_uuid, watch_p))
			if watch_p.MediaUUID == media_uuid {
				target_watch_point = watch_p
				break
			}
		}

		if target_watch_point == nil {
			return nil, dungeon_models.NewLabeledError(fmt.Errorf("Watch point not found"), fmt.Sprintf("Watch point with media_uuid<%s> not found", media_uuid), service_models.ErrWatchPointNotFound)
		}

		return target_watch_point, nil
	}

	return wdb.stream.ReadWatchPoint(requested_ord)
}

func (wdb WatchPointDatabase) GetWatchPointByORD(ctx context.Context, ord uint16) (*service_models.WatchPoint, *dungeon_models.LabeledError) {
	return wdb.stream.ReadWatchPoint(ord)
}

func (wdb *WatchPointDatabase) InsertWatchPoint(ctx context.Context, media_uuid string, start_time uint32) *dungeon_models.LabeledError {
	requested_ord, exists := wdb.mediaToOrd[media_uuid]

	if exists {

		lerr := wdb.stream.UpdateWatchPointTime(requested_ord, start_time)
		if lerr != nil {
			return lerr
		}

	} else {

		new_watch_point := wdb.stream.AddWatchPoint(media_uuid, start_time)

		wdb.mediaToOrd[media_uuid] = new_watch_point.Ord
	}

	wdb.SaveToDisk()

	return nil
}

func (wdb WatchPointDatabase) SaveToDisk() error {
	var watch_points_filename string = getWatchPointsFilename()

	err := os.WriteFile(watch_points_filename, wdb.stream.WatchPoints, 0644)
	if err != nil {
		return dungeon_models.NewLabeledError(err, "Error writing watch points to file", dungeon_models.ErrProcessError)
	}

	return nil
}

func (wdb WatchPointDatabase) Close() error {
	echo.Echo(echo.BlueBG, "Closing watch point database")

	err := wdb.SaveToDisk()
	if err != nil {
		return err
	}

	return nil
}
