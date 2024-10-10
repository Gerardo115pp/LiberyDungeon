package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type httpRejection struct {
	Code  int    `json:"code"`
	Cause string `json:"cause"`
}

type booleanResponse struct {
	Response bool `json:"response"`
}

type ReasonedBooleanResponse struct {
	booleanResponse
	Reason string `json:"reason"`
}

type singleStringResponse struct {
	Response string `json:"response"`
}

func createRejection(code int, message string) *httpRejection {
	var rejection *httpRejection = new(httpRejection)
	rejection.Code = code
	rejection.Cause = message

	return rejection
}

func WriteRejection(response http.ResponseWriter, code int, message string) {
	var rejection *httpRejection = createRejection(code, message)

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(code)

	json.NewEncoder(response).Encode(rejection)
}

func WriteBooleanResponse(response http.ResponseWriter, value bool) {
	var boolean_response *booleanResponse = new(booleanResponse)
	boolean_response.Response = value

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(200)

	json.NewEncoder(response).Encode(boolean_response)
}

func WriteReasonedBooleanResponse(response http.ResponseWriter, value bool, reason string) {
	var reasoned_boolean_response *ReasonedBooleanResponse = new(ReasonedBooleanResponse)
	reasoned_boolean_response.Response = value
	reasoned_boolean_response.Reason = reason

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(200)

	json.NewEncoder(response).Encode(reasoned_boolean_response)
}

func WriteSingleStringResponse(response http.ResponseWriter, value string) {
	var single_string_response *singleStringResponse = new(singleStringResponse)
	single_string_response.Response = value

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(200)

	json.NewEncoder(response).Encode(single_string_response)
}

func WriteSingleStringResponseWithStatus(response http.ResponseWriter, value string, status int) {
	var single_string_response *singleStringResponse = new(singleStringResponse)
	single_string_response.Response = value

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)

	json.NewEncoder(response).Encode(single_string_response)
}

// ========================= Generic Handlers =========================
// These are hanblers for very generic cases like method not allowed, route not found, etc. Avoid using for something that should have a specific context like a cluster not been found
// a user not been allowed to do something, etc. Use your good judgement.

func MethodNotAllowedHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(405)
}

func ResourceNotFoundHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(404)
}

func NotAuthorizedHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(401)
}

func ForbiddenHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(403)
}

func AllowAllHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(200)
}

// =========================  Cookies =========================

func DeleteCookie(response http.ResponseWriter, cookie_name string) {
	http.SetCookie(response, &http.Cookie{
		Name:   cookie_name,
		Value:  "",
		MaxAge: -1,
	})
}

// =========================  MultiPart Form =========================

func CountFilesInMultipart(request *http.Request) (int, error) {
	if request.MultipartForm == nil {
		return 0, fmt.Errorf("In CountFilesInMultipart: MultipartForm is nil. Likely because ParseMultipartForm was not called")
	}

	var count int = 0

	for _, file_headers := range request.MultipartForm.File {
		count += len(file_headers)
	}

	return count, nil
}
