package app_config

import (
	"encoding/json"
	"fmt"
	"libery-dungeon-libs/communication"
	"libery-dungeon-libs/models/platform_services"
	"os"

	"github.com/Gerardo115pp/patriots_lib/echo"
) // Loads the configuration from the environment variables

var SERVICE_PORT string = os.Getenv("SERVICE_PORT")
var SETTINGS_FILE string = os.Getenv("SETTINGS_FILE")
var OPERATION_DATA_PATH string = os.Getenv("OPERATION_DATA_PATH")
var JWT_SECRET string = os.Getenv("JWT_SECRET")
var DOMAIN_SECRET string = os.Getenv("DOMAIN_SECRET")

// --------Settings--------

var SERVICE_NAME platform_services.PlatformServiceName

var service_settings map[string]any = make(map[string]any)

func VerifyConfig() {

	if SERVICE_PORT == "" {
		panic("SERVICE PORT environment variable is required")
	}

	if OPERATION_DATA_PATH == "" {
		panic("OPERATION_DATA_PATH environment variable is required")
	}

	if JWT_SECRET == "" {
		panic("JWT_SECRET environment variable is required")
	}

	if DOMAIN_SECRET == "" {
		panic("DOMAIN_SECRET environment variable is required")
	}

	if SETTINGS_FILE == "" {
		SETTINGS_FILE = "settings.json"
	}

	SERVICE_NAME = platform_services.PWEDITOR_SERVICE

	err := loadSettings()
	if err != nil {
		panic(fmt.Sprintf("Error loading settings file: %s", err.Error()))
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

	return nil
}

func setupPlatformCommunication() (err error) {
	echo.Echo(echo.SkyBlueFG, fmt.Sprintf("Setting up communication:\\nCaPath: %s", communication.SSL_CA_PATH))

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

	err = communication.JD.NotifyServiceOnline(SERVICE_NAME, communication.PWEDITOR_SERVER, SERVICE_PORT)
	if err != nil {
		echo.EchoWarn(fmt.Sprintf("JD did not accept service online notification: %s", err.Error()))
	}

	return nil
}

func ClosePlatformCommunication() {
	communication.JD.NotifyServiceOffline(SERVICE_NAME)
}
