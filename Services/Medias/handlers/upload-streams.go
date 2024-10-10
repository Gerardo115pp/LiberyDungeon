package handlers

import (
	"fmt"
	"io"
	"libery-dungeon-libs/communication"
	"libery-dungeon-libs/dungeonsec/dungeon_middlewares"
	"libery-dungeon-libs/dungeonsec/dungeon_secrets"
	dungeon_helpers "libery-dungeon-libs/helpers"
	"libery-dungeon-libs/libs/libery_networking"
	dungeon_models "libery-dungeon-libs/models"
	app_config "libery_medias_service/Config"
	service_models "libery_medias_service/models"
	"libery_medias_service/repository"
	"libery_medias_service/workflows"
	upload_workflows "libery_medias_service/workflows/uploads"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

func UploadStreamsHandler(service_instance libery_networking.Server) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			getUploadStreamsHandler(response, request)
		case http.MethodPost:
			postUploadStreamsHandler(response, request)
		case http.MethodPatch:
			patchUploadStreamsHandler(response, request)
		case http.MethodDelete:
			deleteUploadStreamsHandler(response, request)
		case http.MethodPut:
			putUploadStreamsHandler(response, request)
		case http.MethodOptions:
			response.WriteHeader(http.StatusOK)
		default:
			response.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func getUploadStreamsHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path
	var handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

	switch resource {
	case "/upload-streams/chunked-ticket":
		handler_func = dungeon_middlewares.CheckUserCan_UploadFiles(getChunkedUploadStreamsTicketHandler)
	case "/upload-streams/stream-ticket":
		// getStreamUploadStreamsTicketHandler(response, request)
		handler_func = dungeon_middlewares.CheckUserCan_UploadFiles(getStreamUploadStreamsTicketHandler)
	}

	handler_func(response, request)
}

func getChunkedUploadStreamsTicketHandler(response http.ResponseWriter, request *http.Request) {
	var upload_stream_ticket *service_models.UploadStreamTicket

	upload_stream_ticket, err := service_models.ParseUploadStreamTicketFromRequest(request, dungeon_secrets.GetDungeonJwtSecret())
	if err != nil {
		echo.Echo(echo.RedBG, fmt.Sprintf("In getChunkedUploadStreamsTicketHandler: Error getting upload ticket because '%s'", err.Error()))
		response.WriteHeader(400)
		return
	}

	if upload_stream_ticket.MissingMedias() < 1 {
		echo.Echo(echo.RedBG, "In getChunkedUploadStreamsTicketHandler: No medias missing")
		response.WriteHeader(400)
		return
	}

	upload_ticket, err := service_models.NewChunkedUploadTicketFromRequest(request)

	upload_ticket_claims, err := service_models.GenerateChunkedUploadTicket(upload_ticket, time.Now().Add(time.Hour*24), app_config.JWT_SECRET)
	if err != nil {
		echo.Echo(echo.RedBG, fmt.Sprintf("In getChunkedUploadStreamsTicketHandler: Error generating upload ticket because '%s'", err.Error()))
		response.WriteHeader(500)
		return
	}

	response.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
	response.Header().Set("Pragma", "no-cache")

	dungeon_helpers.WriteSingleStringResponse(response, upload_ticket_claims)
}

