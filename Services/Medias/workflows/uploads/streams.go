package uploads

import (
	"bytes"
	"context"
	"fmt"
	"io"
	dungeon_helpers "libery-dungeon-libs/helpers"
	dungeon_models "libery-dungeon-libs/models"
	app_config "libery_medias_service/Config"
	service_models "libery_medias_service/models"
	"libery_medias_service/repository"
	"libery_medias_service/workflows"
	workflow_errors "libery_medias_service/workflows/errors"
	"os"
	"path/filepath"
	"time"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

type completeChunkedMediaData struct {
	mediaIdentity *dungeon_models.MediaIdentity
	uploadTicket  *service_models.ChunkedUploadTicket
}

// Writes a chunk to a temporary file
func WriteChunkToFile(chunk_data []byte, upload_ticket *service_models.ChunkedUploadTicket, chunk_serial int) error {
	var chunk_filename string = upload_ticket.GetFileFragmentName(chunk_serial)
	chunk_filename = filepath.Join(app_config.UPLOAD_CHUNKS_PATH, chunk_filename)

	f, err := os.Create(chunk_filename)
	if err != nil {
		labeled_err := dungeon_models.NewLabeledError(err, "In WriteChunkToFile, while creating part file", dungeon_models.ErrIOError)
		return labeled_err
	}
	defer f.Close()

	if chunk_serial == 0 {
		upload_token, err := service_models.GenerateChunkedUploadTicket(upload_ticket, time.Now().Add(time.Hour*24), app_config.DOMAIN_SECRET)
		if err != nil {
			label_err := dungeon_models.NewLabeledError(err, "In WriteChunkToFile, while generating upload ticket for the file header", dungeon_models.ErrIOError)
			return label_err
		}

		token_data := service_models.SerializeChunkedUploadTicketToken(upload_token)

		chunk_data = append(token_data, chunk_data...)
	}

	chunk_buffer := bytes.NewBuffer(chunk_data)

	_, err = io.Copy(f, chunk_buffer)

	return err
}

// Creates a new media from a chunked upload
func CreateMediaFromChunkedUpload(ticket *service_models.ChunkedUploadTicket) error {
	if !CheckChunksPresent(ticket) {
		return fmt.Errorf("Not all chunks are present")
	}

	err := JoinUploadChunks(ticket)
	if err != nil {
		labeled_err := dungeon_models.NewLabeledError(err, "In CreateMediaFromChunkedUpload, while joining chunks", dungeon_models.ErrIOError)
		return labeled_err
	}

	new_media_identity, err := CreateMediaIdentityFromTicket(ticket)
	if err != nil {
		labeled_err := dungeon_models.NewLabeledError(err, "In CreateMediaFromChunkedUpload, while creating media identity", dungeon_models.ErrIOError)
		return labeled_err
	}

	err = workflows.InsertUniqueMediaIdentity(new_media_identity)
	if err != nil {
		labeled_err := dungeon_models.NewLabeledError(err, "In CreateMediaFromChunkedUpload, while inserting media identity", dungeon_models.ErrDB_CouldNotConnectToDB)
		return labeled_err
	}

	// TODO: This function takes too long to execute. we should probably either move it to a goroutine or change the entire process
	// chunked upload system to send the chunks directly to the final location.
	err = MoveChunkedMedia(&completeChunkedMediaData{new_media_identity, ticket})
	if err != nil {
		labeled_err := dungeon_models.NewLabeledError(err, "In CreateMediaFromChunkedUpload, while moving media file", dungeon_models.ErrIOError)
		return labeled_err
	}

	err = DeleteChunks(ticket)

	return err
}

// Create new media identity from a chunked upload
func CreateMediaIdentityFromTicket(ticket *service_models.ChunkedUploadTicket) (*dungeon_models.MediaIdentity, error) {
	var upload_category dungeon_models.Category

	echo.EchoDebug(fmt.Sprintf("Getting category with UUID: %s", ticket.CategoryUUID))
	upload_category, err := repository.CategoriesRepo.GetCategoryByID(context.Background(), ticket.CategoryUUID)
	if err != nil {
		labeled_err := dungeon_models.NewLabeledError(err, "In CreateMediaFromChunkedUpload, while getting category", dungeon_models.ErrDB_CouldNotConnectToDB)
		return nil, labeled_err
	}

	echo.EchoDebug(fmt.Sprintf("Getting cluster with ID: %d", upload_category.Cluster))
	upload_cluster, err := repository.CategoriesClustersRepo.GetClusterByID(context.Background(), upload_category.Cluster)
	if err != nil {
		labeled_err := dungeon_models.NewLabeledError(err, "In CreateMediaFromChunkedUpload, while getting cluster", dungeon_models.ErrDB_CouldNotConnectToDB)
		return nil, labeled_err
	}

	new_media := dungeon_models.CreateNewMedia(ticket.UploadFilename, upload_category.Uuid, dungeon_helpers.IsVideoFile(ticket.UploadFilename), 0)

	media_identity := dungeon_models.CreateNewMediaIdentity(new_media, &upload_category, &upload_cluster)

	return media_identity, nil
}

// Checks that all the chunks of a chunked upload are present
func CheckChunksPresent(ticket *service_models.ChunkedUploadTicket) bool {
	for h := 0; h < ticket.UploadChunks; h++ {
		chunk_filename := ticket.GetFileFragmentName(h)
		chunk_filename = filepath.Join(app_config.UPLOAD_CHUNKS_PATH, chunk_filename)

		if !dungeon_helpers.FileExists(chunk_filename) {
			return false
		}
	}

	return true
}

// Deletes all the chunks of a chunked upload
func DeleteChunks(ticket *service_models.ChunkedUploadTicket) error {
	for h := 0; h < ticket.UploadChunks; h++ {
		chunk_filename := ticket.GetFileFragmentName(h)
		chunk_filename = filepath.Join(app_config.UPLOAD_CHUNKS_PATH, chunk_filename)

		err := os.Remove(chunk_filename)
		if err != nil {
			labeled_err := dungeon_models.NewLabeledError(err, "In DeleteChunks, while deleting chunk file", dungeon_models.ErrFS_NoSuchFileOrDirectory)
			return labeled_err
		}
	}

	return nil
}

// Combines all the chunks of a chunked upload into a single file
func JoinUploadChunks(ticket *service_models.ChunkedUploadTicket) error {
	var file_upload_ticket *service_models.ChunkedUploadTicket
	upload_filename := filepath.Join(app_config.UPLOAD_CHUNKS_PATH, ticket.UploadFilename)

	new_media_file, err := os.Create(upload_filename)
	if err != nil {
		labeled_err := dungeon_models.NewLabeledError(err, "In JoinUploadChunks, while creating new media file", dungeon_models.ErrIOError)
		return labeled_err
	}
	defer new_media_file.Close()

	for h := 0; h < ticket.UploadChunks; h++ {
		chunk_filename := ticket.GetFileFragmentName(h)
		chunk_filename = filepath.Join(app_config.UPLOAD_CHUNKS_PATH, chunk_filename)

		chunk_file, err := os.Open(chunk_filename)
		if err != nil {
			chunk_file.Close()
			labeled_err := dungeon_models.NewLabeledError(err, "In JoinUploadChunks, while opening chunk file", dungeon_models.ErrIOError)
			return labeled_err
		}

		if h == 0 {
			file_upload_ticket, err = service_models.ReadChunkedUploadTicketToken(chunk_file, app_config.DOMAIN_SECRET)
			if err != nil {
				chunk_file.Close()
				labeled_err := dungeon_models.NewLabeledError(err, "In JoinUploadChunks, while reading upload ticket", dungeon_models.ErrIOError)
				return labeled_err
			}

			if !ticket.EqualUploadTicket(*file_upload_ticket) {
				chunk_file.Close()
				return dungeon_models.NewLabeledError(fmt.Errorf("Invalid upload ticket"), "Upload ticket does not match", workflow_errors.ErrUpload_TicketMismatch)
			}
		}

		_, err = io.Copy(new_media_file, chunk_file)
		if err != nil {
			chunk_file.Close()
			labeled_err := dungeon_models.NewLabeledError(err, "In JoinUploadChunks, while copying chunk file", dungeon_models.ErrIOError)
			return labeled_err
		}

		chunk_file.Close()
	}

	return nil
}

// Moves a an chunked upload to the final media location.
func MoveChunkedMedia(upload_data *completeChunkedMediaData) error {
	new_media_filename := filepath.Join(app_config.UPLOAD_CHUNKS_PATH, upload_data.uploadTicket.UploadFilename)
	category_path := filepath.Join(upload_data.mediaIdentity.ClusterPath, upload_data.mediaIdentity.CategoryPath)
	would_be_path := filepath.Join(category_path, upload_data.mediaIdentity.Media.Name)

	is_same_fs, err := dungeon_helpers.IsSameFilesystem(new_media_filename, category_path) // both paths must already exist so we can't use the would_be_path
	if err != nil {
		labeled_err := dungeon_models.NewLabeledError(err, "In CreateMediaFromChunkedUpload, while checking filesystem", dungeon_models.ErrIOError)
		return labeled_err
	}

	if is_same_fs {
		err = os.Rename(new_media_filename, would_be_path) // This is much faster but fails if the files are on different filesystems.
	} else {
		err = dungeon_helpers.MoveFile(new_media_filename, would_be_path) // This is slower but can handle different filesystems even on different network locations.
	}

	return err
}
