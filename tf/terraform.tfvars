// docker images tags

api_tag     = "latest"
backend_tag = "latest"

api_replicas     = 1
backend_replicas = 1


// envs

// update_config.monitor = API_STARTS_COUNT * (API_STARTS_DELAY+1)

// API зависит от Backend 
// Эмулируем медленное завершение Backend - 20 сек

API_SHUTDOWN     = 1
API_STARTS_COUNT = 20
API_STARTS_DELAY = 2

BACKEND_SHUTDOWN     = 20
BACKEND_STARTS_COUNT = 10
BACKEND_STARTS_DELAY = 1


API_VERSION = "v2"


STOP_TIMEOUT = 30