func getStreamUploadStreamsTicketHandler(response http.ResponseWriter, request *http.Request) {
	echo.Echo(echo.YellowFG, "In getStreamUploadStreamsTicketHandler")
	upload_ticket, err := service_models.NewUploadStreamTicketFromRequest(request)
	if err != nil {
		echo.Echo(echo.RedBG, fmt.Sprintf("In getStreamUploadStreamsTicketHandler: Error getting upload ticket because '%s'", err.Error()))
		response.WriteHeader(400)
		return
	}

	category_uuid := request.URL.Query().Get("category_uuid")

	if category_uuid == "" {
		echo.Echo(echo.RedBG, "In getStreamUploadStreamsTicketHandler: Missing category_uuid query parameter")
		response.WriteHeader(400)
		return
	}

	existent_category, err := repository.CategoriesRepo.GetCategoryByID(request.Context(), category_uuid)
	if err != nil {
		echo.Echo(echo.RedBG, fmt.Sprintf("In getStreamUploadStreamsTicketHandler: Error getting category because '%s'", err.Error()))
		response.WriteHeader(404)
		return
	}

	category_cluster, err := repository.CategoriesClustersRepo.GetClusterByID(request.Context(), existent_category.Cluster)
	if err != nil {
		echo.Echo(echo.RedBG, fmt.Sprintf("In getStreamUploadStreamsTicketHandler: Error getting cluster because '%s'", err.Error()))
		response.WriteHeader(404)
		return
	}

	category_identity := dungeon_models.CreateNewCategoryIdentity(&existent_category, &category_cluster)
	if err != nil {
		echo.Echo(echo.RedBG, fmt.Sprintf("In getStreamUploadStreamsTicketHandler: Error creating category identity because '%s'", err.Error()))
		response.WriteHeader(500)
		return
	}

	upload_ticket.UploadCategoryIdentity = category_identity.ToWeakIdentity()

	err = upload_workflows.SetUploadStreamTicketCookie(response, upload_ticket)
	if err != nil {
		echo.Echo(echo.RedBG, fmt.Sprintf("In getStreamUploadStreamsTicketHandler: Error setting upload stream ticket cookie because '%s'", err.Error()))
		response.WriteHeader(500)
		return
	}

	response.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
	response.Header().Set("Pragma", "no-cache")

	dungeon_helpers.WriteBooleanResponse(response, true)
}

func postUploadStreamsHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path
	var handler_fund http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

	switch resource {
	case "/upload-streams/chunked-upload":
		handler_fund = postChunkedUploadStreamsHandler
	case "/upload-streams/stream-fragment":
		handler_fund = postStreamFragmentUploadStreamsHandler
	}

	handler_fund(response, request)
}

