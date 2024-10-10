package models

import (
	"fmt"
	"io"
	dungeon_models "libery-dungeon-libs/models"
	app_config "libery_medias_service/Config"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

// ========================= Upload Stream Ticket =========================

type UploadStreamTicket struct {
	UploadUUID             string                               `json:"upload_uuid"`
	TotalMedias            int                                  `json:"total_medias"`
	UploadedMedias         int                                  `json:"uploaded_medias"`
	UploadCategoryIdentity *dungeon_models.CategoryWeakIdentity `json:"upload_category_identity"`
}

func (ut UploadStreamTicket) UploadComplete() bool {
	return ut.UploadedMedias == ut.TotalMedias
}

func (ut UploadStreamTicket) MissingMedias() int {
	return ut.TotalMedias - ut.UploadedMedias
}

func (ut UploadStreamTicket) GenerousTimeToUpload() time.Time {
	var missing_medias int = ut.MissingMedias()
	var time_to_upload_per_media time.Duration = 10 * time.Minute

	var expire_at time.Time = time.Now().Add(time.Duration(missing_medias) * time_to_upload_per_media)

	return expire_at
}

func (ut *UploadStreamTicket) ProgressUploadedMedias() error {
	if (ut.UploadedMedias + 1) > ut.TotalMedias {
		return fmt.Errorf("In ProgressUploadedMedias: Upload ticket has already uploaded all medias")
	}

	ut.UploadedMedias += 1

	return nil
}

// Signs the upload stream ticket with the given secret key and returns the signed token. The expiration time for the token is set to based on
// the amount of medias left to upload.
func (ut UploadStreamTicket) Sign(sk string) (string, error) {
	var expires_at time.Time = ut.GenerousTimeToUpload()
	return GenerateUploadStreamTicket(&ut, expires_at, sk)
}

type UploadStreamTicketClaims struct {
	UploadStreamTicket
	jwt.StandardClaims
}

func NewUploadStreamTicketFromRequest(request *http.Request) (*UploadStreamTicket, error) {
	var upload_ticket *UploadStreamTicket = new(UploadStreamTicket)
	upload_ticket.UploadUUID = request.URL.Query().Get("upload_uuid")

	upload_ticket_total_medias, err := strconv.Atoi(request.URL.Query().Get("total_medias"))
	if err != nil {
		return nil, fmt.Errorf("In getUploadStreamsTicketHandler: Invalid total_medias query parameter because '%s'", err.Error())
	}

	upload_ticket.TotalMedias = upload_ticket_total_medias

	upload_ticket.UploadedMedias = 0

	return upload_ticket, nil
}

func ParseUploadStreamTicketFromRequest(request *http.Request, sk string) (*UploadStreamTicket, error) {
	var upload_stream_ticket_cookie *http.Cookie
	upload_stream_ticket_cookie, err := request.Cookie(app_config.UPLOAD_STREAM_TICKET_COOKIE_NAME)
	if err != nil {
		return nil, fmt.Errorf("In ParseUploadStreamTicketFromRequest: Error getting upload stream ticket cookie because '%s'", err.Error())
	}

	return ParseUploadStreamTicket(upload_stream_ticket_cookie.Value, sk)
}

func GenerateUploadStreamTicket(upload_ticket *UploadStreamTicket, expires_at time.Time, sk string) (string, error) {
	upload_ticket_claims := &UploadStreamTicketClaims{
		UploadStreamTicket: *upload_ticket,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires_at.Unix(),
		},
	}

	token := jwt.NewWithClaims(dungeon_models.JwtSigningMethod, upload_ticket_claims)
	return token.SignedString([]byte(sk))
}

func ParseUploadStreamTicket(token string, sk string) (*UploadStreamTicket, error) {
	claims := &UploadStreamTicketClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(sk), nil
	})

	return &claims.UploadStreamTicket, err
}

// ========================= Chunked Upload Ticket =========================

