from os import getenv, path
import json

#### GENERAL ####

SERVICE_PORT = getenv('SERVICE_PORT', "")
assert SERVICE_PORT != "", "SERVICE_PORT is not set"

SERVICE_NAME = getenv('SERVICE_NAME', "")
assert SERVICE_NAME != "", "SERVICE_NAME is not set"

OPERATION_DATA_PATH = getenv('OPERATION_DATA_PATH', "")
assert OPERATION_DATA_PATH != "", "OPERATION_DATA_PATH is not set"

DEVELOPMENT_MODE = getenv('DEVELOPMENT', "0") == "1"    

DOWNLOAD_DEFAULT_CATEGORY = getenv('DOWNLOAD_DEFAULT_CATEGORY', "")
assert DOWNLOAD_DEFAULT_CATEGORY != "", "DOWNLOAD_DEFAULT_CATEGORY is not set"

#### SECRETS ####

DOMAIN_SECRET = getenv('DOMAIN_SECRET', "")
assert DOMAIN_SECRET != "", "DOMAIN_SECRET is not set"

JWT_SECRET = getenv('JWT_SECRET', "")
assert JWT_SECRET != "", "JWT_SECRET is not set"

#### MYSQL ####

MYSQL_DB = getenv('MYSQL_DB', "")
assert MYSQL_DB != "", "MYSQL_DB is not set"

MYSQL_PORT = getenv('MYSQL_PORT', "")
assert MYSQL_PORT != "", "MYSQL_PORT is not set"

MYSQL_HOST = getenv('MYSQL_HOST', "")
assert MYSQL_HOST != "", "MYSQL_HOST is not set"

MYSQL_USER = getenv('MYSQL_USER', "")
assert MYSQL_USER != "", "MYSQL_USER is not set"

MYSQL_PASSWORD = getenv('MYSQL_PASSWORD', "")
assert MYSQL_PASSWORD != "", "MYSQL_PASSWORD is not set"

#### GRPC ####  

GRPC_SERVER = getenv('GRPC_SERVER', "")
assert GRPC_SERVER != "", "GRPC_SERVER is not set"


_local_settings = getenv('LOCAL_SETTINGS_PATH', "") or path.join(OPERATION_DATA_PATH, "settings.json")

def loadLocalSettings(config_path) -> dict[str, str]:
    if not path.exists(config_path):
        raise Exception(f"Local config file not found on path: {config_path}")
    
    config = {}    
    with open(_local_settings, 'r') as f:
        config = json.load(f)
    
    return config

def saveLocalSettings(config: dict[str, str]):
    with open(_local_settings, 'w') as f:
        json.dump(config, f)

LOCAL_SETTINGS = loadLocalSettings(_local_settings)
assert 'tracked_boards' in LOCAL_SETTINGS, "tracked_boards is not set in local settings"
TRACKED_BOARDS = LOCAL_SETTINGS['tracked_boards']
assert 'board_catalog_template' in LOCAL_SETTINGS, "board_catalog_template is not set in local settings"
BOARD_CATALOG_TEMPLATE = LOCAL_SETTINGS['board_catalog_template']
assert 'thread_cover_image_template' in LOCAL_SETTINGS, "thread_cover_image_template is not set in local settings"
THREAD_COVER_IMAGE_TEMPLATE = LOCAL_SETTINGS['thread_cover_image_template']
assert 'thread_template' in LOCAL_SETTINGS, "thread_template is not set in local settings"
THREAD_TEMPLATE = LOCAL_SETTINGS['thread_template']
assert 'catalog_threads_lookup_string' in LOCAL_SETTINGS, "catalog_threads_lookup_string is not set in local settings"
CATALOG_THREADS_LOOKUP_STRING = LOCAL_SETTINGS['catalog_threads_lookup_string']