func postChunkedUploadStreamsHandler(response http.ResponseWriter, request *http.Request) {
	var upload_ticket_token string = request.Header.Get("Authorization")
	var upload_chunk_serial int
	if upload_ticket_token == "" {
		echo.Echo(echo.RedBG, "In postUploadStreamsHandler: Missing Authorization header")
		response.WriteHeader(401)
		return
	}

	// Check if the chunk serial was passed and is a valid number
	upload_chunk_serial, err := strconv.Atoi(request.URL.Query().Get("chunk_serial"))
	if err != nil {
		echo.Echo(echo.RedBG, "In postUploadStreamsHandler: Invalid chunk_serial query parameter")
		response.WriteHeader(400)
		return
	}

	// Check if the upload ticket is valid
	upload_ticket_claims, err := service_models.ParseChunkedUploadTicket(upload_ticket_token, app_config.JWT_SECRET)
	if err != nil {
		echo.Echo(echo.RedBG, fmt.Sprintf("In postUploadStreamsHandler: Token found<%s> but it was invalid because '%s'", upload_ticket_token, err.Error()))
		response.WriteHeader(403)
		return
	}

	// Check if the chunk serial is within the range of the accepted amount of chunks
	if upload_chunk_serial >= upload_ticket_claims.UploadChunks || upload_chunk_serial < 0 {
		echo.Echo(echo.RedBG, fmt.Sprintf("In postUploadStreamsHandler: Invalid chunk_serial query parameter<%d> for upload ticket<%s> with a maximum of %d chunks", upload_chunk_serial, upload_ticket_claims.UploadUUID, upload_ticket_claims.UploadChunks))
		response.WriteHeader(400)
		return
	}

	request.ParseMultipartForm(50 << 20) // 50 MB

	file, _, err := request.FormFile("chunk")
	if err != nil {
		echo.Echo(echo.RedBG, fmt.Sprintf("In postUploadStreamsHandler: Error getting file from form because '%s'", err.Error()))
		response.WriteHeader(400)
		return
	}
	defer file.Close()

	file_data, err := io.ReadAll(file)
	if err != nil {
		echo.Echo(echo.RedBG, fmt.Sprintf("In postUploadStreamsHandler: Error reading file data because '%s'", err.Error()))
		response.WriteHeader(500)
		return
	}

	err = upload_workflows.WriteChunkToFile(file_data, upload_ticket_claims, upload_chunk_serial)
	if err != nil {
		echo.Echo(echo.RedBG, fmt.Sprintf("In postUploadStreamsHandler: Error writing chunk to file because '%s'", err.Error()))
		response.WriteHeader(500)
		return
	}

	var is_last_chunk bool = upload_chunk_serial == (upload_ticket_claims.UploadChunks - 1)

	if is_last_chunk {
		err = upload_workflows.CreateMediaFromChunkedUpload(upload_ticket_claims)
		if err != nil {
			echo.Echo(echo.RedBG, fmt.Sprintf("In postUploadStreamsHandler: Error creating media from chunked upload because '%s'", err.Error()))
			response.WriteHeader(500)
			return
		}

		upload_stream_ticket, err := service_models.ParseUploadStreamTicketFromRequest(request, dungeon_secrets.GetDungeonJwtSecret())
		if err != nil {
			echo.Echo(echo.RedBG, fmt.Sprintf("In postUploadStreamsHandler: Error getting upload ticket because '%s'", err.Error()))
			response.WriteHeader(400)
			return
		}

		upload_stream_ticket.UploadedMedias++

		if upload_stream_ticket.UploadComplete() {
			fs_change_event := communication.NewClusterFSChangeEvent(dungeon_secrets.GetDungeonJwtSecret(), upload_stream_ticket.UploadCategoryIdentity.ClusterUUID, 0, upload_stream_ticket.TotalMedias, 0)

			err = fs_change_event.Emit()
			if err != nil {
				echo.Echo(echo.RedBG, fmt.Sprintf("In postUploadStreamsHandler: Error emitting fs_change_event because '%s'", err.Error()))
			}

			upload_workflows.DeleteUploadStreamTicketCookie(response)
		} else {
			err = upload_workflows.SetUploadStreamTicketCookie(response, upload_stream_ticket)
			if err != nil {
				echo.Echo(echo.RedBG, fmt.Sprintf("In postUploadStreamsHandler: Error setting upload stream ticket cookie because '%s'", err.Error()))
				response.WriteHeader(500)
				return
			}
		}

		response.WriteHeader(201)
	} else {
		response.WriteHeader(204)
	}
}