type ChunkedUploadTicket struct {
	UploadUUID     string `json:"upload_uuid"`
	UploadFilename string `json:"upload_filename"`
	UploadSize     int64  `json:"upload_size"`
	UploadChunks   int    `json:"upload_chunks"`
	CategoryUUID   string `json:"category_uuid"`
}

func (ut ChunkedUploadTicket) GetFileFragmentName(chunk_serial int) string {
	var file_fragment_name string = fmt.Sprintf("%s.part%d", ut.UploadUUID, chunk_serial)

	return file_fragment_name
}

func (ut ChunkedUploadTicket) EqualUploadTicket(other ChunkedUploadTicket) bool {
	var equal bool = ut.UploadUUID == other.UploadUUID

	equal = equal && ut.UploadFilename == other.UploadFilename
	equal = equal && ut.UploadSize == other.UploadSize
	equal = equal && ut.UploadChunks == other.UploadChunks
	equal = equal && ut.CategoryUUID == other.CategoryUUID

	return equal
}

func NewChunkedUploadTicketFromRequest(request *http.Request) (*ChunkedUploadTicket, error) {
	var upload_ticket *ChunkedUploadTicket = new(ChunkedUploadTicket)
	upload_ticket.UploadUUID = request.URL.Query().Get("upload_uuid")
	upload_ticket.UploadFilename = request.URL.Query().Get("upload_filename")

	upload_ticket_upload_size, err := strconv.ParseInt(request.URL.Query().Get("upload_size"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("In getUploadStreamsTicketHandler: Invalid upload_size query parameter because '%s'", err.Error())
	}

	upload_ticket.UploadSize = upload_ticket_upload_size

	upload_ticket_upload_chunks, err := strconv.Atoi(request.URL.Query().Get("upload_chunks"))
	if err != nil {
		return nil, fmt.Errorf("In getUploadStreamsTicketHandler: Invalid upload_chunks query parameter because '%s'", err.Error())
	}

	upload_ticket.UploadChunks = upload_ticket_upload_chunks

	upload_ticket.CategoryUUID = request.URL.Query().Get("category_uuid")

	return upload_ticket, nil
}

type ChunkedUploadTicketClaims struct {
	ChunkedUploadTicket
	jwt.StandardClaims
}

func GenerateChunkedUploadTicket(upload_ticket *ChunkedUploadTicket, expires_at time.Time, sk string) (string, error) {
	upload_ticket_claims := &ChunkedUploadTicketClaims{
		ChunkedUploadTicket: *upload_ticket,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires_at.Unix(),
		},
	}

	token := jwt.NewWithClaims(dungeon_models.JwtSigningMethod, upload_ticket_claims)
	return token.SignedString([]byte(sk))
}

func ParseChunkedUploadTicket(token string, sk string) (*ChunkedUploadTicket, error) {
	claims := &ChunkedUploadTicketClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(sk), nil
	})

	return &claims.ChunkedUploadTicket, err
}

func SerializeChunkedUploadTicketToken(token string) []byte {
	var token_length uint16 = uint16(len(token))
	var token_bytes = make([]byte, 2)

	token_bytes[0] = byte(token_length)
	token_bytes[1] = byte(token_length >> 8) // Little endian

	token_bytes = append(token_bytes, []byte(token)...)

	return token_bytes
}

// Reads the upload ticket token from a file(or any io.ReadSeeker) that has the token serialized at the beginning.
// Returns the upload ticket and an error if any. The seek pointer will be at the end of the token.
func ReadChunkedUploadTicketToken(r io.ReadSeeker, sk string) (*ChunkedUploadTicket, error) {
	var token_length_bytes = make([]byte, 2)
	_, err := r.Read(token_length_bytes)
	if err != nil {
		return nil, err
	}

	var token_length uint16 = uint16(token_length_bytes[0]) | uint16(token_length_bytes[1])<<8

	var token_bytes = make([]byte, token_length)

	_, err = r.Read(token_bytes)
	if err != nil {
		return nil, err
	}

	var token string = string(token_bytes)

	return ParseChunkedUploadTicket(token, sk)
}
