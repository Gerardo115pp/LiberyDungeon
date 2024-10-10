package service_clients

import (
	"bytes"
	"fmt"
	"io"
	"libery-dungeon-libs/dungeonsec"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/textproto"
)

type MediaServiceClient struct {
	BaseServiceClient
	tickets_cookie_jar *cookiejar.Jar
}

func (medias_client MediaServiceClient) Alive() (bool, error) {
	var metadata_endpoint string

	metadata_endpoint = medias_client.getHttpsEndpoint()

	metadata_endpoint += "/alive"

	request, err := http.NewRequest("GET", metadata_endpoint, nil)
	if err != nil {
		return false, err
	}

	client := &http.Client{
		Transport: medias_client.HttpTransport,
	}

	response, err := client.Do(request)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	return response.StatusCode >= 200 && response.StatusCode < 300, nil
}

func (medias_client MediaServiceClient) getHttpsEndpoint() string {
	return fmt.Sprintf("https://%s%s", medias_client.BaseDomain, medias_client.HttpAddress)
}

// Requests the media service to send an upload ticket for a new stream of medias. This ticket is sent as a cookie
// which is stored in a cookie jar and that cookie jar in the MediaServiceClient.tickets_cookie_jar field.
func (medias_client *MediaServiceClient) GetUploadStreamTicket(upload_uuid, category_uuid string, total_medias int) error {
	var endpoint string = medias_client.getHttpsEndpoint()
	request_url := fmt.Sprintf("%s/upload-streams/stream-ticket?upload_uuid=%s&total_medias=%d&category_uuid=%s", endpoint, upload_uuid, total_medias, category_uuid)

	request, err := http.NewRequest("GET", request_url, nil)
	if err != nil {
		return fmt.Errorf("Error creating request: %s", err.Error())
	}

	err = dungeonsec.SignInternalHTTPRequest(request)
	if err != nil {
		return fmt.Errorf("Error signing internal request: %s\nThis is likely because someone skipped a configuration step", err.Error())
	}

	var clear_cookiejar *cookiejar.Jar

	clear_cookiejar, err = cookiejar.New(nil)
	if err != nil {
		return fmt.Errorf("Error creating cookie jar: %s", err.Error())
	}

	client := &http.Client{
		Transport: medias_client.HttpTransport,
		Jar:       clear_cookiejar,
	}

	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("Error sending request: %s", err.Error())
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return fmt.Errorf("Error sending request: %s", response.Status)
	}

	medias_client.tickets_cookie_jar = clear_cookiejar

	return nil
}

func (medias_client *MediaServiceClient) UploadMediaFile(file io.ReadCloser, filename string) error {
	message_body := &bytes.Buffer{}

	file_type_buffer := make([]byte, 512)
	n, err := file.Read(file_type_buffer)
	if err != nil {
		return fmt.Errorf("Error reading file type: %s", err.Error())
	}

	content_type := http.DetectContentType(file_type_buffer[:n])

	multi_reader := io.MultiReader(bytes.NewReader(file_type_buffer[:n]), file)

	writer := multipart.NewWriter(message_body)

	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filename))
	header.Set("Content-Type", content_type)

	part, err := writer.CreatePart(header)
	if err != nil {
		return fmt.Errorf("Error adding file to multipart message: %s", err.Error())
	}

	_, err = io.Copy(part, multi_reader)
	if err != nil {
		return fmt.Errorf("Error copying file content to multipart message: %s", err.Error())
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("Error closing multipart message: %s", err.Error())
	}

	return medias_client.sendNewMediaFile(message_body, writer.FormDataContentType())
}

func (medias_client *MediaServiceClient) sendNewMediaFile(message_body *bytes.Buffer, content_type string) error {
	var endpoint string = medias_client.getHttpsEndpoint()
	request_url := fmt.Sprintf("%s/upload-streams/stream-fragment", endpoint)

	request, err := http.NewRequest("POST", request_url, message_body)
	if err != nil {
		return fmt.Errorf("Error creating request: %s", err.Error())
	}

	request.Header.Set("Content-Type", content_type)

	client := &http.Client{
		Transport: medias_client.HttpTransport,
		Jar:       medias_client.tickets_cookie_jar,
	}

	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("Error sending request: %s", err.Error())
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return fmt.Errorf("Error sending request: %s", response.Status)
	}

	return nil
}
