package uploads

import (
	"fmt"
	app_config "libery_medias_service/Config"
	service_models "libery_medias_service/models"
	"net/http"
	"time"
)

func SetUploadStreamTicketCookie(response http.ResponseWriter, upload_ticket *service_models.UploadStreamTicket) error {
	if upload_ticket.UploadComplete() {
		return fmt.Errorf("In SetUploadStreamTicket: Upload ticket has already uploaded all medias")
	}

	var expire_at time.Time = upload_ticket.GenerousTimeToUpload()

	progressed_token, err := upload_ticket.Sign(app_config.JWT_SECRET)
	if err != nil {
		return fmt.Errorf("In SetUploadStreamTicket: Error signing upload ticket because '%s'", err.Error())
	}

	http.SetCookie(response, &http.Cookie{
		Name:     app_config.UPLOAD_STREAM_TICKET_COOKIE_NAME,
		Value:    progressed_token,
		Expires:  expire_at,
		HttpOnly: true,
		Path:     "/",
	})

	return nil
}

func DeleteUploadStreamTicketCookie(response http.ResponseWriter) {
	var expired_cookie *http.Cookie = &http.Cookie{
		Name:     app_config.UPLOAD_STREAM_TICKET_COOKIE_NAME,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Path:     "/",
	}

	http.SetCookie(response, expired_cookie)
}
