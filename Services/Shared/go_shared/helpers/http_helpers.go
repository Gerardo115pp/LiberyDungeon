package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

type singleIntResponse struct {
	Response int `json:"response"`
}

type PaginatedResponseList[T any] struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalPages int `json:"total_pages"`
	TotalItems int `json:"total_items"`
	Content    []T `json:"content"`
}

type PaginatedResponse[T any] struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalPages int `json:"total_pages"`
	TotalItems int `json:"total_items"`
	Content    T   `json:"content"`
}

func createPaginatedResponseList[T any](content []T, page, total_pages, total_items int) *PaginatedResponseList[T] {
	var paginated *PaginatedResponseList[T] = new(PaginatedResponseList[T])

	paginated.Page = page
	paginated.PageSize = len(content)
	paginated.TotalPages = total_pages
	paginated.TotalItems = total_items
	paginated.Content = content

	return paginated
}

func createPaginatedResponse[T any](content T, page, total_pages, total_items int) *PaginatedResponse[T] {
	var paginated *PaginatedResponse[T] = new(PaginatedResponse[T])

	paginated.Page = page
	paginated.PageSize = 1
	paginated.TotalPages = total_pages
	paginated.TotalItems = total_items
	paginated.Content = content

	return paginated
}

func WritePaginatedResponseList[T any](response http.ResponseWriter, content []T, page, total_pages, total_items int) {
	paginated_response := createPaginatedResponseList(content, page, total_pages, total_items)

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(200)

	json.NewEncoder(response).Encode(paginated_response)
}

func WritePaginatedResponse[T any](response http.ResponseWriter, content T, page, total_pages, total_items int) {
	paginated_response := createPaginatedResponse(content, page, total_pages, total_items)

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(200)

	json.NewEncoder(response).Encode(paginated_response)
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

func WriteSingleIntResponse(response http.ResponseWriter, value int) {
	var single_int_response *singleIntResponse = new(singleIntResponse)
	single_int_response.Response = value

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(200)

	json.NewEncoder(response).Encode(single_int_response)
}

func WriteSingleIntResponseWithStatus(response http.ResponseWriter, value int, status int) {
	var single_int_response *singleIntResponse = new(singleIntResponse)
	single_int_response.Response = value

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)

	json.NewEncoder(response).Encode(single_int_response)
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

// =========================  Query Parameters Parsing =========================

// It parses a query parameter matched by the given key, is expected to be a string of the form "a1,a2,a3" and contain
// only numbers. It returns a slice of integers.
func ParseQueryParameterAsIntSlice(request *http.Request, key string) ([]int, error) {
	var query_param string = request.URL.Query().Get(key)

	if query_param == "" {
		return nil, fmt.Errorf("Query parameter %s is empty", key)
	}

	var values []int

	var numeric_members []string = strings.FieldsFunc(query_param, func(r rune) bool {
		return r == ','
	})

	for _, member := range numeric_members {
		var value int
		value, err := strconv.Atoi(member)
		if err != nil {
			return nil, fmt.Errorf("Error parsing query parameter %s: %s", key, err.Error())
		}

		values = append(values, value)
	}

	return values, nil
}
