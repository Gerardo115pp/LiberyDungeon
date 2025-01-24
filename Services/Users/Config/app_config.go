package app_config

import (
	"encoding/json"
	"fmt"
	"libery-dungeon-libs/communication"
	"libery-dungeon-libs/models/platform_services"
	"os"
	"time"

	"github.com/Gerardo115pp/patriots_lib/echo"
) // Loads the configuration from the environment variables

var SERVICE_PORT string = os.Getenv("SERVICE_PORT")
var SETTINGS_FILE string = os.Getenv("SETTINGS_FILE")
var OPERATION_DATA_PATH string = os.Getenv("OPERATION_DATA_PATH")
var SERVICE_NAME platform_services.PlatformServiceName
var DEBUG_MODE bool = os.Getenv("EDEBUG") != ""
var USERS_GRPC_SERVER_PORT string = os.Getenv("USERS_GRPC_SERVER_PORT")
var DOMAIN_SECRET string = os.Getenv("DOMAIN_SECRET")
var JWT_SECRET string = os.Getenv("JWT_SECRET")

// --------Settings--------

var service_settings map[string]any = make(map[string]any)

var LOCALTIME string = "America/Mexico_City"

var USER_CLAIMS_EXPIRATION_HOURS time.Duration = 7 * (24 * time.Hour)

var USER_CLAIMS_COOKIE_NAME string = "auth_user_claims"
var INITIAL_SETUP_SECRET string = ""

func VerifyConfig() {

	if SERVICE_PORT == "" {
		panic("SERVICE PORT environment variable is required")
	}

	if OPERATION_DATA_PATH == "" {
		panic("OPERATION_DATA_PATH environment variable is required")
	}

	if SETTINGS_FILE == "" {
		SETTINGS_FILE = "settings.json"
	}

	if USERS_GRPC_SERVER_PORT == "" {
		panic("USERS_GRPC_SERVER_PORT environment variable is required")
	}

	if DOMAIN_SECRET == "" {
		panic("DOMAIN_SECRET environment variable is required")
	}

	if JWT_SECRET == "" {
		panic("JWT_SECRET environment variable is required")
	}

	if custom_localtime := os.Getenv("LOCALTIME"); custom_localtime != "" {
		LOCALTIME = custom_localtime
	}

	SERVICE_NAME = platform_services.USERS_SERVICE

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

	if _, exists := settings["LOCALTIME"]; exists {
		LOCALTIME = settings["LOCALTIME"].(string)
	}

	if _, exists := settings["INITIAL_SETUP_SECRET"]; exists {
		INITIAL_SETUP_SECRET = settings["INITIAL_SETUP_SECRET"].(string)
	} else {
		panic(fmt.Sprintf("INITIAL_SETUP_SECRET is not set on '%s'", settings_path))
	}

	service_settings = settings

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

	err = communication.JD.NotifyServiceOnline(SERVICE_NAME, communication.USERS_SERVER, SERVICE_PORT)
	if err != nil {
		echo.EchoWarn(fmt.Sprintf("JD did not accept service online notification: %s", err.Error()))
	}

	return nil
}

func ClosePlatformCommunication() {
	communication.JD.NotifyServiceOffline(SERVICE_NAME)
}