func postStreamFragmentUploadStreamsHandler(response http.ResponseWriter, request *http.Request) {
	var downloaded_from_uuid int64 = 0
	var err error

	upload_ticket, err := service_models.ParseUploadStreamTicketFromRequest(request, app_config.JWT_SECRET)
	if err != nil {
		echo.Echo(echo.RedBG, fmt.Sprintf("In postStreamFragmentUploadStreamsHandler: Error getting upload ticket because '%s'", err.Error()))
		response.WriteHeader(400)
		return
	}

	if upload_ticket.UploadComplete() {
		upload_workflows.DeleteUploadStreamTicketCookie(response)
		echo.Echo(echo.RedBG, "In postStreamFragmentUploadStreamsHandler: The used upload ticket is complete")
		response.WriteHeader(429)
		return
	}

	echo.Echo(echo.CyanFG, fmt.Sprintf("Processing upload %d of %d", upload_ticket.UploadedMedias+1, upload_ticket.TotalMedias))
	echo.EchoDebug(fmt.Sprintf("Uploading with claims: %+v", upload_ticket))

	download_from := request.URL.Query().Get("download_from") // ? is this being used for anything? i was under the impression that this was a junk property

	if download_from != "" {
		downloaded_from_uuid, err = strconv.ParseInt(download_from, 10, 64)
		if err != nil {
			http.Error(response, "Invalid download_from query parameter", 400)
			echo.Echo(echo.RedBG, fmt.Sprintf("Invalid download_from query parameter: %s", err.Error()))
			return
		}
	}

	err = request.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(response, "Error parsing multipart form", 400)
		echo.Echo(echo.RedBG, fmt.Sprintf("Error parsing multipart form: %s", err.Error()))
		return
	}

	var files_in_request int
	files_in_request, err = dungeon_helpers.CountFilesInMultipart(request)
	if err != nil {
		echo.Echo(echo.RedBG, fmt.Sprintf("Error counting files in multipart form: %s", err.Error()))
		response.WriteHeader(400)
		return
	}

	if files_in_request == 0 {
		echo.Echo(echo.RedBG, "In postStreamFragmentUploadStreamsHandler: No files in request")
		dungeon_helpers.WriteRejection(response, 400, "No files in request")
		return
	}

	if (upload_ticket.UploadedMedias + files_in_request) > upload_ticket.TotalMedias {
		echo.Echo(echo.RedBG, fmt.Sprintf("In postStreamFragmentUploadStreamsHandler: Too many files in request. Expected %d files, got %d", upload_ticket.TotalMedias, upload_ticket.UploadedMedias+files_in_request))
		dungeon_helpers.WriteRejection(response, 403, "Too many files in request")
		return
	}

	for _, file_headers := range request.MultipartForm.File {
		echo.Echo(echo.CyanFG, fmt.Sprintf("Processing %d files", len(file_headers)))
		for _, file_header := range file_headers {
			echo.Echo(echo.CyanFG, fmt.Sprintf("Processing file: %s", file_header.Filename))
			file, err := file_header.Open()
			if err != nil {
				http.Error(response, "Error opening file", 400)
				echo.Echo(echo.RedBG, fmt.Sprintf("Error opening file: %s", err.Error()))
				return
			}
			defer file.Close()

			// Get mime type
			mime_type := file_header.Header.Get("Content-Type")

			is_video := strings.HasPrefix(mime_type, "video")

			media := dungeon_models.CreateNewMedia(file_header.Filename, upload_ticket.UploadCategoryIdentity.CategoryUUID, is_video, downloaded_from_uuid)

			echo.Echo(echo.CyanFG, fmt.Sprintf("About to insert media: %s", media.Uuid))

			media_identity := upload_ticket.UploadCategoryIdentity.ToMediaIdentity(media)

			err = workflows.SaveMediaFile(media_identity, &file)
			if err != nil {
				http.Error(response, "Error saving media file", 500)
				echo.Echo(echo.RedBG, fmt.Sprintf("Error saving media file: %s", err.Error()))
				return
			}

			err = repository.MediasRepo.InsertMedia(request.Context(), media)
			if err != nil {
				http.Error(response, "Error inserting media", 500)
				echo.Echo(echo.RedBG, fmt.Sprintf("Error inserting media: %s", err.Error()))
				return
			}

			upload_ticket.UploadedMedias++
		}
	}

	if upload_ticket.UploadComplete() {
		fs_change_event := communication.NewClusterFSChangeEvent(app_config.JWT_SECRET, upload_ticket.UploadCategoryIdentity.ClusterUUID, 0, upload_ticket.TotalMedias, 0)

		err = fs_change_event.Emit()
		if err != nil {
			echo.Echo(echo.RedBG, fmt.Sprintf("In postStreamFragmentUploadStreamsHandler: Error emitting fs_change_event because '%s'", err.Error()))
		}

		echo.Echo(echo.GreenFG, "Finished uploading all medias")
		upload_workflows.DeleteUploadStreamTicketCookie(response)
		response.WriteHeader(201)
	} else {
		echo.Echo(echo.GreenFG, fmt.Sprintf("Uploaded %d of %d medias", upload_ticket.UploadedMedias, upload_ticket.TotalMedias))
		err = upload_workflows.SetUploadStreamTicketCookie(response, upload_ticket)
		if err != nil {
			echo.Echo(echo.RedBG, fmt.Sprintf("In postStreamFragmentUploadStreamsHandler: Error setting upload stream ticket cookie because '%s'", err.Error()))
			response.WriteHeader(500)
			return
		}

		response.WriteHeader(204)
	}
}

func patchUploadStreamsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
func deleteUploadStreamsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
func putUploadStreamsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
