package handlers

import (
	"encoding/json"
	"fmt"
	"libery-dungeon-libs/communication/service_requests/metadata_requests"
	"libery-dungeon-libs/dungeonsec/dungeon_middlewares"
	dungeon_helpers "libery-dungeon-libs/helpers"
	"libery-dungeon-libs/libs/libery_networking"
	video_moment_models "libery-metadata-service/models/video_moments"
	"libery-metadata-service/repository"
	"net/http"

	"github.com/Gerardo115pp/patriot_router"
	"github.com/Gerardo115pp/patriots_lib/echo"
)

var video_moments_path string = "/video-moments"

var VIDEO_MOMENTS_ROUTE *patriot_router.Route = patriot_router.NewRoute(fmt.Sprintf("%s(/.+)?", video_moments_path), false)

func VideoMomentsHandler(service_instance libery_networking.Server) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		var request_handler_func http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

		switch request.Method {
		case http.MethodGet:
			request_handler_func = getVideoMomentsHandler
		case http.MethodPost:
			request_handler_func = postVideoMomentsHandler
		case http.MethodPatch:
			request_handler_func = patchVideoMomentsHandler
		case http.MethodDelete:
			request_handler_func = deleteVideoMomentsHandler
		case http.MethodPut:
			request_handler_func = putVideoMomentsHandler
		case http.MethodOptions:
			request_handler_func = dungeon_helpers.AllowAllHandler
		default:
			request_handler_func = dungeon_helpers.MethodNotAllowedHandler
		}

		request_handler_func(response, request)
	}
}

func getVideoMomentsHandler(response http.ResponseWriter, request *http.Request) {
	var resource_path string = request.URL.Path
	var resource_handler http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

	switch resource_path {
	case fmt.Sprintf("%s/video-moments", video_moments_path):
		resource_handler = dungeon_middlewares.CheckUserCan_ViewContent(get__VideoMomentsHandler)
	}

	resource_handler(response, request)
}

func get__VideoMomentsHandler(response http.ResponseWriter, request *http.Request) {
	var request_params *metadata_requests.VideoMoments_VideoIdentifier = metadata_requests.ParseVideoIdentifierParams(request)

	if request_params == nil || request_params.VideoUUID == "" || request_params.VideoCluster == "" {
		echo.Echo(echo.RedFG, "In handlers/video_moments.get__VideoMomentsHandler: request was malformed, either video_uuid or video_cluster was empty")
		dungeon_helpers.WriteRejection(response, 400, "Missing parameters")
		return
	}

	video_instance, err := repository.VideoMomentsRepo.GetVideoCTX(request.Context(), request_params.VideoUUID, request_params.VideoCluster)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In handlers/video_moments.get__VideoMomentsHandler: While searching for the video. Video likely doesn't exists.\n\n%s", err))
		dungeon_helpers.WriteRejection(response, 404, "Video not found")
		return
	}

	moments, err := repository.VideoMomentsRepo.GetVideoMomentsCTX(request.Context(), video_instance)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In handlers/video_moments.get__VideoMomentsHandler: error getting video moments\n\n%s", err))
		dungeon_helpers.WriteRejection(response, 500, "Error getting video moments")
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(200)

	err = json.NewEncoder(response).Encode(moments)
}

func postVideoMomentsHandler(response http.ResponseWriter, request *http.Request) {
	var resource_path string = request.URL.Path
	var resource_handler http.HandlerFunc = dungeon_helpers.ResourceNotFoundHandler

	switch resource_path {
	case fmt.Sprintf("%s/moments", video_moments_path):
		resource_handler = dungeon_middlewares.CheckUserCan_ViewContent(post__NewVideoMomentHandler)
	}

	resource_handler(response, request)
}

func post__NewVideoMomentHandler(response http.ResponseWriter, request *http.Request) {
	var request_body *metadata_requests.VideoMoments_NewVideoMoment

	err := json.NewDecoder(request.Body).Decode(&request_body)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In handlers/video_moments.post__NewVideoMomentHandler: error decoding request body\n\n%s", err))
		dungeon_helpers.WriteRejection(response, 400, "Malformed request")
		return
	}

	if request_body.MomentTitle == "" || request_body.VideoCluster == "" || request_body.VideoUUID == "" {
		echo.Echo(echo.RedFG, "In handlers/video_moments.post__NewVideoMomentHandler: request was malformed, either moment_title, video_cluster or video_uuid was empty")
		dungeon_helpers.WriteRejection(response, 400, "Missing parameters")
		return
	}

	video_instance, err := repository.VideoMomentsRepo.GetVideoCTX(request.Context(), request_body.VideoUUID, request_body.VideoCluster)
	if err != nil {
		// If the video doesn't exist, create it.

		video_instance = &video_moment_models.Video{
			VideoUUID:    request_body.VideoUUID,
			VideoCluster: request_body.VideoCluster,
		}

		err = repository.VideoMomentsRepo.AddVideoCTX(request.Context(), *video_instance)
		if err != nil {
			echo.Echo(echo.RedFG, "In handlers/video_moments.post__NewVideoMomentHandler: error adding video\n\n%s", err)
			dungeon_helpers.WriteRejection(response, 500, "Error adding video")
			return
		}
	}

	video_moment_instance := video_moment_models.VideoMoment{
		VideoUUID:   video_instance.VideoUUID,
		MomentTime:  request_body.MomentTime,
		MomentTitle: request_body.MomentTitle,
	}

	moment_id, err := repository.VideoMomentsRepo.AddVideoMomentCTX(request.Context(), video_moment_instance)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In handlers/video_moments.post__NewVideoMomentHandler: error adding video moment\n\n%s", err))
		dungeon_helpers.WriteRejection(response, 500, "Error adding video moment")
		return
	}

	dungeon_helpers.WriteSingleIntResponseWithStatus(response, moment_id, 201)
}

func patchVideoMomentsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func deleteVideoMomentsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func putVideoMomentsHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
