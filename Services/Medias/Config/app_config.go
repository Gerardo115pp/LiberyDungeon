package app_config

import (
	"encoding/json"
	"fmt"
	"libery-dungeon-libs/communication"
	"libery-dungeon-libs/models/platform_services"
	"os"
	"path/filepath"

	"github.com/Gerardo115pp/patriots_lib/echo"
) // Loads the configuration from the environment variables

var SERVICE_PORT string = os.Getenv("SERVICE_PORT")
var SERVICE_NAME platform_services.PlatformServiceName
var SETTINGS_FILE string = os.Getenv("SETTINGS_FILE")
var OPERATION_DATA_PATH string = os.Getenv("OPERATION_DATA_PATH")
var UPLOAD_CHUNKS_PATH string
var TRASH_STORAGE_PATH string = os.Getenv("TRASH_STORAGE_PATH")
var TRASH_MEDIA_PATH string = filepath.Join(TRASH_STORAGE_PATH, "medias")

// --------Secrets--------

var DOMAIN_SECRET string = os.Getenv("DOMAIN_SECRET")
var JWT_SECRET string = os.Getenv("JWT_SECRET")

// --------Internal headers--------

// IH = Internal Header.
var IH_DOWNLOAD_RECIPIENT_PATH string = os.Getenv("IH_DOWNLOAD_RECIPIENT_PATH")

// --------Cookies names--------
// Cookies shared between services are environment variables.
// Cookies that are unique to a service are defined
// as hard coded values.

// The name of the cookie that contains the cluster access jwt token(a json representation of a category cluster).
var CATEGORIES_CLUSTER_ACCESS_COOKIE_NAME string = os.Getenv("CATEGORIES_CLUSTER_ACCESS_COOKIE_NAME")
var UPLOAD_STREAM_TICKET_COOKIE_NAME string = "upload-stream-ticket"

// --------Settings--------

var service_settings map[string]any = make(map[string]any)

// The default thumbnail width
var THUMBNAIL_WIDTH int = 100
var MIN_THUMBNAIL_WIDTH int = 50
var MOBILE_MAX_WIDTH int = 580

func VerifyConfig() {

	if SERVICE_PORT == "" {
		panic("SERVICE PORT environment variable is required")
	}

	SERVICE_NAME = platform_services.MEDIAS_SERVICE

	if OPERATION_DATA_PATH == "" {
		panic("OPERATION_DATA_PATH environment variable is required")
	}

	UPLOAD_CHUNKS_PATH = filepath.Join(OPERATION_DATA_PATH, "upload_chunks")

	if SETTINGS_FILE == "" {
		SETTINGS_FILE = "settings.json"
	}

	if DOMAIN_SECRET == "" {
		panic("DOMAIN_SECRET environment variable is required")
	}

	if JWT_SECRET == "" {
		panic("JWT_SECRET environment variable is required")
	}

	if TRASH_STORAGE_PATH == "" {
		panic("TRASH_STORAGE_PATH environment variable is required")
	}

	err := loadSettings()
	if err != nil {
		panic(fmt.Sprintf("Error loading settings file: %s", err.Error()))
	}

	verifyDirectoryStructure()

	err = setupPlatformCommunication()
	if err != nil {
		panic(fmt.Sprintf("Error setting up platform communication: %s", err.Error()))
	}
}

func verifyDirectoryStructure() {
	echo.EchoDebug(fmt.Sprintf("Verifying directory Operation Data Path<%s>", OPERATION_DATA_PATH))
	if !verifyDirectoryExists(OPERATION_DATA_PATH) {
		os.Mkdir(OPERATION_DATA_PATH, 0755)
		echo.Echo(echo.YellowFG, fmt.Sprintf("Directory<%s> does not existed. Created", OPERATION_DATA_PATH))
	}

	echo.EchoDebug(fmt.Sprintf("Verifying directory Upload Chunks Path<%s>", UPLOAD_CHUNKS_PATH))
	if !verifyDirectoryExists(UPLOAD_CHUNKS_PATH) {
		os.Mkdir(UPLOAD_CHUNKS_PATH, 0755)
		echo.Echo(echo.YellowFG, fmt.Sprintf("Directory<%s> does not existed. Created", UPLOAD_CHUNKS_PATH))
	}
}

// Checks if a directory exists. If it does exist but is not a directory it panics.
func verifyDirectoryExists(directory_path string) bool {
	var directory_exists bool = true
	stat, err := os.Stat(directory_path)
	if os.IsNotExist(err) {
		directory_exists = false
	}

	if directory_exists && !stat.IsDir() {
		panic(fmt.Sprintf("Path<%s> exists but is not a directory", directory_path))
	}

	return directory_exists
}

func loadSettings() error {
	// Load settings from file
	var settings_path string = fmt.Sprintf("%s/%s", OPERATION_DATA_PATH, SETTINGS_FILE)
	var err error

	if _, err = os.Stat(settings_path); os.IsNotExist(err) {
		return fmt.Errorf("Settings file<%s> not found", settings_path)
	}

	var setting_content []byte

	setting_content, err = os.ReadFile(settings_path)

	var settings map[string]any = make(map[string]any)

	err = json.Unmarshal(setting_content, &settings)
	if err != nil {
		return fmt.Errorf("While unmarshaling settings data, found error <%s>", err.Error())
	}

	service_settings = settings

	if _, exists := service_settings["THUMBNAIL_WIDTH"]; exists {
		THUMBNAIL_WIDTH = int(service_settings["THUMBNAIL_WIDTH"].(float64))
	}

	if _, exists := service_settings["MIN_THUMBNAIL_WIDTH"]; exists {
		MIN_THUMBNAIL_WIDTH = int(service_settings["MIN_THUMBNAIL_WIDTH"].(float64))
	}

	if _, exists := service_settings["MOBILE_MAX_WIDTH"]; exists {
		MOBILE_MAX_WIDTH = int(service_settings["MOBILE_MAX_WIDTH"].(float64))
	}

	return nil
}

func setupPlatformCommunication() (err error) {
	echo.Echo(echo.SkyBlueFG, fmt.Sprintf("Setting up communication:\nCaPath: %s", communication.SSL_CA_PATH))

	err = communication.Setup()
	if err != nil {
		return
	}

	jd_alive, err := communication.JD.Alive()
	if err != nil {
		return
	}

	if !jd_alive {
		return fmt.Errorf("JD is unreachable")
	}

	err = communication.JD.NotifyServiceOnline(SERVICE_NAME, communication.DOWNLOADS_SERVER, SERVICE_PORT)
	if err != nil {
		echo.EchoWarn(fmt.Sprintf("JD did not accept service online notification: %s", err.Error()))
	}

	return nil
}

func ClosePlatformCommunication() {
	communication.JD.NotifyServiceOffline(SERVICE_NAME)
}
