package app_config

import (
	"encoding/json"
	"fmt"
	"libery-dungeon-libs/communication"
	dungeon_helpers "libery-dungeon-libs/helpers"
	"libery-dungeon-libs/models/platform_services"
	"os"

	"github.com/Gerardo115pp/patriots_lib/echo"
) // Loads the configuration from the environment variables

var SERVICE_PORT string = os.Getenv("SERVICE_PORT")
var SERVICE_NAME platform_services.PlatformServiceName
var DEBUG_MODE bool = os.Getenv("EDEBUG") != ""
var CATEGORIES_GRPC_SERVER_PORT string = os.Getenv("CATEGORIES_GRPC_SERVER_PORT")
var SETTINGS_FILE string = os.Getenv("SETTINGS_FILE")
var OPERATION_DATA_PATH string = os.Getenv("OPERATION_DATA_PATH")
var SERVICE_CLUSTERS_ROOT string = os.Getenv("SERVICE_CLUSTERS_ROOT")
var MYSQL_USER string = os.Getenv("MYSQL_USER")
var MYSQL_PASSWORD string = os.Getenv("MYSQL_PASSWORD")
var MYSQL_HOST string = os.Getenv("MYSQL_HOST")
var MYSQL_DB string = os.Getenv("MYSQL_DB")
var MYSQL_PORT string = os.Getenv("MYSQL_PORT")
var MAIN_CATEGORY_PROXY_ID string = os.Getenv("MAIN_CATEGORY_PROXY_ID")
var DOMAIN_SECRET string = os.Getenv("DOMAIN_SECRET")
var JWT_SECRET string = os.Getenv("JWT_SECRET")
var TRASH_STORAGE_PATH string = os.Getenv("TRASH_STORAGE_PATH")

// --------Cookies names--------

var CATEGORIES_CLUSTER_ACCESS_COOKIE_NAME string = os.Getenv("CATEGORIES_CLUSTER_ACCESS_COOKIE_NAME")

// --------App routes--------

var MEDIAS_APP_DUNGEON_EXPLORER_ROUTE string = os.Getenv("MEDIAS_APP_DUNGEON_EXPLORER_ROUTE")

// --------Settings--------

var service_settings map[string]any = make(map[string]any)

var LOCALTIME string = "America/Mexico_City"
var SHARED_MEDIA_EXPIRATION_SECS int64 = 3600

func VerifyConfig() {

	if SERVICE_PORT == "" {
		panic("SERVICE PORT environment variable is required")
	}

	if CATEGORIES_GRPC_SERVER_PORT == "" {
		panic("CATEGORIES_GRPC_SERVER_PORT environment variable is required")
	}

	if OPERATION_DATA_PATH == "" {
		panic("OPERATION_DATA_PATH environment variable is required")
	}

	if SERVICE_CLUSTERS_ROOT == "" {
		panic("SERVICE_CLUSTERS_ROOT environment variable is required")
	}

	SERVICE_CLUSTERS_ROOT = dungeon_helpers.NormalizePath(SERVICE_CLUSTERS_ROOT)

	if SETTINGS_FILE == "" {
		SETTINGS_FILE = "settings.json"
	}

	if MYSQL_USER == "" {
		panic("MYSQL_USER environment variable is required")
	}

	if MYSQL_PASSWORD == "" {
		panic("MYSQL_PASSWORD environment variable is required")
	}

	if MYSQL_HOST == "" {
		panic("MYSQL_HOST environment variable is required")
	}

	if MYSQL_DB == "" {
		panic("MYSQL_DB environment variable is required")
	}

	if MYSQL_PORT == "" {
		panic("MYSQL_PORT environment variable is required")
	}

	if MAIN_CATEGORY_PROXY_ID == "" {
		MAIN_CATEGORY_PROXY_ID = "main"
	}

	if DOMAIN_SECRET == "" {
		panic("DOMAIN_SECRET environment variable is required")
	}

	if JWT_SECRET == "" {
		panic("JWT_SECRET environment variable is required")
	}

	if CATEGORIES_CLUSTER_ACCESS_COOKIE_NAME == "" {
		panic("CATEGORIES_CLUSTER_ACCESS_COOKIE_NAME environment variable is required")
	}

	if MEDIAS_APP_DUNGEON_EXPLORER_ROUTE == "" {
		MEDIAS_APP_DUNGEON_EXPLORER_ROUTE = "/dungeon-explorer"
		fmt.Printf("WARNING: MEDIAS_APP_DUNGEON_EXPLORER_ROUTE not set, using default: %s\n", MEDIAS_APP_DUNGEON_EXPLORER_ROUTE)
	}

	SERVICE_NAME = platform_services.CATEGORIES_SERVICE

	err := loadSettings()
	if err != nil {
		panic(fmt.Sprintf("Error loading settings file: %s", err.Error()))
	}

	err = setupPlatformCommunication()
	if err != nil {
		panic(fmt.Sprintf("Error setting up platform communication: %s", err.Error()))
	}
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

	if _, exists := service_settings["LOCALTIME"]; exists {
		LOCALTIME = service_settings["LOCALTIME"].(string)
	}

	if _, exists := service_settings["SHARED_MEDIA_EXPIRATION_SECS"]; exists {
		SHARED_MEDIA_EXPIRATION_SECS = int64(service_settings["SHARED_MEDIA_EXPIRATION_SECS"].(float64))
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
