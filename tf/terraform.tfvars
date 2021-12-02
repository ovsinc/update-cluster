// docker

api_tag     = "latest"
backend_tag = "latest"

api_replicas     = 1
backend_replicas = 1


// envs

LISTEN_PORT = 80
API_PORT    = 8000

// используется в stop_grace_period сервисов и
// для программного прерывания процесса остановки (gracefull shutdown)
STOP_TIMEOUT = 30


// API
// API зависит от Backend

// Эмулирует время остановки сервиса API
API_SHUTDOWN = 1
// количество попыток запуска
API_STARTS_COUNT = 4
// перерыв между попытками
API_STARTS_DELAY = 2


// BACKEND

// Эмулирует время остановки сервиса Backend
BACKEND_SHUTDOWN = 2

// количество попыток запуска
BACKEND_STARTS_COUNT = 2
// перерыв между попытками
BACKEND_STARTS_DELAY = 1
